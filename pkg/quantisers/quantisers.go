package quantisers

import (
	LMQ "github.com/nadav-rahimi/dominant-colour/internal/lmq"
	Otsu "github.com/nadav-rahimi/dominant-colour/internal/otsu"
	PNN "github.com/nadav-rahimi/dominant-colour/internal/pnn"
	PNNLAB "github.com/nadav-rahimi/dominant-colour/internal/pnnlab"
	"image"
	"image/color"
)

// Quantiser interface used to define common quantiser behaviour
type Quantiser interface {
	Greyscale(img image.Image, m int) (color.Palette, error)
	Colour(img image.Image, m int) (color.Palette, error)
}

type OtsuQuantiser struct {
	*Otsu.Otsu
}
type LMQQuantiser struct {
	*LMQ.LMQ
}
type PNNQuantiser struct {
	*PNN.PNN
}
type PNNLABQuantiser struct {
	*PNNLAB.PNNLAB
}

func NewOtsuQuantiser() *OtsuQuantiser {
	return &OtsuQuantiser{&Otsu.Otsu{}}
}
func NewLMQQuantiser() *LMQQuantiser {
	return &LMQQuantiser{&LMQ.LMQ{}}
}
func NewPNNQuantiser() *PNNQuantiser {
	return &PNNQuantiser{&PNN.PNN{}}
}
func NewPNNLABQuantiser() *PNNLABQuantiser {
	return &PNNLABQuantiser{&PNNLAB.PNNLAB{}}
}
