package main

import (
	"log"
	"net/http"
	"os"

	"github.com/akazwz/file-upload/initialize"
	"github.com/joho/godotenv"
)

func main() {
	r := initialize.InitRouter()
	// 读取环境变量配置
	InitEnvConfig()

	port := os.Getenv("API_PORT")
	log.Println("PORT:" + port)

	s := &http.Server{
		Addr:    port,
		Handler: r,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalln("Api启动失败")
	}
}

func InitEnvConfig() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatalln("读取配置文件失败")
		}
	}
}
