package micro

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

type ServiceConfig struct {
	Name       string
	Port       int
	TargetPort int
}
