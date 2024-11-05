package logic

import (
	"math/rand"
	"time"
)

func SimulateVehicle(parkingLot *ParkingLot) {
	for {
		time.Sleep(time.Duration(rand.ExpFloat64()) * time.Second) // Poisson arrival
		parkingLot.entranceChan <- struct{}{}                      // Request entrance
		spaceIndex := parkingLot.FindAvailableSpace()

		if spaceIndex != -1 {
			// Simulate parking
			time.Sleep(time.Duration(rand.Intn(3)+3) * time.Second)
			parkingLot.FreeSpace(spaceIndex)
			parkingLot.exitChan <- struct{}{} // Signal exit
		} else {
			<-parkingLot.entranceChan // Exit if no space available
		}
	}
}
