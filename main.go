package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	 "encoding/json"
	 "net"
	 "sync"
	"github.com/joho/godotenv"

	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/sd"

	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	"github.com/luraproject/lura/v2/router/gin"
)
var mu = sync.Mutex{}
type configj struct {
	Version   int         `json:"version"`
	Port      int         `json:"port"`
	Cache_ttl string      `json:"cache_ttl"`
	Timeout   string      `json:"timeout"`
	Endpoints []endpoints `json:"endpoints"`
}
type endpoints struct {
	Endpoint string    `json:"endpoint"`
	Method   string    `json:"method"`
	Backend  []backend `json:"backend"`
}
type backend struct {
	Url_pattern  string      `json:"url_pattern"`
	Encoding     string      `json:"encoding"`
	Method       string      `json:"method"`
	Host         []string    `json:"host"`
	Extra_config Extraconfig `json:"extra_config"`
}
type Extraconfig map[string]interface{}

type Server struct {
	url string
}
type customProxyFactory struct {
	logger  logging.Logger
	factory proxy.Factory
}

type defaultFactory struct {
	backendFactory    proxy.BackendFactory
	logger            logging.Logger
	subscriberFactory sd.SubscriberFactory
}

var (
	LOGGING_TYPE  string
	LOGGING_NAME  string
	CONFIGURATION string
	DEBUG         string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	LOGGING_TYPE = os.Getenv("LOGGING_TYPE")
	LOGGING_NAME = os.Getenv("LOGGING_NAME")
	CONFIGURATION = os.Getenv("CONFIGURATION")
	DEBUG = os.Getenv("DEBUG")
}

func main() {
	// Loand config json
	parser := config.NewParser()

	dir, err := os.Getwd()
	fmt.Println(dir)
	if err != nil {
		log.Panicln(err.Error())
	}

	file := fmt.Sprintf("%v/%v", dir, CONFIGURATION)
	fmt.Println(file)
	serviceConfig, err := parser.Parse(file)
	if err != nil {
		log.Fatalf("Config error. Detail: %v\n", err.Error())
	}

	// Config logs level
	debug, _ := strconv.ParseBool(DEBUG)
	serviceConfig.Debug = serviceConfig.Debug || debug

	logger, err := logging.NewLogger(LOGGING_TYPE, os.Stdout, LOGGING_NAME)
	if err != nil {
		log.Fatalf("Log error. Detail: %v\n", err.Error())
	}

	// Middlewares
	/*secureMiddleware := secure.New(secure.Options{
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
	})
	*/

	cfg := gin.DefaultFactory(proxy.DefaultFactory(logger), logger)
	go tcpConn(serviceConfig)
	cfg.New().Run(serviceConfig)

}

func tcpConn(cfg config.ServiceConfig) {
	l, _ := net.Listen("tcp", ":7001")
	fmt.Println("Listen tcp")
	defer l.Close()

	for {
		conn, _ := l.Accept()

		defer conn.Close()
		go ConnHandler(conn, cfg)
	}

}

func ConnHandler(conn net.Conn, cfg config.ServiceConfig) {
	defer conn.Close()
	fmt.Println("Ok")
	data := map[string]string{}

	json.NewDecoder(conn).Decode(&data)
	cleaner := config.NewURIParser()
	mu.Lock()
	   for _, j := range cfg.Endpoints {
                                                           for k, i := range j.Backend {
                                                                                           if k == 0 {
												   for key, _:= range data {
                                                                                                 for num, host := range i.Host  {
                                                                                                                        if key == "zombieIp"{
                                                                                                                           if host ==  cleaner.CleanHost(data["zombieIp"]) {    
                                                                                                                                        i.Host = append(i.Host[:num], i.Host[num+1:]...)
                                                                                                                                       // i.Host = append(i.Host,i.Host[len(i.Host)-1])                          
                                                                                                                                                                                                    
                                                                                                                                     }
                                                                                                                                        }else if key == "newIp" && num == len(i.Host)-1{
                                                                                                                                                i.Host = append(i.Host, cleaner.CleanHost(data["newIp"]))      
  
                                                                                                                                        }
                                                                                                                                }
           }
                  
                   }                                                                                                                                                               }                           
                                                                     
         } 
	mu.Unlock()
	fmt.Println(cfg.Endpoints[0].Backend[0].Host)
	conn.Write([]byte("res"))	
}

