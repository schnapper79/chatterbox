package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/schnapper79/chatterbox"
)

func main() {
	ModelPath := "./models"
	PathToLLama := ""

	ModelPath_ENV := os.Getenv("MODEL_PATH")
	if ModelPath_ENV != "" {
		ModelPath = ModelPath_ENV
	}

	if _, err := os.Stat(ModelPath); os.IsNotExist(err) {
		err := os.MkdirAll(ModelPath, 0755) // 0755 are the UNIX permissions
		if err != nil {
			panic(err)
		}
		fmt.Println("Directory created.")
	} else {
		fmt.Println("Directory exists.")
	}

	var startmodel string
	var host string
	flag.StringVar(&startmodel, "s", "", "Start model")
	flag.StringVar(&PathToLLama, "llama", "./llama.cpp", "Path to llama.cpp")
	flag.StringVar(&host, "host", ":8080", "Host")

	flag.Parse()

	server := chatterbox.GetServer(ModelPath, PathToLLama, host)

	if startmodel != "" {
		server.LoadModellFromFile(startmodel)
	}

	log.Printf("Server started on %s\n", host)
	err := server.Server.ListenAndServe()
	if err != nil {
		if err.Error() == "http: Server closed" {
			log.Println("Server closed")
		} else {
			log.Fatal(err)
		}
	}
}
