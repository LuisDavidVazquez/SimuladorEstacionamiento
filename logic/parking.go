package logic

import (
	"sync"
)

const TotalSpaces = 20

type ParkingLot struct {
	mu           sync.Mutex
	spaces       []bool
	entranceChan chan struct{}
	exitChan     chan struct{}
}

func NewParkingLot() *ParkingLot {
	return &ParkingLot{
		spaces:       make([]bool, TotalSpaces),
		entranceChan: make(chan struct{}, 1), // Semaphore for entrance
		exitChan:     make(chan struct{}, 1), // Semaphore for exit
	}
}

func (p *ParkingLot) FindAvailableSpace() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < TotalSpaces; i++ {
		if !p.spaces[i] {
			p.spaces[i] = true
			return i
		}
	}
	return -1 // No space available
}

func (p *ParkingLot) FreeSpace(index int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.spaces[index] = false
}
