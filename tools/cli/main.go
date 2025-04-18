package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime/debug"
)

var (
	serialPortFlag = &cli.StringFlag{
		Name:     "serial-port",
		Usage:    "serial port",
		Required: true,
		EnvVars:  []string{"SERIAL_PORT"},
		Category: "Serial",
	}
	baudRateFlag = &cli.UintFlag{
		Name:     "baud-rate",
		Usage:    "baud rate",
		Value:    9600,
		EnvVars:  []string{"BAUD_RATE"},
		Category: "Serial",
	}
	dataBitsFlag = &cli.UintFlag{
		Name:     "data-bits",
		Usage:    "data bits",
		Value:    8,
		EnvVars:  []string{"DATA_BITS"},
		Category: "Serial",
	}
	parityFlag = &cli.StringFlag{
		Name:    "parity",
		Usage:   "parity",
		Value:   "N",
		EnvVars: []string{"PARITY"},
		Action: func(cCtx *cli.Context, s string) error {
			_, err := parseParity(s)
			return err
		},
		Category: "Serial",
	}
	stopBitsFlag = &cli.UintFlag{
		Name:     "stop-bits",
		Usage:    "stop bits",
		Value:    2,
		EnvVars:  []string{"STOP_BITS"},
		Category: "Serial",
	}
	modbusUnitIdFlag = &cli.UintFlag{
		Name:    "modbus-unit-id",
		Usage:   "ModBus unit ID",
		Value:   1,
		EnvVars: []string{"MODBUS_UNIT_ID"},
		Action: func(cCtx *cli.Context, v uint) error {
			if (v < 1) || (v > 247) {
				return fmt.Errorf("invalid modbus-unit-id: %d", v)
			}
			return nil
		},
		Category: "Modbus",
	}

	app = &cli.App{
		Name:  "prostar-pwm",
		Usage: "ProStar PWM CLI",
		Flags: []cli.Flag{
			serialPortFlag,
			baudRateFlag,
			dataBitsFlag,
			parityFlag,
			stopBitsFlag,
			modbusUnitIdFlag,
		},
	}
)

func main() {
	buildInfo, _ := debug.ReadBuildInfo()
	if buildInfo != nil {
		app.Version = buildInfo.Main.Version
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
