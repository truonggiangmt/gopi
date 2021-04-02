package main

import (
	"os"

	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

func main() {
	os.Exit(tool.Server("chromecast", os.Args[1:], new(app)))
}
