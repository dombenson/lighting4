package channellight

import (
	"github.com/brutella/hc/service"
	"math"
)

type baseChannelLight struct {
	targetColor HSLColor
}

type baseChannelColor struct {
	Brightness int
}

func (this *baseChannelLight) SetColor(lightbulb *service.Lightbulb) bool {
	newColor := HSLColor{}
	newColor.Hue = lightbulb.Hue.GetValue()
	newColor.Saturation = lightbulb.Saturation.GetValue()
	newColor.Brightness = lightbulb.Brightness.GetValue()
	if newColor != this.targetColor {
		this.targetColor = newColor
		return true
	}
	return false
}

func (this *baseChannelLight) GetColor() HSLColor {
	return this.targetColor
}

func (this *baseChannelColor) calculateBright(color HSLColor) {
	this.Brightness = int(math.Floor(0.5 + 2.55*float64(color.Brightness)))
}

func (this *baseChannelColor) SetColor(color HSLColor) {
	this.calculateBright(color)
}
