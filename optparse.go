package main

import (
	"os"
	"fmt"
	"flag"
)

type Config struct {
	IndexDomain string
	RegistryDomain string
}

type command struct {
	name string
	argsdesc string
	desc string
	nargs int
	fn func (*Config, []string)
}

var commands = [...]command {
	command{"info", "<repos_name>", "lookup repos meta-data", 1, CmdInfo},
	command{"layerinfo", "<repos_name> <layer_id>", "lookup layer meta-data", 2, CmdLayerInfo},
	command{"curlme", "<repos_name> <layer_id>", "print a curl command for fetching the layer", 2, CmdCurlme},
}

func OptParse() {
	flag.Usage = func () {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] command\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "commands:")
		for _, c := range commands {
			fmt.Fprintf(os.Stderr, "  %s %s: %s\n", c.name, c.argsdesc, c.desc)
		}
	}

	config := &Config{}
	flag.StringVar(&config.RepositoryName, "repos", "", "repository name (\"foo/bar\")")
	flag.StringVar(&config.ImageId, "image", "", "image id to test (64 chars long)")
	flag.StringVar(&config.IndexDomain, "index", "https://index.docker.io", "override index endpoint")
	flag.StringVar(&config.RegistryDomain, "registry", "", "override registry endpoint")
	flag.Parse()

	for _, c := range commands {
		if flag.Arg(0) == c.name {
			args := flag.Args()[1:]
			if len(args) != c.nargs {
				fmt.Fprintf(os.Stderr, "%s takes %d arguments: %s\n", c.name, c.nargs, c.argsdesc)
				os.Exit(2)
			}
			c.fn(config, args)
			return
		}
	}
	flag.Usage()
	os.Exit(2)
}

