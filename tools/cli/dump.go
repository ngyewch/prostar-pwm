package main

import (
	"github.com/yassinebenaid/godump"
)

var (
	dumper godump.Dumper
)

func init() {
	dumpTheme := godump.DefaultTheme
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
