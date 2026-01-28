package main

import (
	"github.com/gin-gonic/gin"
	"github.com/laveshdeshwal/go-notify/internal/storage"
	"github.com/spf13/viper"
	"log"
	//"os"
)

func main() {
	// Load Environment Variables
	//env := os.Getenv("APP_ENV")
	envPath := ".env"

	viper.SetConfigFile(envPath)
	viper.ReadInConfig()
	viper.AutomaticEnv()

	app := gin.Default()

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Server": "go-notify-server",
			"Status": "running",
		})
	})

	if err := storage.Init(); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	database := storage.GetDB()

	var one int
	if err := database.QueryRow("SELECT 1").Scan(&one); err != nil {
		log.Fatal(err)
	}

	log.Println("DB query OK:", one)

	app.Run(":8080")
}
