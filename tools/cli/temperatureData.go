package main

import "github.com/urfave/cli/v2"

func doTemperatureData(cCtx *cli.Context) error {
	dev, err := newDev(cCtx)
	if err != nil {
		return err
	}

	result, err := dev.ReadTemperatureData()
	if err != nil {
		return err
	}

	err = dump(result)
	if err != nil {
		return err
	}

	return nil
}
