package main

import (
	"fmt"
	"net"
	"sync"

	"gopkg.in/natefinch/npipe.v2"
)

var (
	pipeClients []net.Conn
	pipeMutex   sync.Mutex
)

func startPipeServer() {
	ln, err := npipe.Listen(`\\.\pipe\SpaceWorkspace`)
	if err != nil {
		fmt.Println("Failed to start pipe server:", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		pipeMutex.Lock()
		pipeClients = append(pipeClients, conn)
		pipeMutex.Unlock()
	}
}

func broadcastWorkspace(ws int) {
	pipeMutex.Lock()
	defer pipeMutex.Unlock()

	msg := fmt.Appendf(nil, "%d\n", ws)
	var activeClients []net.Conn

	for _, conn := range pipeClients {
		_, err := conn.Write(msg)
		if err == nil {
			activeClients = append(activeClients, conn)
		} else {
			conn.Close()
		}
	}
	pipeClients = activeClients
}
