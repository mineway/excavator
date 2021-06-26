package pipeline

import (
	"context"
	"github.com/mineway/excavator/internal/pkg/pipeline/pipes"
	"github.com/mineway/excavator/internal/pkg/rig"
)

type Piper interface {
	GetName() string
	Run(ctx context.Context, r *rig.Core) error
}

func Run (ctx context.Context, r *rig.Core) error {
	for _, pipe := range pipelines {
		if err := pipe.Run(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// Pipelines
var pipelines = []Piper{
	pipes.Computer{},
}
