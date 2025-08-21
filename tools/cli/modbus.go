package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/urfave/cli/v3"
)

func parseParity(s string) (uint, error) {
	s = strings.ToUpper(s)
	switch s {
	case "N", "NONE":
		return modbus.PARITY_NONE, nil
	case "E", "EVEN":
		return modbus.PARITY_EVEN, nil
	case "O", "ODD":
		return modbus.PARITY_ODD, nil
	}
	return 0, fmt.Errorf("invalid parity: %s", s)
}

func newModbusClient(cmd *cli.Command, configurer func(cfg *modbus.ClientConfiguration)) (*modbus.ModbusClient, error) {
	serialPort := cmd.String(serialPortFlag.Name)
	baudRate := cmd.Uint(baudRateFlag.Name)
	dataBits := cmd.Uint(dataBitsFlag.Name)
	parityString := cmd.String(parityFlag.Name)
	stopBits := cmd.Uint(stopBitsFlag.Name)

	parity, err := parseParity(parityString)
	if err != nil {
		return nil, err
	}

	config := &modbus.ClientConfiguration{
		URL:      "rtu://" + serialPort,
		Speed:    baudRate,
		DataBits: dataBits,
		Parity:   parity,
		StopBits: stopBits,
		Timeout:  1 * time.Second,
	}
	if configurer != nil {
		configurer(config)
	}

	client, err := modbus.NewClient(config)
	if err != nil {
		return nil, err
	}

	err = client.Open()
	if err != nil {
		return nil, err
	}

	return client, nil
}
