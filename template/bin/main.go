package main

import (
	"fmt"
	"github.com/sqjian/go-kit/splash"
	"github.com/sqjian/go-kit/template/bin/cmd"
)

func init() {
	fmt.Println(splash.Stringify())
}
func main() {
	cmd.Execute()
}
