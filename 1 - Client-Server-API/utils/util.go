package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func ExitIfError(err error, title string) {
	if err != nil {
		c := color.New(color.BgRed)
		c.Println(title)
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func ShowError(err error, message string) {
	if err != nil {
		c := color.New(color.BgRed)
		c.Println(message)
		fmt.Println(err.Error())
	}
}
