package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
)

func main() {
	log.SetFlags(log.LstdFlags)
	flag.Parse()
	{
		if len(*typeNames) == 0 {
			flag.Usage()
			os.Exit(2)
		}
		_type := strings.Split(*typeNames, ",")
		log.Printf("type:%v", _type)
	}

	{
		args := flag.Args()
		if len(args) == 0 {
			args = []string{"."}
		}
		var file string
		var dir string
		if len(args) == 1 && func(name string) bool {
			info, err := os.Stat(name)
			if err != nil {
				log.Fatal(err)
			}
			return info.IsDir()
		}(args[0]) {
			dir = args[0]
		} else {
			file = args[0]
			dir = filepath.Dir(args[0])
		}
		log.Printf("dir:%v,file:%v", dir, file)
	}
}
