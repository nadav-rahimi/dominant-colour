package main

import (
	"github.com/fiwippi/go-quantise/pkg/quantisers/lmq"
	"github.com/fiwippi/go-quantise/pkg/quantisers/otsu"
	"github.com/fiwippi/go-quantise/pkg/quantisers/pnn"
	"github.com/fiwippi/go-quantise/pkg/quantisers/pnnlab"
	"testing"
)

var benchImg, _ = ReadImage("fish.jpg")

func BenchmarkOtsuGreySingle(b *testing.B) {
	otsu.QuantiseGreyscale(benchImg)
}

func BenchmarkLMQGreySingle(b *testing.B) {
	lmq.QuantiseGreyscale(benchImg, 1)
}

func BenchmarkLMQGreyMulti(b *testing.B) {
	lmq.QuantiseGreyscale(benchImg, 6)
}

func BenchmarkPNNGreySingle(b *testing.B) {
	pnn.QuantiseGreyscale(benchImg, 1)
}

func BenchmarkPNNGreyMulti(b *testing.B) {
	pnn.QuantiseGreyscale(benchImg, 6)
}

func BenchmarkPNNColourMulti(b *testing.B) {
	pnn.QuantiseColour(benchImg, 6)
}

func BenchmarkPNNLABColourMulti(b *testing.B) {
	pnnlab.QuantiseColour(benchImg, 6)
}
