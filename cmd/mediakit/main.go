package main

import (
	"os"

	"github.com/djthorpe/gopi/v3/pkg/tool"
)

func main() {
	os.Exit(tool.CommandLine("mediakit", os.Args[1:], new(app)))
}
