package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"runtime/debug"
)

var (
	serialPortFlag = &cli.StringFlag{
		Name:     "serial-port",
		Usage:    "serial port",
		Required: true,
		Sources:  cli.EnvVars("SERIAL_PORT"),
		Category: "Serial",
	}
	baudRateFlag = &cli.UintFlag{
		Name:     "baud-rate",
		Usage:    "baud rate",
		Value:    9600,
		Sources:  cli.EnvVars("BAUD_RATE"),
		Category: "Serial",
	}
	dataBitsFlag = &cli.UintFlag{
		Name:     "data-bits",
		Usage:    "data bits",
		Value:    8,
		Sources:  cli.EnvVars("DATA_BITS"),
		Category: "Serial",
	}
	parityFlag = &cli.StringFlag{
		Name:    "parity",
		Usage:   "parity",
		Value:   "N",
		Sources: cli.EnvVars("PARITY"),
		Action: func(ctx context.Context, cmd *cli.Command, s string) error {
			_, err := parseParity(s)
			return err
		},
		Category: "Serial",
	}
	stopBitsFlag = &cli.UintFlag{
		Name:     "stop-bits",
		Usage:    "stop bits",
		Value:    2,
		Sources:  cli.EnvVars("STOP_BITS"),
		Category: "Serial",
	}
	modbusUnitIdFlag = &cli.UintFlag{
		Name:    "modbus-unit-id",
		Usage:   "ModBus unit ID",
		Value:   1,
		Sources: cli.EnvVars("MODBUS_UNIT_ID"),
		Action: func(ctx context.Context, cmd *cli.Command, v uint) error {
			if (v < 1) || (v > 247) {
				return fmt.Errorf("invalid modbus-unit-id: %d", v)
			}
			return nil
		},
		Category: "Modbus",
	}

	app = &cli.Command{
		Name:  "prostar-pwm",
		Usage: "ProStar PWM CLI",
		Commands: []*cli.Command{
			{
				Name:   "raw-adc-data",
				Usage:  "raw ADC data",
				Action: doRawADCData,
			},
			{
				Name:   "filtered-adc-data",
				Usage:  "filtered ADC data",
				Action: doFilteredADCData,
			},
			{
				Name:   "temperature-data",
				Usage:  "temperature data",
				Action: doTemperatureData,
			},
			{
				Name:   "charger-status",
				Usage:  "charger status",
				Action: doChargerStatus,
			},
			{
				Name:   "load-status",
				Usage:  "load status",
				Action: doLoadStatus,
			},
			{
				Name:   "misc-data",
				Usage:  "misc data",
				Action: doMiscData,
			},
			{
				Name:   "charge-settings",
				Usage:  "charge settings",
				Action: doChargeSettings,
			},
			{
				Name:   "load-settings",
				Usage:  "load settings",
				Action: doLoadSettings,
			},
			{
				Name:   "misc-settings",
				Usage:  "misc settings",
				Action: doMiscSettings,
			},
			{
				Name:   "pwm-settings",
				Usage:  "pwm settings",
				Action: doPWMSettings,
			},
			{
				Name:   "statistics",
				Usage:  "statistics",
				Action: doStatistics,
			},
			{
				Name:   "logged-data",
				Usage:  "logged data",
				Action: doLoggedData,
			},
		},
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

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
