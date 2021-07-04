package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/ermos/annotation"
	"github.com/ermos/annotation/parser"
	"github.com/ermos/dotenv"
	"github.com/mineway/excavator/api/routes"
	"github.com/mineway/excavator/api/server"
	"github.com/mineway/excavator/internal/pkg/config"
	"github.com/mineway/excavator/internal/pkg/logger"
	"github.com/mineway/excavator/internal/pkg/pipeline"
	"github.com/mineway/excavator/utils"
	"log"
	"os"
	"time"
)

//go:embed config.json
var configBytes []byte

//go:embed routes.json
var routesByte []byte

func main()  {
	ctx := context.Background()

	// Build Mode
	if utils.InArrayString(os.Args, "build") {
		var annotationResult []parser.API

		err := annotation.Fetch("api/routes", &annotationResult, parser.ToAPI)
		if err != nil {
			log.Fatal(err)
		}

		err = annotation.Save(annotationResult, "cmd/excavator/routes.json")
		if err != nil {
			log.Fatal(err)
		}

		return
	}

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

	// Start API
	var routesAPI []parser.API

	if err = json.Unmarshal(routesByte, &routesAPI); err != nil {
		log.Fatal(err)
	}

	server.Serve(c.ApiChan, "80", routesAPI, routes.Handler{})
	
	for {
		//fmt.Println("test")
		time.Sleep(1 * time.Second)
	}
}
