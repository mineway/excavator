package pipes

import (
	"context"
	"fmt"
	"github.com/jaypipes/ghw"
	"github.com/mineway/excavator/internal/pkg/logger"
	"github.com/mineway/excavator/internal/pkg/rig"
)

var pipeName = "computer"

type Computer struct {}

func (Computer) GetName() string {
	return pipeName
}

func (Computer) Run(ctx context.Context, r *rig.Core) (err error) {
	logger.Info("[%s] extract computer informations..", pipeName)

	cpu, err := ghw.CPU()
	if err != nil {
		return
	}

	if len(cpu.Processors) != 0 {
		r.SetCPU(cpu)
		for _, processor := range cpu.Processors {
			logger.Success("[%s] CPU found : %s", pipeName, processor.Model)
		}
	} else {
		return fmt.Errorf("no one processor found in this computer")
	}

	gpu, err := ghw.GPU()
	if err != nil {
		return
	}

	if len(gpu.GraphicsCards) != 0 {
		r.SetGPU(gpu)
		for _, graphicsCard := range gpu.GraphicsCards {
			logger.Success("[%s] GPU found : %s", pipeName, graphicsCard.DeviceInfo.Product.Name)
		}
	} else {
		logger.Warning("[%s] no one GPU found in this computer", pipeName)
	}

	return
}