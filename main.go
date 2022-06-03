package main

import (
	"github.com/home-sol/homectl/cmd"
	"github.com/home-sol/homectl/pkg/utils"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		utils.PrintErrorToStdErrorAndExit(err)
	}
}
