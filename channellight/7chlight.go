package channellight

import "github.com/brutella/hc/service"

type SevenChannelColor struct {
	FourChannelColor
	Amber int
	Uv    int
	White int
}

func (this *SevenChannelColor) SetColor(color HSLColor) {
	this.FourChannelColor.SetColor(color)
	// Red and Green have different ranges to accommodate Amber
	this.calculateRed(color)
	this.calculateGreen(color)
	this.calculateBlue(color)
	this.calculateWhite(color)
	this.calculateAmber(color)
	this.calculateUv(color)
}

type SevenChannelLight struct {
	baseChannelLight
	outputColor SevenChannelColor
}

func (this *SevenChannelLight) SetColor(lightbulb *service.Lightbulb) {
	this.baseChannelLight.SetColor(lightbulb)
	newOutputColor := SevenChannelColor{}
	newOutputColor.SetColor(this.targetColor)
	this.outputColor = newOutputColor
}

func (this *SevenChannelLight) GetOutputColor() SevenChannelColor {
	return this.outputColor
}

func (this *SevenChannelColor) calculateRed(color HSLColor) {
	this.Red = scaleSat(color.Saturation, multiStop(color.Hue, 30, 240, 300, 15))
}

func (this *SevenChannelColor) calculateGreen(color HSLColor) {
	this.Green = scaleSat(color.Saturation, multiStop(color.Hue, 240, 30, 90, 180))
}

func (this *SevenChannelColor) calculateBlue(color HSLColor) {
	this.Blue = scaleSat(color.Saturation, multiStop(color.Hue, 0, 120, 180, 300))
}

func (this *SevenChannelColor) calculateWhite(color HSLColor) {
	if color.Saturation < 50 {
		this.White = 255
		return
	}
	this.White = linearScale(color.Saturation, 100, 50)
}

func (this *SevenChannelColor) calculateAmber(color HSLColor) {
	this.Amber = scaleSat(color.Saturation, multiStop(color.Hue, 120, 0, 15, 60))
}

func (this *SevenChannelColor) calculateUv(color HSLColor) {
	this.Uv = scaleSat(color.Saturation, multiStop(color.Hue, 350, 240, 290, 310))
}
