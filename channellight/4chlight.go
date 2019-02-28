package channellight

import (
	"github.com/brutella/hc/service"
	"math"
)

type FourChannelColor struct {
	baseChannelColor
	Red   int
	Blue  int
	Green int
}

func (this *FourChannelColor) SetColor(color HSLColor) {
	this.baseChannelColor.SetColor(color)
	this.calculateRed(color)
	this.calculateGreen(color)
	this.calculateBlue(color)
}

type FourChannelLight struct {
	baseChannelLight
	outputColor FourChannelColor
}

func (this *FourChannelLight) SetColor(lightbulb *service.Lightbulb) {
	if this.baseChannelLight.SetColor(lightbulb) {
		newOutputColor := FourChannelColor{}
		newOutputColor.baseChannelColor.SetColor(this.targetColor)
		newOutputColor.calculateRed(this.targetColor)
		newOutputColor.calculateGreen(this.targetColor)
		newOutputColor.calculateBlue(this.targetColor)
		this.outputColor = newOutputColor
	}
}

func (this *FourChannelLight) GetOutputColor() FourChannelColor {
	return this.outputColor
}

func (this *FourChannelColor) scaleSat(saturation float64, base int) int {
	baseFloat := float64(base)
	floatVal := (100-saturation)*2.55 + saturation*baseFloat/100
	return int(math.Floor(floatVal + 0.5))
}

func (this *FourChannelColor) calculateRed(color HSLColor) {
	this.Red = this.scaleSat(color.Saturation, this.calculateBaseRed(color))
}

func (this *FourChannelColor) calculateGreen(color HSLColor) {
	this.Green = this.scaleSat(color.Saturation, this.calculateBaseGreen(color))
}

func (this *FourChannelColor) calculateBlue(color HSLColor) {
	this.Blue = this.scaleSat(color.Saturation, this.calculateBaseBlue(color))
}


func (this *FourChannelColor) calculateBaseRed(color HSLColor) int {
	return multiStop(color.Hue, 120, 240, 300, 60)
}

func (this *FourChannelColor) calculateBaseGreen(color HSLColor) int {
	return multiStop(color.Hue, 240, 0, 90, 180)
}

func (this *FourChannelColor) calculateBaseBlue(color HSLColor) int {
	return multiStop(color.Hue, 0, 120, 180, 300)
}
