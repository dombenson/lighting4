package channellight

import (
	"github.com/brutella/hc/service"
	"math"
)

type FiveChannelColor struct {
	FourChannelColor
	White int
}

func (this *FiveChannelColor) SetColor(color HSLColor) {
	this.FourChannelColor.SetColor(color)
	this.calculateRed(color)
	this.calculateGreen(color)
	this.calculateBlue(color)
	this.calculateWhite(color)
}

type FiveChannelLight struct {
	baseChannelLight
	outputColor FiveChannelColor
}

func (this *FiveChannelLight) SetColor(lightbulb *service.Lightbulb) {
	if this.baseChannelLight.SetColor(lightbulb) {
		newOutputColor := FiveChannelColor{}
		newOutputColor.SetColor(this.targetColor)
		this.outputColor = newOutputColor
	}
}

func (this *FiveChannelLight) GetOutputColor() FiveChannelColor {
	return this.outputColor
}

func (this *FiveChannelColor) calculateRed(color HSLColor) {
	this.Red = this.scaleSat(color.Saturation, this.calculateBaseRed(color))
}

func (this *FiveChannelColor) calculateGreen(color HSLColor) {
	this.Green = this.scaleSat(color.Saturation, this.calculateBaseGreen(color))
}

func (this *FiveChannelColor) calculateBlue(color HSLColor) {
	this.Blue = this.scaleSat(color.Saturation, this.calculateBaseBlue(color))
}

func (this *FiveChannelColor) calculateWhite(color HSLColor) {
	if color.Saturation < 50 {
		this.White = 255
		return
	}
	this.White = linearScale(color.Saturation, 100, 50)
}


func (this *FiveChannelColor) scaleSat(saturation float64, baseVal int) int {
	if saturation > 50 {
		return baseVal
	}
	return int(math.Floor(saturation*float64(baseVal)/50 + 0.5))
}
