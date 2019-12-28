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
	nodeCount        int
	natBuffer        packet // stores a packet that the NAT may route as needed in order to manage traffic
	natBufferHistory int    // remember the last y value sent by the NAT
	nodes            []intcode.Computer
}

func makeNetwork(nodeCount int, NIC string) network {

	net := network{
		nodes:            make([]intcode.Computer, nodeCount),
		nodeCount:        nodeCount,
		natBufferHistory: -9999,
	}

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

func (net *network) routePackets() bool {
	// an easy way of letting the caller know whether the
	// puzzle condition has been met.
	exitConditionFlag := false

	// we need to keep track of what nodes received messages
	// and modify some of them in-flight depending on the
	// network circumstances,
	// so we accumulate all messages before distributing them
	msgQueues := make(map[int][]int)
	for i := 0; i < net.nodeCount; i++ {
		msgQueues[i] = make([]int, 0)
	}

	for i := 0; i < net.nodeCount; i++ {
		buffer := net.nodes[i].FlushOutput()
		packets := unpackPacketBuffer(buffer)

		for _, p := range packets {
			if (p.destination < 0 || p.destination >= net.nodeCount) && p.destination != 255 {
				// this should not happen
				fmt.Println("Invalid destination address.", p)
				continue
			} else if p.destination == 255 {
				// Packages routed to the NAT.
				// Only remember the last one received.
				net.natBuffer.destination = p.destination
				net.natBuffer.x = p.x
				net.natBuffer.y = p.y
				continue
			}
			msgQueues[p.destination] = append(msgQueues[p.destination], p.x)
			msgQueues[p.destination] = append(msgQueues[p.destination], p.y)
		}

	}

	emptyCount := 0
	for i := net.nodeCount - 1; i >= 0; i-- {

		// any nodes that did not receive input during this cycle get a -1
		// EXCEPT for the case where ALL nodes are idle.
		// If all nodes are idle (i.e., if all nodes will receive a -1),
		// then node 0 receives the package that is currently stored in the NAT
		if emptyCount < net.nodeCount-1 && len(msgQueues[i]) == 0 {
			msgQueues[i] = append(msgQueues[i], -1)
			emptyCount++

		} else if emptyCount == net.nodeCount-1 && len(msgQueues[i]) == 0 {

			msgQueues[i] = append(msgQueues[i], net.natBuffer.x)
			msgQueues[i] = append(msgQueues[i], net.natBuffer.y)
			emptyCount++

			// fmt.Println(net.natBuffer)
			exitConditionFlag = net.natBufferHistory == net.natBuffer.y
			net.natBufferHistory = net.natBuffer.y
		}

		net.nodes[i].QueueInput(msgQueues[i]...)
	}

	return exitConditionFlag
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

	for true {
		net.tick()
		net.routePackets()
		if net.natBuffer.y != 0 {
			break
		}
	}

	fmt.Println("[Part 1] The y value of the first packet sent to the NAT is", net.natBuffer.y)
}

func part2() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Can't read input file.")
	}

	net := makeNetwork(50, string(data))

	for true {
		net.tick()
		if net.routePackets() {
			break
		}
	}

	fmt.Println("[Part 2] The first repeated y value sent by the NAT is", net.natBuffer.y)
}

func main() {
	part1()
	part2()
}
