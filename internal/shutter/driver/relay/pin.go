package relay

import (
	"context"
	"github.com/racerxdl/go-mcp23017"
	"time"
)

type Mcp23017Pin struct {
	device *mcp23017.Device
	pin    uint8
}

func NewMcp23017Pin(device *mcp23017.Device, pin uint8) (p *Mcp23017Pin, err error) {
	p = &Mcp23017Pin{}
	p.device = device
	p.pin = pin
	err = p.device.PinMode(pin, mcp23017.OUTPUT)
	return p, err
}

func (m *Mcp23017Pin) High() error {
	return m.device.DigitalWrite(m.pin, mcp23017.HIGH)
}

func (m *Mcp23017Pin) Low() error {
	return m.device.DigitalWrite(m.pin, mcp23017.LOW)
}

type SetPin interface {
	High() error
	Low() error
}

type Wired struct {
	Pin         SetPin
	NormalClose bool
}

func (p *Wired) EnableFor(ctx context.Context, duration time.Duration) error {
	after := time.After(duration)
	if err := p.enable(); err != nil {
		return err
	}
	defer p.disable()

	for {
		select {
		case <-after:
		case <-ctx.Done():
			return nil
		}

		time.Sleep(time.Millisecond)
	}
}

func (p *Wired) enable() error {
	if !p.NormalClose {
		return p.Pin.High()
	}

	return p.Pin.Low()
}

func (p *Wired) disable() error {
	if !p.NormalClose {
		return p.Pin.Low()
	}

	return p.Pin.High()
}