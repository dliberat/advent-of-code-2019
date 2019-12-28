package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

type packet struct {
	destination int
	x           int
	y           int
}

type network struct {
	nodes     []intcode.Computer
	nodeCount int
}

func makeNetwork(nodeCount int, NIC string) network {

	net := network{nodes: make([]intcode.Computer, nodeCount), nodeCount: nodeCount}

	for i := 0; i < nodeCount; i++ {
		net.nodes[i] = intcode.MakeComputer(NIC, nil, nil)
		net.nodes[i].QueueInput(i) // set the network address
	}

	return net
}

func (net *network) tick() {
	for i := 0; i < net.nodeCount; i++ {
		net.nodes[i].Run()
	}
}

func (net *network) routePackets() {
	msgQueues := make(map[int][]int)

	for i := 0; i < net.nodeCount; i++ {
		msgQueues[i] = make([]int, 0)
	}

	for i := 0; i < net.nodeCount; i++ {
		buffer := net.nodes[i].FlushOutput()
		packets := unpackPacketBuffer(buffer)

		for _, p := range packets {
			if p.destination < 0 || p.destination >= net.nodeCount {
				fmt.Println("Invalid destination address.", p)
				continue
			}
			msgQueues[p.destination] = append(msgQueues[p.destination], p.x)
			msgQueues[p.destination] = append(msgQueues[p.destination], p.y)
		}

	}

	for i := 0; i < net.nodeCount; i++ {
		// any nodes that did not receive input during this cycle get a -1
		if len(msgQueues[i]) == 0 {
			msgQueues[i] = append(msgQueues[i], -1)
		}

		net.nodes[i].QueueInput(msgQueues[i]...)
	}
}

func unpackPacketBuffer(buffer []int) []packet {
	packets := make([]packet, 0)

	i := 0
	for i < len(buffer) {
		p := packet{}
		p.destination = buffer[i]
		p.x = buffer[i+1]
		p.y = buffer[i+2]
		i += 3

		packets = append(packets, p)
	}
	return packets
}

func part1() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Can't read input file.")
	}

	net := makeNetwork(50, string(data))

	for i := 0; i < 10000; i++ {
		net.tick()
		net.routePackets()
	}
}

func main() {
	part1()
}
