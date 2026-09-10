package main

import (
	"basic-app/config"
	"basic-app/database"
	mongorepo "basic-app/repository/mongo"
	"basic-app/router"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	client, db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("DB Error: %v", err)
	}

	defer func() {
		if err := database.Disconnect(client); err != nil {
			log.Printf("mongo disconnect error: %v", err)
		}
	}()

	// Ensure MongoDB indexes exist
	if err := mongorepo.EnsureUserIndexes(
		context.Background(),
		db,
	); err != nil {
		log.Fatalf("Failed to create user indexes: %v", err)
	}

	gin.SetMode(cfg.GinMode)

	// middleware.StartCleanup()

	engine := router.NewRouter(db, cfg)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)

	log.Printf("Server listening on http://localhost%s", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("Server Failed: %v", err)
	}
}
