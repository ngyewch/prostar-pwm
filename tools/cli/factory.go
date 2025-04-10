package main

import (
	prostar_pwm "github.com/ngyewch/prostar-pwm"
	"github.com/urfave/cli/v2"
	"sync"
)

func newDev(cCtx *cli.Context) (*prostar_pwm.Dev, error) {
	client, err := newModbusClient(cCtx, nil)
	if err != nil {
		return nil, err
	}

	modbusUnitId := modbusUnitIdFlag.Get(cCtx)

	err = client.Open()
	if err != nil {
		return nil, err
	}

	var mutex sync.Mutex

	dev := prostar_pwm.New(client, uint8(modbusUnitId), &mutex)

	return dev, nil
}
