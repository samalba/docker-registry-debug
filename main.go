package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"os"
	"strings"
	"text/tabwriter"
)

var Quiet = false

func fatal(msg error) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}

func log(format string, args ...interface{}) {
	if Quiet == true {
		return
	}
	format = fmt.Sprintf("# %s\n", format)
	fmt.Fprintf(os.Stderr, format, args...)
}

func readUserCredentials() (string, string) {
	log("Reading user/passwd from env var \"USER_CREDS\"")
	v := os.Getenv("USER_CREDS")
	if v == "" {
		log("No password provided, disabling auth")
		return "", ""
	}
	user := strings.SplitN(v, ":", 2)
	if len(user) < 2 {
		return user[0], ""
	}
	return user[0], user[1]
}

func initRegistry(config *Config, reposName string) (*Registry, string) {
	registry, e := NewRegistry(config.IndexDomain, config.RegistryDomain)
	if e != nil {
		fatal(e)
	}
	registry.Logger = log
	user, passwd := readUserCredentials()
	log("Getting token from %s", config.IndexDomain)
	token, e := registry.GetToken(user, passwd, reposName)
	if e != nil {
		fatal(e)
	}
	return registry, token
}

func CmdInfo(config *Config, args []string) {
	registry, token := initRegistry(config, args[0])
	tags, err := registry.ReposTags(token, args[0])
	if err != nil {
		fatal(err)
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Println("- Repository:", args[0])
	fmt.Println("- Tags:")
	for k, v := range tags {
		fmt.Fprintf(w, "\t%s\t%s\n", k, v)
	}
	w.Flush()
}

func CmdLayerInfo(config *Config, args []string) {
	registry, token := initRegistry(config, args[0])
	info, err := registry.LayerJson(token, args[1])
	if err != nil {
		fatal(err)
	}
	ancestry, err := registry.LayerAncestry(token, args[1])
	if err != nil {
		fatal(err)
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintf(w, "- Id\t%s\n", info.Id)
	fmt.Fprintf(w, "- Parent\t%s\n", info.Parent)
	fmt.Fprintf(w, "- Size\t%s\n", humanize.Bytes(uint64(info.Size)))
	fmt.Fprintf(w, "- Created\t%s\n", info.Created)
	fmt.Fprintf(w, "- DockerVersion\t%s\n", info.DockerVersion)
	fmt.Fprintf(w, "- Author\t%s\n", info.Author)
	fmt.Fprintf(w, "- Ancestry:")
	for _, id := range *ancestry {
		fmt.Fprintf(w, "\t%s\n", id)
	}
	w.Flush()
}

func CmdCurlme(config *Config, args []string) {
	registry, token := initRegistry(config, args[0])
	fmt.Printf("curl -i --location-trusted -I -X GET -H \"Authorization: Token %s\" %s/v1/images/%s/layer\n",
		token, registry.RegistryHost, args[1])
}

func main() {
	OptParse()
}
