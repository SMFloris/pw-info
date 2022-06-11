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
			id   = cmd.IntArg("ID", 0, "id of the node")
			name = cmd.StringArg("NAME", "", "property name to return")
		)

		cmd.Action = func() {
			if *name == "all" {
				properties, err := pipewire.ListProperties(nodes, *id)
				if err != nil {
					log.Fatal(err)
				}
				PrintProperties(properties)
				return
			}
			property, err := pipewire.GetProperty(nodes, *id, *name)
			if err != nil {
				log.Fatal(err)
			}
			PrintProperty(property)
		}
	})

	app.Command("find-node", "Returns a node id based on a property value", func(cmd *cli.Cmd) {
		var (
			name  = cmd.StringArg("NAME", "", "name of the property")
			value = cmd.StringArg("VALUE", "", "value of the property")
		)

		cmd.Action = func() {
			fmt.Println(*name)
			node, err := pipewire.FindNode(nodes, *name, *value)
			if err != nil {
				log.Fatal(err)
			}
			PrintNode(node)
		}
	})

	app.Run(os.Args)
}

func PrintNode(node pipewire.PipewireNode) {
	fmt.Printf("%s: %d\n", "id", node.Id)
	fmt.Printf("%s: %s\n", "type", node.Type)
}

func PrintProperties(properties []pipewire.PipewireNodeProperty) {
	for _, property := range properties {
		PrintProperty(property)
	}
}

func PrintProperty(property pipewire.PipewireNodeProperty) {
	fmt.Printf("%s: %s\n", property.Key, property.Value)
}
