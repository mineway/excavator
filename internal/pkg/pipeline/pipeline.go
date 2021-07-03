package pipeline

import (
	"context"
	"github.com/mineway/excavator/internal/pkg/config"
	"github.com/mineway/excavator/internal/pkg/pipeline/pipes"
)

type Piper interface {
	GetName() string
	Run(ctx context.Context, c *config.Config) error
}

func Run (ctx context.Context, c *config.Config) error {
	for _, pipe := range pipelines {
		if err := pipe.Run(ctx, c); err != nil {
			return err
		}
	}
	return nil
}

// Pipelines
var pipelines = []Piper{
	pipes.Computer{},
}
