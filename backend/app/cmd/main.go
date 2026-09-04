package main

import (
	"basic-app/backend/app/config"
	"basic-app/backend/app/database"
	"basic-app/backend/app/router"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		//log.Fatalf("Config Error")
		log.Fatalf("Config Error: %v", err)
	}

	client, db, err := database.Connect(cfg)
	if err != nil {
		//log.Fatalf("db error")
		log.Fatalf("DB Error: %v", err)
	}

	defer func() {
		if err := database.Disconnect(client); err != nil {
			log.Printf("mongo disconnect error: %v", err)
		}
	}()

	gin.SetMode(cfg.GinMode)

	// middleware.StartCleanup() // Enable when rate limiting is introduced (v1.1)

	engine := router.NewRouter(db, cfg)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)

	log.Printf("Server listening on http://localhost%s", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("Server Failed")
	}

}
