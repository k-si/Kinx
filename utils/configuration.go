package utils

//var Config *Configuration
//
//type Configuration struct {
//	// Server配置
//	Name    string
//	Host    string
//	TcpPort int
//
//	// Kinx配置
//	Version           string
//	MaxConn           int    // 最多可有多少连接
//	MaxPackage        uint32 // 一次可发送包的最大size
//	WorkerPoolSize    uint32 // 多少个worker
//	MaxWorkerTaskSize uint32 //
//}
//
//func isExist(filePath string) bool {
//	_, err := os.Stat(filePath)
//	if os.IsNotExist(err) {
//		return false
//	} else {
//		return true
//	}
//}
//
//func (c *Configuration) reload() {
//	filePath := "conf/kinx.json"
//	if isExist(filePath) {
//		data, err := ioutil.ReadFile(filePath)
//		if err != nil {
//			fmt.Println("read kinx.json fail")
//		}
//		err = json.Unmarshal(data, &Config)
//		if err != nil {
//			fmt.Println("parse kinx.json fail")
//		}
//	}
//}
//
//func (c *Configuration) banner() {
//	filePath := "conf/kinx.banner"
//	if isExist(filePath) {
//		data, err := ioutil.ReadFile(filePath)
//		if err != nil {
//			fmt.Println("read kinx.banner fail")
//		}
//		fmt.Println(string(data))
//	}
//}
//
//func init() {
//	Config = &Configuration{
//		Name:              "TcpServerApp",
//		Host:              "0.0.0.0",
//		TcpPort:           8888,
//		Version:           "0.8",
//		MaxConn:           3,
//		MaxPackage:        1024,
//		WorkerPoolSize:    5,
//		MaxWorkerTaskSize: 1024,
//	}
//
//	// 加载配置文件
//	Config.reload()
//
//	// 显示banner
//	Config.banner()
//
//	fmt.Println("[server load configuration...]")
//	fmt.Println("[server MAX connection number]:", Config.MaxConn)
//	fmt.Println("[server worker pool size]:", Config.WorkerPoolSize)
//}
