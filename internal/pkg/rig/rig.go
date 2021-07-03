package rig

import (
	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/gpu"
	"net"
)

type Data struct {
	UdpAddr 	*net.UDPAddr
	Cpu 		*cpu.Info
	Gpu 		*gpu.Info
}

func New() (d *Data, err error) {
	d = new(Data)

	if err = d.setOutboundIpAddr(); err != nil {
		return
	}

	return
}

func (d *Data) setOutboundIpAddr() (err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err
	}
	defer conn.Close()

	d.UdpAddr = conn.LocalAddr().(*net.UDPAddr)

	return
}