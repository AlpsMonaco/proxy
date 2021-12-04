package vpn

import (
	"math/rand"
	"time"
)

var globalRandomNumberGenerator RandomNumberGenerator

func init() {
	var r = &BasicRandomNumberGenerator{}
	r.Init()
	globalRandomNumberGenerator = r
}

type BasicRandomNumberGenerator struct{}

func (brng *BasicRandomNumberGenerator) Init() {
	rand.Seed(time.Now().Unix())
}

func (brng *BasicRandomNumberGenerator) Read(b []byte) {
	for k := range b {
		b[k] = byte(rand.Intn(0xFF))
	}
}

type RandomNumberGenerator interface {
	Read([]byte)
}
