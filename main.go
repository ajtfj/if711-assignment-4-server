package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/ajtfj/graph"
)

const (
	GRAPH_FILE = "graph.txt"
)

type GraphRPC struct {
	Graph *graph.Graph
}

func NewGraphRPC() *GraphRPC {
	return &GraphRPC{
		Graph: graph.NewGraph(),
	}
}

func (r *GraphRPC) ShortestPath(args *ShortestPathArgs, reply *ShortestPathReply) error {
	path, err := r.Graph.ShortestPath(args.Ori, args.Dest)
	if err != nil {
		return err
	}
	reply.Path = path
	return nil
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("undefined PORT")
	}

	graphRPC := NewGraphRPC()

	if err := setupGraph(graphRPC.Graph); err != nil {
		log.Fatal(err)
	}

	server := rpc.NewServer()
	if err := server.RegisterName("Graph", graphRPC); err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("localhost:%s", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	server.Accept(ln)
}

func setupGraph(graph *graph.Graph) error {
	file, err := os.Open(GRAPH_FILE)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		inputLine := scanner.Text()
		u, v, weight, err := parceGraphInputLine(inputLine)
		if err != nil {
			return err
		}
		graph.AddEdge(u, v, weight)
	}

	return nil
}

func parceGraphInputLine(inputLine string) (graph.Node, graph.Node, int, error) {
	matches := strings.Split(inputLine, " ")
	if len(matches) < 3 {
		return graph.Node(""), graph.Node(""), 0, fmt.Errorf("invalid input")
	}

	weight, err := strconv.ParseInt(matches[2], 10, 0)
	if err != nil {
		return graph.Node(""), graph.Node(""), 0, err
	}

	return graph.Node(matches[0]), graph.Node(matches[1]), int(weight), nil
}

type ShortestPathArgs struct {
	Ori  graph.Node
	Dest graph.Node
}

type ShortestPathReply struct {
	Path []graph.Node
}
