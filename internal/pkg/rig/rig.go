package rig

import (
	"fmt"
	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/gpu"
	"github.com/mineway/excavator/internal/pkg/miner"
	"net"
	"strings"
)

type Core struct {
	udpAddr 	*net.UDPAddr
	name 		string
	miners 		[]miner.Core
	cpu 		*cpu.Info
	gpu 		*gpu.Info
}

func New() (c *Core, err error) {
	c = new(Core)

	if err = c.setOutboundIpAddr(); err != nil {
		return
	}

	return
}

func (c *Core) GetName() string {
	return c.name
}

func (c *Core) SetName(name string) error {
	if len(strings.TrimSpace(name)) < 3 {
		return fmt.Errorf("rig's name need to be greater than three characters")
	}

	c.name = name

	return nil
}

func (c *Core) SetGPU(gpu *gpu.Info) {
	c.gpu = gpu
}

func (c *Core) SetCPU(cpu *cpu.Info) {
	c.cpu = cpu
}

func (c *Core) setOutboundIpAddr() (err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err
	}
	defer conn.Close()

	c.udpAddr = conn.LocalAddr().(*net.UDPAddr)

	return
}