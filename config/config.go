package config

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"os"
)

var EnvApplication string
var RootPath string

var AppConfig *appConfig
var RouterConfig *routerConfig
var MuxRoute *mux.Router
var ViewConfigModule *viewConfigModule

type appConfig struct {
	Port          string
	RootPath      string
	StaticPath    string
	UploadPath    string
	SiteUrl       string
	Protocol      string
	StaticUrl     string
	StaticVersion string
	UploadUrl     string
}

type routerConfig struct {
	App struct {
		Index struct {
			Name    string
			Pattern string
			Params  []string
			Queries []string
		}
		Cate struct {
			Name    string
			Pattern string
			Params  []string
			Queries []string
		}
		Post struct {
			Name    string
			Pattern string
			Params  []string
			Queries []string
		}
		Chapter struct {
			Name    string
			Pattern string
			Params  []string
			Queries []string
		}
		Author struct {
			Name    string
			Pattern string
			Params  []string
			Queries []string
		}
		Search struct {
			Name    string
			Pattern string
			Params  []string
			Queries []string
		}
	}
	Admin struct {
		Index struct {
			Name    string
			Pattern string
		}
		Cate struct {
			Name    string
			Pattern string
		}
		Post struct {
			Name    string
			Pattern string
		}
		Chapter struct {
			Name    string
			Pattern string
		}
	}
}

type Database struct {
	Host     string
	Port     string
	Driver   string
	UserName string
	Password string
	DBName   string
}

type Elastic struct {
	Host     string
	Port     string
	Type     string
	Protocol string
	Prefix   string
}

type Redis struct {
	Host     string
	Port     string
	DB       int
	Password string
}

type RabbitMQ struct {
	URL    string
	Prefix string
	Queue  []map[string]interface{}
}

type Route struct {
	PathView string
	Port     string
	Driver   string
	UserName string
	Password string
	DBName   string
}

type ArgumentsWorker struct {
	Action string
	Params map[string]interface{}
}

type ViewConfig struct {
	RootTemplate string
	Path         string
	Template     []string
}

type viewConfigModule struct {
	App struct {
		Templates    []string
		RootTemplate string
		View         struct {
			Index   string
			Cate    string
			Post    string
			Chapter string
			Author  string
			Search  string
		}
	}
	Admin struct{}
}

func InitLoad() {
	if AppConfig == nil {
		args := os.Args
		if len(args) < 2 {
			fmt.Println("Please input args")
			os.Exit(1)
		}

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		EnvApplication = args[1]
		RootPath = dir
		getAppConfig()
	}

	if RouterConfig == nil {
		getRouterConfig()
	}
}

func getRouterConfig() *routerConfig {
	if RouterConfig != nil {
		return RouterConfig
	}
	path := "config/router.json"
	plan, _ := ioutil.ReadFile(path)

	data := routerConfig{}
	err := json.Unmarshal(plan, &data)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	RouterConfig = &data
	return RouterConfig
}

func getAppConfig() *appConfig {
	if AppConfig != nil {
		return AppConfig
	}
	path := "config/" + EnvApplication + "/appConfig.json"
	plan, _ := ioutil.ReadFile(path)

	data := appConfig{}
	err := json.Unmarshal(plan, &data)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	data.RootPath = RootPath
	data.StaticPath = RootPath + "/static"
	data.UploadPath = data.StaticPath + "/static"
	AppConfig = &data
	return AppConfig
}

func GetConfigDB() Database {
	path := "config/" + EnvApplication + "/database.json"
	plan, _ := ioutil.ReadFile(path)
	data := Database{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

//
func GetConfigES() Elastic {
	path := "config/" + EnvApplication + "/elastic.json"
	plan, _ := ioutil.ReadFile(path)
	data := Elastic{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

//
func GetConfigRedis() Redis {
	path := "config/" + EnvApplication + "/redis.json"
	plan, _ := ioutil.ReadFile(path)
	data := Redis{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

//
func GetConfigRabbitMQ() RabbitMQ {
	path := "config/" + EnvApplication + "/rabbitMQ.json"
	plan, _ := ioutil.ReadFile(path)
	data := RabbitMQ{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

/**
Get config from file json
*/
func GetViewConfigModule() *viewConfigModule {
	if ViewConfigModule != nil {
		return ViewConfigModule
	}

	path := "config/view.json"
	plan, _ := ioutil.ReadFile(path)

	data := viewConfigModule{}
	err := json.Unmarshal(plan, &data)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ViewConfigModule = &data
	return ViewConfigModule
}
