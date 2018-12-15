// +build ignore

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/cmd"
	"github.com/99designs/gqlgen/codegen"
)

func main() {
	config, err := codegen.LoadConfigFromDefaultLocations()
	if err != nil {
		log.Fatalln("cannot find gqlgen.yml anywhere from current working directory:", err)
	}

	p, err := filepath.Abs(config.Resolver.Filename)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Remove(p)
	if err != nil {
		log.Fatalln(err)
	}
	cmd.Execute()
}
