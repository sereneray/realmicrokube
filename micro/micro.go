package micro

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"realmicrokube/utils"
	"reflect"
	"strconv"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"k8s.io/api/apps/v1beta2"
	kbapiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var clientset *kubernetes.Clientset

func init() {
	log.Println("Initializing micro...")
	initKubeInCluster()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func initKubeOutofCluster() {
	// Out of cluster
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("Create client set error.")
		panic(err.Error())
	}
}

func initKubeInCluster() {
	// In cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

type Service struct {
	Config       *ServiceConfig
	NewClientRef interface{}
	KubeService  *KubeService
}

type KubeServiceDeployConfig struct {
	Namespace  string
	Name       string
	Port       int32
	TargetPort int32
	Image      string
	Replicas   int32
}

type KubeService struct {
	Namespace  string
	Name       string
	Port       int32
	TargetPort int32
	Endpoints  kbapiv1.Endpoints
}

func NewService(config *ServiceConfig, server interface{}, grpcRegisterServer interface{}) {
	listeningAddr := config.Host + ":" + strconv.Itoa(config.Port)
	listener, err := net.Listen("tcp", listeningAddr)
	log.Println("Service listening at =>", listeningAddr)
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}

	grpcServer := grpc.NewServer()
	var args []reflect.Value
	args = append(args, reflect.ValueOf(grpcServer))
	args = append(args, reflect.ValueOf(server))
	reflect.ValueOf(grpcRegisterServer).Call(args)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}
}

func DeployKubeService(deployment *KubeServiceDeployConfig) (success bool, desc string) {
	deploy, err := newKubeDeployment(deployment)
	if err != nil {
		return false, err.Error()
	}
	log.Println("Successfully deployed kube deployment, uid => ", deploy.GetUID())
	svc, err := newKubeService(&KubeService{
		Name:       deployment.Name,
		Namespace:  deployment.Namespace,
		Port:       deployment.Port,
		TargetPort: deployment.TargetPort,
	})
	log.Println("Successfully create kube service, uid => ", svc.GetObjectMeta().GetUID())
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func newKubeDeployment(deploy *KubeServiceDeployConfig) (*v1beta2.Deployment, error) {
	deployment := &v1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploy.Name,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
			Labels: map[string]string{
				"app": deploy.Name,
			},
		},
		Spec: v1beta2.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploy.Name,
				},
			},
			Replicas: utils.Int32Ptr(deploy.Replicas),
			Template: kbapiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploy.Name,
					},
				},
				Spec: kbapiv1.PodSpec{
					Containers: []kbapiv1.Container{
						{
							Name:  deploy.Name,
							Image: deploy.Image,
							Ports: []kbapiv1.ContainerPort{
								{
									Name:          "port",
									Protocol:      kbapiv1.ProtocolTCP,
									ContainerPort: int32(deploy.Port),
								},
							},
						},
					},
				},
			},
		},
	}
	dep, err := clientset.Apps().Deployments(deploy.Namespace).Create(deployment)
	if err != nil {
		return nil, err
	}
	return dep, nil
}

func newKubeService(service *KubeService) (*kbapiv1.Service, error) {
	svc, err := clientset.CoreV1().Services(service.Namespace).Create(&kbapiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: service.Name,
			Labels: map[string]string{
				"app": service.Name,
			},
		},
		Spec: kbapiv1.ServiceSpec{
			Type: kbapiv1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": service.Name,
			},
			Ports: []kbapiv1.ServicePort{
				{
					Port: int32(service.Port),
					TargetPort: intstr.IntOrString{
						IntVal: service.TargetPort,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func NewServiceClient(service string, newClientRef interface{}) (*Service, error) {
	if service == "" || newClientRef == nil {
		return nil, errors.New("Create service client error. Arguments nil.")
	}
	srv, err := queryKubeService("default", service)

	srvConf := &ServiceConfig{
		Host: srv.Spec.ClusterIP,
		Port: int(srv.Spec.Ports[0].Port),
	}
	if err != nil {
		return nil, err
	}

	var endpoints kbapiv1.Endpoints
	resp, err := http.Get("/api/v1/endpoints/" + service)
	content, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(content, &endpoints)
	log.Println("Service end points => ", endpoints)

	kubesvc := &KubeService{
		Namespace: srv.Namespace,
		Name:      srv.Name,
		Endpoints: endpoints,
	}

	return &Service{Config: srvConf, NewClientRef: newClientRef, KubeService: kubesvc}, nil
}

func queryKubeService(namespace, service string) (*kbapiv1.Service, error) {
	srv, err := clientset.CoreV1().Services(namespace).Get(service, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("I got the service => ", service)
	return srv, nil
}

func (s *Service) Call(method string, ctx context.Context, reqObj interface{}) (interface{}, error) {
	// Use grpc locad balancing strategy.
	// grpc.WithBalancer(grpc.RoundRobin(r))
	// if err := c.do(ctx, "GET", c.nsEndpoint()+"endpoints/"+serviceName, &res); err != nil {
	// 	return nil, err
	// }
	address := s.Config.Host + ":" + strconv.Itoa(s.Config.Port)
	conn, err := grpc.Dial(address)
	if err != nil {
		log.Println("Connection to server error.")
		return nil, err
	}
	if conn == nil {
		return nil, errors.New("Connection cannot be established.")
	}
	defer conn.Close()

	var client reflect.Value
	var newClientArgs []reflect.Value
	newClientArgs = append(newClientArgs, reflect.ValueOf(conn))
	newClientVals := reflect.ValueOf(s.NewClientRef).Call(newClientArgs)
	if newClientVals != nil && len(newClientVals) > 0 {
		client = newClientVals[0]
	}

	if client.IsNil() {
		return nil, errors.New("Parse grpc client error.")
	}

	var methodArgs []reflect.Value
	methodArgs = append(methodArgs, reflect.ValueOf(ctx))
	methodArgs = append(methodArgs, reflect.ValueOf(reqObj))
	// Call grpc method
	methodVals := client.MethodByName(method).Call(methodArgs)

	var respResult interface{}
	var respError error
	if methodVals != nil && len(methodVals) > 0 {
		if methodVals[0].CanInterface() {
			if methodVals[0].Interface() != nil {
				respResult = methodVals[0].Interface()
			}
		}
	}
	if methodVals != nil && len(methodVals) > 1 {
		if methodVals[1].CanInterface() {
			if methodVals[1].Interface() != nil {
				respError = methodVals[1].Interface().(error)
			}
		}
	}

	log.Printf("RespResult => %#v RespError => %#v", respResult, respError)

	return respResult, respError
}
