package channellight

import "github.com/brutella/hc/service"

type SevenChannelColor struct {
	SixChannelColor
	Uv    int
}

func (this *SevenChannelColor) SetColor(color HSLColor) {
	this.SixChannelColor.SetColor(color)
	this.calculateUv(color)
}

type SevenChannelLight struct {
	baseChannelLight
	outputColor SevenChannelColor
}

func (this *SevenChannelLight) SetColor(lightbulb *service.Lightbulb) {
	if this.baseChannelLight.SetColor(lightbulb) {
		newOutputColor := SevenChannelColor{}
		newOutputColor.SetColor(this.targetColor)
		this.outputColor = newOutputColor
	}
}

func (this *SevenChannelLight) GetOutputColor() SevenChannelColor {
	return this.outputColor
}

func (this *SevenChannelColor) calculateUv(color HSLColor) {
	this.Uv = this.scaleSat(color.Saturation, multiStop(color.Hue, 350, 240, 290, 310))
}
