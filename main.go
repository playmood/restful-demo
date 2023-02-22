package main

import (
	"fmt"
	"github.com/playmood/restful-demo/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
