package main

import (
	"fmt"
)

func CmdMetadata(config *Config, args []string) {
	registry := NewRegistry(config.IndexEndpoint)
}

func CmdCurlme(*Config, []string) {
}

func main() {
	OptParse()
}
