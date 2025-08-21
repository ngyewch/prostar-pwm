package main

import (
	"context"

	"github.com/urfave/cli/v3"
)

func doChargeSettings(ctx context.Context, cmd *cli.Command) error {
	dev, err := newDev(cmd)
	if err != nil {
		return err
	}

	result, err := dev.ReadChargeSettings()
	if err != nil {
		return err
	}

	err = dump(result)
	if err != nil {
		return err
	}

	return nil
}
