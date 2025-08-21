package main

import (
	"sync"

	prostar_pwm "github.com/ngyewch/prostar-pwm"
	"github.com/urfave/cli/v3"
)

func newDev(cmd *cli.Command) (*prostar_pwm.Dev, error) {
	client, err := newModbusClient(cmd, nil)
	if err != nil {
		return nil, err
	}

	modbusUnitId := cmd.Uint(modbusUnitIdFlag.Name)

	err = client.Open()
	if err != nil {
		return nil, err
	}

	var mutex sync.Mutex

	dev := prostar_pwm.New(client, uint8(modbusUnitId), &mutex)

	return dev, nil
}
