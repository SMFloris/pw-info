package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	cli "github.com/jawher/mow.cli"
	"github.com/smfloris/pw-info/pipewire"
)

func main() {
	app := cli.App("pw-info", "Filters pw-cli output in order to get objects based on properties")

	pwCli := exec.Command("pw-cli", "dump", "all")
	out, err := pwCli.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	var nodes []pipewire.PipewireNode
	nodes = pipewire.ParseOutput(string(out))

	app.Command("find-node-property", "gets the value of a specific property from a node", func(cmd *cli.Cmd) {
		var (
			id    = cmd.IntArg("ID", 0, "id of the node")
			name  = cmd.StringArg("NAME", "", "property name to return")
			short = cmd.BoolOpt("s short", false, "do not output key names, only values")
		)

		cmd.Action = func() {
			if *name == "all" {
				properties, err := pipewire.ListProperties(nodes, *id)
				if err != nil {
					log.Fatal(err)
				}
				PrintProperties(properties, *short)
				return
			}
			property, err := pipewire.GetProperty(nodes, *id, *name)
			if err != nil {
				log.Fatal(err)
			}
			PrintProperty(property, *short)
		}
	})

	app.Command("find-node", "Returns a node based on a property value", func(cmd *cli.Cmd) {
		var (
			name  = cmd.StringArg("NAME", "", "name of the property")
			value = cmd.StringArg("VALUE", "", "value of the property")
			short = cmd.BoolOpt("s short", false, "do not output key names, only values")
		)

		cmd.Action = func() {
			node, err := pipewire.FindNode(nodes, *name, *value)
			if err != nil {
				log.Fatal(err)
			}
			PrintNode(node, *short)
		}
	})

	app.Run(os.Args)
}

func PrintNode(node pipewire.PipewireNode, short bool) {
	fmt.Printf("%s%d\n", getShortKey("id: ", short), node.Id)
	fmt.Printf("%s%s\n", getShortKey("type: ", short), node.Type)
}

func PrintProperties(properties []pipewire.PipewireNodeProperty, short bool) {
	for _, property := range properties {
		PrintProperty(property, short)
	}
}

func PrintProperty(property pipewire.PipewireNodeProperty, short bool) {
	fmt.Printf("%s%s\n", getShortKey(property.Key, short), property.Value)
}

func getShortKey(key string, short bool) string {
	if short {
		return ""
	}

	return key + ": "
}
