package channellight

import "github.com/brutella/hc/service"

type SixChannelColor struct {
	FiveChannelColor
	Amber int
}

func (this *SixChannelColor) SetColor(color HSLColor) {
	this.FiveChannelColor.SetColor(color)
	// Red and Green have different ranges to accommodate Amber
	this.calculateRed(color)
	this.calculateGreen(color)
	this.calculateAmber(color)
}

type SixChannelLight struct {
	baseChannelLight
	outputColor SixChannelColor
}

func (this *SixChannelLight) SetColor(lightbulb *service.Lightbulb) {
	this.baseChannelLight.SetColor(lightbulb)
	newOutputColor := SixChannelColor{}
	newOutputColor.SetColor(this.targetColor)
	this.outputColor = newOutputColor
}

func (this *SixChannelLight) GetOutputColor() SixChannelColor {
	return this.outputColor
}

func (this *SixChannelColor) calculateRed(color HSLColor) {
	this.Red = scaleSat(color.Saturation, multiStop(color.Hue, 30, 240, 300, 15))
}

func (this *SixChannelColor) calculateGreen(color HSLColor) {
	this.Green = scaleSat(color.Saturation, multiStop(color.Hue, 240, 30, 90, 180))
}

func (this *SixChannelColor) calculateAmber(color HSLColor) {
	this.Amber = scaleSat(color.Saturation, multiStop(color.Hue, 120, 0, 15, 70))
}
