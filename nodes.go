package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var nodes = map[string]*node{}

type node struct {
	name               string
	contributionPoints int
	owned              bool
	closestWorker      string
	assignedWorker     string
	produces           []string
}

func addNode(name string, cp int) {
	nodes[name] = &node{name: name, contributionPoints: cp}
	if cp == 0 {
		nodes[name].owned = true
	}
}

func addProductionNode(parent string, name string, cp int, closestWorker string, produces ...string) {
	name = parent + ": " + name
	addNode(name, cp)
	addConnection(parent, name)
	nodes[name].closestWorker = closestWorker
	for _, p := range produces {
		nodes[name].produces = append(nodes[name].produces, p)
	}
	addConnection(name, parent)
}

func (n *node) String() string {
	var s string
	if n.owned {
		s = fmt.Sprintf("%s (%d) owned", n.name, n.contributionPoints)
	} else {
		s = fmt.Sprintf("%s [%d]", n.name, n.contributionPoints)
	}
	if n.closestWorker != "" {
		s += ", closest worker from " + n.closestWorker
	}
	if len(n.produces) > 0 {
		s += ", produces:"
		for i, p := range n.produces {
			if i != 0 {
				s += ","
			}
			s += " " + p
		}
	}
	return s
}

var connections = map[string]map[string]struct{}{}

func addConnection(a string, b string) {
	if _, ok := connections[a]; !ok {
		connections[a] = map[string]struct{}{b: struct{}{}}
	} else {
		connections[a][b] = struct{}{}
	}
}

type costNode struct {
	cost int
	node *node
}

type costNodes []*costNode

func (cns costNodes) Len() int {
	return len(cns)
}

func (cns costNodes) Swap(x, y int) {
	cns[x], cns[y] = cns[y], cns[x]
}

func (cns costNodes) Less(x, y int) bool {
	v := cns[x].cost - cns[y].cost
	if v < 0 {
		return true
	} else if v > 0 {
		return false
	}
	return cns[x].node.name < cns[y].node.name
}

func nodesCommand(args []string) {
	f, err := os.Open("owned")
	if err == nil {
		scanner := bufio.NewScanner(f)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			name := line
			var worker string
			t := strings.SplitN(line, " -- ", 2)
			if len(t) == 2 {
				name = t[0]
				workerL := strings.ToLower(t[1])
				for n := range nodes {
					if strings.ToLower(n) == workerL {
						worker = n
					}
				}
				if worker == "" {
					help(fmt.Sprintf("Could not find node %q referenced on line %d: %q\n", t[1], lineNumber, line), 1)
				}
			}
			nameL := strings.ToLower(name)
			for n := range nodes {
				if strings.ToLower(n) == nameL {
					nodes[n].owned = true
					if worker != "" {
						nodes[n].assignedWorker = worker
					}
					nameL = ""
					break
				}
			}
			if nameL != "" {
				help(fmt.Sprintf("Could not find node %q referenced on line %d: %q\n", name, lineNumber, line), 1)
			}
		}
		f.Close()
	}
	var cmd string
	if len(args) > 0 {
		cmd = args[0]
		args = args[1:]
	}
	switch cmd {
	case "path":
		if len(args) < 1 {
			help(fmt.Sprintf("Path request had no parameters"), 1)
		}
		if len(args) > 2 {
			help(fmt.Sprintf("Path request had too many parameters"), 1)
		}
		var nodeA string
		var nodeB string
		nodeAL := strings.ToLower(args[0])
		var nodeBL string
		if len(args) == 2 {
			nodeBL = strings.ToLower(args[1])
		}
		for n := range nodes {
			nL := strings.ToLower(n)
			if nL == nodeAL {
				nodeA = n
			}
			if nL == nodeBL {
				nodeB = n
			}
			if nodeA != "" && (nodeB != "" || nodeBL == "") {
				break
			}
		}
		if nodeA == "" {
			help(fmt.Sprintf("Could not find node %q.", args[0]), 1)
		}
		if nodeB == "" && nodeBL != "" {
			help(fmt.Sprintf("Could not find node %q.", args[1]), 1)
		}
		if nodeA == nodeB {
			help(fmt.Sprintf("Both nodes seem to be the same node: %q %q.", args[0], args[1]), 1)
		}
		bestCost, bestPaths := bestPaths(nodeA, nodeB)
		if nodeB == "" {
			fmt.Printf("%d contribution points are needed to connect to %s.\n", bestCost, nodeA)
		} else {
			fmt.Printf("%d contribution points are needed to connect %s to %s.\n", bestCost, nodeA, nodeB)
		}
		for i, pth := range bestPaths {
			if len(bestPaths) > 1 {
				fmt.Printf("Option %d:\n", i+1)
			}
			for j := len(pth) - 1; j >= 0; j-- {
				node := nodes[pth[j]]
				if node.owned {
					if node.contributionPoints == 0 {
						fmt.Printf("          %s (always owned)\n", node.name)
					} else {
						fmt.Printf("          %s (already owned for %d)\n", node.name, node.contributionPoints)
					}
				} else {
					fmt.Printf("   %2d for %s\n", node.contributionPoints, node.name)
				}
			}
		}
	case "search":
		if len(args) < 1 {
			help("No search phrase given.", 1)
		}
		var costs bool
		var search string
		if args[0] == "costs" {
			costs = true
			if len(args) < 2 {
				help("No search phrase given.", 1)
			}
			search = strings.ToLower(strings.Join(args[1:], " "))
		} else {
			search = strings.ToLower(strings.Join(args, " "))
		}
		var matches []string
		for _, n := range nodes {
			if strings.Contains(strings.ToLower(n.name), search) {
				matches = append(matches, n.name)
				continue
			}
			for _, p := range n.produces {
				if strings.Contains(strings.ToLower(p), search) {
					matches = append(matches, n.name)
					break
				}
			}
		}
		if costs {
			var cns costNodes
			for _, n := range matches {
				bestCost, _ := bestPaths(n, "")
				cns = append(cns, &costNode{cost: bestCost, node: nodes[n]})
			}
			sort.Sort(cns)
			for _, cn := range cns {
				fmt.Printf("[%d] %s\n", cn.cost, cn.node)
			}
		} else {
			sort.Strings(matches)
			for _, n := range matches {
				fmt.Println(nodes[n])
			}
		}
	default:
		var filter string
		if cmd != "" {
			filterL := strings.ToLower(cmd)
			for n := range nodes {
				if filterL == strings.ToLower(n) {
					filter = n
					break
				}
			}
			if filter == "" {
				help(fmt.Sprintf("Could not find node %q.", cmd), 1)
			}
		}
		count := 0
		cp := 0
		production := 0
		produces := map[string]int{}
		var notProducing []*node
		workers := 0
		for _, node := range nodes {
			if node.owned {
				count++
				cp += node.contributionPoints
				if len(node.produces) > 0 {
					production++
					if filter == "" {
						if node.assignedWorker == "" {
							notProducing = append(notProducing, node)
						} else {
							workers++
							for _, p := range node.produces {
								produces[p]++
							}
						}
					} else if node.assignedWorker == filter {
						workers++
						for _, p := range node.produces {
							produces[p]++
						}
					}
				}
			}
		}
		if filter == "" {
			fmt.Printf("You own %d nodes for %d contribution points.\n", count, cp)
			if production > 0 {
				fmt.Printf("\n%d are production nodes, of which %d are assigned workers producing the following items:\n", production, workers)
				ps := make([]string, 0, len(produces))
				for p := range produces {
					ps = append(ps, p)
				}
				sort.Strings(ps)
				for _, p := range ps {
					if produces[p] > 1 {
						fmt.Printf("    %s x%d\n", p, produces[p])
					} else {
						fmt.Printf("    %s\n", p)
					}
				}
			}
			if len(notProducing) > 0 {
				fmt.Printf("\nYou have %d production nodes without assigned workers:\n", len(notProducing))
				for _, n := range notProducing {
					fmt.Printf("    %s could produce: %s\n", n.name, strings.Join(n.produces, ", "))
				}
			}
		} else {
			fmt.Printf("%d workers from %s are producing the following items:\n", workers, filter)
			ps := make([]string, 0, len(produces))
			for p := range produces {
				ps = append(ps, p)
			}
			sort.Strings(ps)
			for _, p := range ps {
				if produces[p] > 1 {
					fmt.Printf("    %s x%d\n", p, produces[p])
				} else {
					fmt.Printf("    %s\n", p)
				}
			}
		}
	}
}

func bestPaths(nodeA string, nodeB string) (int, [][]string) {
	deadends := map[string]struct{}{}
	visited := map[string]struct{}{nodeA: struct{}{}}
	cost := 0
	bestCost := int(^uint(0) >> 1)
	var bestPaths [][]string
	var branch func(pth []string, n string) bool
	branch = func(pth []string, n string) bool {
		if !nodes[n].owned {
			cost += nodes[n].contributionPoints
		}
		visited[n] = struct{}{}
		defer func() {
			if !nodes[n].owned {
				cost -= nodes[n].contributionPoints
			}
			delete(visited, n)
		}()
		deadend := true
		for n2 := range connections[n] {
			if _, ok := deadends[n2]; ok {
				continue
			}
			if _, ok := visited[n2]; ok {
				continue
			}
			if n2 == nodeB || (nodeB == "" && nodes[n2].owned) {
				deadend = false
				npth := make([]string, len(pth)+2)
				copy(npth, pth)
				npth[len(pth)] = n
				npth[len(pth)+1] = n2
				ncost := cost
				if !nodes[n2].owned {
					ncost += nodes[n2].contributionPoints
				}
				if bestCost > ncost {
					bestCost = ncost
					bestPaths = bestPaths[:0]
					bestPaths = append(bestPaths, npth)
				} else if bestCost == ncost {
					bestPaths = append(bestPaths, npth)
				}
				continue
			}
			npth := make([]string, len(pth)+1)
			copy(npth, pth)
			npth[len(pth)] = n
			if branch(npth, n2) {
				deadend = false
			}
		}
		if deadend {
			deadends[n] = struct{}{}
			return false
		}
		return true
	}
	branch(nil, nodeA)
	return bestCost, bestPaths
}
