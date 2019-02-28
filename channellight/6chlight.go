package channellight

import (
	"github.com/brutella/hc/service"
)

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

func (this *SixChannelColor) GetColor() (color HSLColor) {
	normalisedColor := *this
	minRgb := normalisedColor.Blue
	if normalisedColor.Red < minRgb {
		minRgb = normalisedColor.Red
	}
	if normalisedColor.Green < minRgb {
		minRgb = normalisedColor.Green
	}
	if minRgb > 0 {
		normalisedColor.White = normalisedColor.White + minRgb
		normalisedColor.Red = normalisedColor.Red - minRgb
		normalisedColor.Green = normalisedColor.Green - minRgb
		normalisedColor.Blue = normalisedColor.Blue - minRgb
	}
	if normalisedColor.Amber > 0 && normalisedColor.Blue > 0 {
		minRgb = normalisedColor.Amber
		if normalisedColor.Blue < minRgb {
			minRgb = normalisedColor.Blue
		}
		normalisedColor.Amber = normalisedColor.Amber - minRgb
		normalisedColor.Blue = normalisedColor.Blue - minRgb
		normalisedColor.White = normalisedColor.White + minRgb
	}
	if normalisedColor.Red > 0 && normalisedColor.Green > 0 {
		minRgb = normalisedColor.Red
		if normalisedColor.Green < minRgb {
			minRgb = normalisedColor.Green
		}
		normalisedColor.Red = normalisedColor.Red - minRgb
		normalisedColor.Green = normalisedColor.Green - minRgb
		normalisedColor.Amber = normalisedColor.Amber + minRgb
	}

	maxChannel := 0
	maxColorChannel := 0
	if normalisedColor.Red > maxColorChannel {
		maxColorChannel = normalisedColor.Red
	}
	if normalisedColor.Amber > maxColorChannel {
		maxColorChannel = normalisedColor.Amber
	}
	if normalisedColor.Green > maxColorChannel {
		maxColorChannel = normalisedColor.Green
	}
	if normalisedColor.Blue > maxColorChannel {
		maxColorChannel = normalisedColor.Blue
	}
	maxChannel = maxColorChannel
	if normalisedColor.White > maxChannel {
		maxChannel = normalisedColor.White
	}
	if normalisedColor.White > 0 {
		if maxColorChannel > 0 {

		} else {
			color.Saturation = 0
		}
	} else {
		color.Saturation = 255
	}


	if maxChannel != 255 {
		scaling := float64(255)/float64(maxChannel)
		normalisedColor.White = int(scaling*float64(normalisedColor.White))
		normalisedColor.Red = int(scaling*float64(normalisedColor.White))
		normalisedColor.Green = int(scaling*float64(normalisedColor.White))
		normalisedColor.Blue = int(scaling*float64(normalisedColor.White))
		normalisedColor.Amber = int(scaling*float64(normalisedColor.White))
		normalisedColor.Brightness = int(float64(normalisedColor.Brightness)/scaling)
	}

	if normalisedColor.Brightness > 255 {
		normalisedColor.Brightness = 255
	}

	color.Brightness = normalisedColor.Brightness


	if normalisedColor.Red > 0 {
		if normalisedColor.Amber > 0 {

		} else if normalisedColor.Blue > 0 {

		} else {
			color.Hue = 0
		}
	} else if normalisedColor.Amber > 0 {
		if normalisedColor.Green > 0 {

		} else {
			color.Hue = 30
		}
	} else if normalisedColor.Green > 0 {
		if normalisedColor.Blue > 0 {

		} else {
			color.Hue = 120
		}
	} else if normalisedColor.Blue > 0 {
		color.Hue = 240
	} else {
		color.Hue = 0
		color.Saturation = 0
	}


	return
}

type SixChannelLight struct {
	baseChannelLight
	outputColor SixChannelColor
}

func (this *SixChannelLight) SetColor(lightbulb *service.Lightbulb) {
	if this.baseChannelLight.SetColor(lightbulb) {
		newOutputColor := SixChannelColor{}
		newOutputColor.SetColor(this.targetColor)
		this.outputColor = newOutputColor
	}
}

func (this *SixChannelLight) GetOutputColor() SixChannelColor {
	return this.outputColor
}

func (this *SixChannelColor) calculateRed(color HSLColor) {
	this.Red = this.scaleSat(color.Saturation, multiStop(color.Hue, 30, 240, 300, 15))
}

func (this *SixChannelColor) calculateGreen(color HSLColor) {
	this.Green = this.scaleSat(color.Saturation, multiStop(color.Hue, 240, 30, 90, 180))
}

func (this *SixChannelColor) calculateAmber(color HSLColor) {
	this.Amber = this.scaleSat(color.Saturation, multiStop(color.Hue, 120, 0, 15, 70))
}

