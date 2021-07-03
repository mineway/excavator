package pipes

import (
	"context"
	"fmt"
	"github.com/jaypipes/ghw"
	"github.com/mineway/excavator/internal/pkg/config"
	"github.com/mineway/excavator/internal/pkg/logger"
	"runtime"
)

var pipeName = "computer"

type Computer struct {}

func (Computer) GetName() string {
	return pipeName
}

func (Computer) Run(ctx context.Context, c *config.Config) (err error) {
	logger.Info("[%s] extract computer informations..", pipeName)

	c.OS = runtime.GOOS
	c.Arch = runtime.GOARCH

	logger.Success("[%s] OS found : %s (%s)", pipeName, c.OS, c.Arch)

	cpu, err := ghw.CPU()
	if err != nil {
		return
	}

	if len(cpu.Processors) != 0 {
		c.RigData.Cpu = cpu
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
		c.RigData.Gpu = gpu
		for _, graphicsCard := range gpu.GraphicsCards {
			logger.Success("[%s] GPU found : %s", pipeName, graphicsCard.DeviceInfo.Product.Name)
		}
	} else {
		logger.Warning("[%s] no one GPU found in this computer", pipeName)
	}

	return
}