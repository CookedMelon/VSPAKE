package main

import (
	"fmt"
	"sync"
)

type ConnectPool struct {
	mu sync.Mutex
	m  map[string]*connection
}
type connection struct {
	SessionKey []byte
	Akey       []byte
}

func (cp *ConnectPool) Get(socket string) (*connection, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	connBytes, ok := cp.m[socket]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return connBytes, nil
}
func (cp *ConnectPool) Set(socket string, conn *connection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.m[socket] = conn
}
func (cp *ConnectPool) Delete(socket string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	delete(cp.m, socket)
}
func CtearConnectPool() *ConnectPool {
	return &ConnectPool{
		m: make(map[string]*connection),
	}
}
