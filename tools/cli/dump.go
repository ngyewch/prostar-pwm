package main

import (
	"github.com/yassinebenaid/godump"
	"os"
)

var (
	dumper godump.Dumper
)

func init() {
	noColor := os.Getenv("NO_COLOR") == "1"
	dumpTheme := godump.DefaultTheme
	if noColor {
		dumpTheme = godump.Theme{}
	}
	dumpTheme.Address = &suppressedStyle{}

	dumper = godump.Dumper{
		Indentation:       "  ",
		HidePrivateFields: true,
		Theme:             dumpTheme,
	}
}

type suppressedStyle struct {
}

func (style suppressedStyle) Apply(s string) string {
	return ""
}

func dump(v any) error {
	return doDump(v)
}

func doDump(v any) error {
	return dumper.Println(v)
}
