package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func init() {
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

func checkPods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking pods...")
	for {
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		realPods := pods.Items
		for i := 0; i < len(realPods); i++ {
			if realPods[i].GetNamespace() == "default" {
				pod, err := clientset.CoreV1().Pods("default").Get(realPods[i].GetName(), metav1.GetOptions{})
				if errors.IsNotFound(err) {
					fmt.Printf("Pod not found\n")
				} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
					fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
				} else if err != nil {
					panic(err.Error())
				} else {
					fmt.Printf("Found pod\n")
				}
				fmt.Println("pod cluster name => ", pod.GetClusterName(), "pod name => ", pod.GetName(), "pod namespace => ", pod.GetNamespace())
			}
		}

		time.Sleep(5 * time.Second)
	}
}

type PodObj struct {
	Name      string `json:"name"`
	NameSpace string `json:"name_space"`
}

func showPods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Show pods.")
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if pods == nil || pods.Items == nil || len(pods.Items) <= 0 {
		return
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	var podList []*PodObj
	for _, pod := range pods.Items {
		fmt.Println("Pod Name => ", pod.GetName())
		podObj := &PodObj{
			Name:      pod.GetName(),
			NameSpace: pod.GetNamespace(),
		}
		podList = append(podList, podObj)
	}
	podListMarshaled, _ := json.Marshal(podList)
	w.Write(podListMarshaled)
}

func int32Ptr(i int32) *int32 { return &i }

func newDeployment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New deployment.")
	deploymentClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "consul",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "consul",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "consul",
							Image: "ray-xyz.com:9090/consul",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "micro",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8500,
								},
							},
						},
					},
				},
			},
		},
	}
	log.Println("Creating consul deployment...")
	deploy, err := deploymentClient.Create(deployment)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Successfully deployed a deployment. Deployment name => ", deploy.GetObjectMeta().GetName())
}

func main() {
	if clientset == nil {
		panic("Client set is nil.")
	}

	port := "7878"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am in real micro kube server."))
	})
	http.HandleFunc("/showpods", showPods)
	http.HandleFunc("/checkpods", checkPods)
	http.HandleFunc("/newdeploy", newDeployment)
	log.Println("Server running on port => ", port)
	http.ListenAndServe(":"+port, nil)
}
