package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/ermos/dotenv"
	"github.com/mineway/excavator/api/middleware"
	"github.com/mineway/excavator/api/routes"
	"github.com/mineway/excavator/api/server"
	"github.com/mineway/excavator/internal/pkg/config"
	"github.com/mineway/excavator/internal/pkg/logger"
	"github.com/mineway/excavator/internal/pkg/pipeline"
	"log"
	"time"
)

//go:embed config.json
var configBytes []byte

func main()  {
	ctx := context.Background()

	// Load .env if exist
	_ = dotenv.Parse(".env")

	// Init Config
	c, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Import Embed Configuration
	if err := json.Unmarshal(configBytes, &c); err != nil {
		log.Fatal(err)
	}

	// Import User Current Setting
	if err := c.Init(); err != nil {
		log.Fatal(err)
	}

	// Start pipeline process
	err = pipeline.Run(ctx, c)
	if err != nil {
		logger.Fatal("pipeline failed: %s", err.Error())
	}

	server.Serve(c.ApiChan, "80", "dist/t.json", routes.Handler{}, middleware.Handler{})
	
	for {
		//fmt.Println("test")
		time.Sleep(1 * time.Second)
	}
}
