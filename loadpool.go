package main

import "sync"

type ConnectPool struct {
	mu sync.Mutex
	m  map[string]*connection
}
type connection struct {
	SessionKey []byte
	Akey       []byte
}

func (cp *ConnectPool) Get(key []byte) (*connection, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	return cp.m[string(key)], nil
}
