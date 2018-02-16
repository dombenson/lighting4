package channellight

import (
	"github.com/brutella/hc/service"
)

type ChannelLight interface {
	SetColor(lightbulb *service.Lightbulb)
}

type HSLColor struct {
	Hue        float64
	Saturation float64
	Brightness int
}
