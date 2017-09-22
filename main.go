package main

import (
	"fmt"
	"os"
)

func init() {
	nodesinit()
}

func help(msg string, exitCode int) {
	fmt.Printf(`%s <command> [args]

This tool was written to serve as a personal Black Desert Database. It is
missing a ton of information, likely has some incorrect information, and
probably is only useful to me.

nodes [worker city]
    Shows information about your node network. You can provide a [worker city]
    to just display what is being produced by workers from that city.

nodes path <node a> [node b]
    Shows the best way to connect <node a> to your network, or to [node b] if
    that is given.

nodes search [costs] <phrase>
    Shows information about the nodes that match the search <phrase> given.
    If you give the "costs" option, the contribution points needed to connect
    each matching nodes to your network will be shown as well.

table search <file> <phrase>
    Shows the lines in the table <file> that match the search <phrase> given.

table search-column <file> <column> <phrase>
    Shows the lines in the table <file> that match the search <phrase> given,
    but only within the <column> given.

csv
    This will translate a CSV file from stdin to a table file to stdout.

If you have a file named "owned" in the current directory, it will be read as
the list of nodes you own, one node per line. If a line ends with
" -- <worker city>" it will mark the node as having a worker assigned to it
from the <worker city>.

Example "owned" file showing a common case where a Velian worker is working on
the Ancient Stone Chamber excavation node:

Velia
Bartali Farm
Toscani Farm
Forest of Seclusion
Ancient Stone Chamber
Ancient Stone Chamber: A -- Velia
`, os.Args[0])
	if msg != "" {
		fmt.Println("")
		fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(exitCode)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		help("", 0)
	}
	switch args[0] {
	case "nodes":
		nodesCommand(args[1:])
	case "table":
		table(args[1:])
	case "csv":
		csvToTable(args[1:])
	default:
		help(fmt.Sprintf("Unknown command %q.", args[0]), 1)
	}
}
