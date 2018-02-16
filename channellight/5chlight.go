package channellight

import "github.com/brutella/hc/service"

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
	this.baseChannelLight.SetColor(lightbulb)
	newOutputColor := FiveChannelColor{}
	newOutputColor.SetColor(this.targetColor)
	this.outputColor = newOutputColor
}

func (this *FiveChannelLight) GetOutputColor() FiveChannelColor {
	return this.outputColor
}

func (this *FiveChannelColor) calculateRed(color HSLColor) {
	this.Red = scaleSat(color.Saturation, this.calculateBaseRed(color))
}

func (this *FiveChannelColor) calculateGreen(color HSLColor) {
	this.Green = scaleSat(color.Saturation, this.calculateBaseGreen(color))
}

func (this *FiveChannelColor) calculateBlue(color HSLColor) {
	this.Blue = scaleSat(color.Saturation, this.calculateBaseBlue(color))
}

func (this *FiveChannelColor) calculateWhite(color HSLColor) {
	if color.Saturation < 50 {
		this.White = 255
		return
	}
	this.White = linearScale(color.Saturation, 100, 50)
}
