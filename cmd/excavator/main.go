package main

import (
	"context"
	"github.com/mineway/excavator/internal/pkg/logger"
	"github.com/mineway/excavator/internal/pkg/pipeline"
	"github.com/mineway/excavator/internal/pkg/rig"
	"time"
)

func main()  {
	ctx := context.Background()

	r, err := rig.New()
	if err != nil {
		logger.Fatal("rig instance failed: %s", err.Error())
	}

	err = pipeline.Run(ctx, r)
	if err != nil {
		logger.Fatal("pipeline failed: %s", err.Error())
	}

	for {
		//fmt.Println("test")
		time.Sleep(1 * time.Second)
	}
}
