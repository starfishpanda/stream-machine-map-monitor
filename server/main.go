package main

import (
	"context"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"os/signal"
	pb "stream-machine-map-monitor/proto"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// MachineManager handles all machines through goroutines
type MachineManager struct {
	pb.UnimplementedMachineMapServer
	machines map[uint32]*Machine // map of machines id to machine pointers
	mu sync.RWMutex // thread-locking map of machines
	nextID uint32
	stopChans map[uint32]chan struct{} // signal goroutines to stop
	updateRate time.Duration // rate at which time elapses for every machine
}

// Creates the MachineManager at startup
func NewMachineManager() *MachineManager {
	return &MachineManager{
		machines: make(map[uint32]*Machine),
		stopChans: make(map[uint32]chan struct{}),
		updateRate: 1000 * time.Millisecond,
	}
}

// Machine represents a robot and its state
type Machine struct {
	ID uint32
	Location *pb.GPS 
	IsPaused bool
	mutex sync.RWMutex
	FuelLevel float32
	brownian *BrownianMotion
}

type BrownianMotion struct{
	stepSizeLatLon float64
	stepSizeAlt float64
	fuelDrainRate float32
}

// createMachine creates a new machine with initial position
func (mm *MachineManager) createMachine() *Machine {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	machine := &Machine{
		ID: mm.nextID,
		Location: &pb.GPS{
			Lat: 47.695185 + (0.0001 * (float64(mm.nextID%5) - 2)), // Start machine around Sammamish Valley, and nudge based on manipulation of ID
			Lon: -122.145161 + (0.0001 * (float64(mm.nextID%5) - 2)),
			Alt: float32(0 + mm.nextID%50), // Different starting altitudes from sea level
		},
		IsPaused: true,
		FuelLevel: 100.0, // Initially 100% FuelLevel
		brownian: &BrownianMotion{
			stepSizeLatLon: 0.0001,
			stepSizeAlt: 1.0,
			fuelDrainRate: 0.1,
		},
	}

	mm.machines[mm.nextID] = machine
	mm.nextID++

	return machine
}

func (mm *MachineManager) machineToProto(machine *Machine) *pb.Machine {
	machine.mutex.RLock()
	defer machine.mutex.RUnlock()

	return &pb.Machine{
		Id: machine.ID,
		Location: machine.Location,
		FuelLevel: machine.FuelLevel,
		IsPaused: machine.IsPaused,
	}
}

func (mm *MachineManager) startMachineMovement(machine *Machine) {
	// Create a new channel for the new machine so goroutine for movement can be stopped
	stopChan := make(chan struct{})
	mm.stopChans[machine.ID] = stopChan

	// goroutine to update machine GPS location
	go func() {
		ticker := time.NewTicker(mm.updateRate)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return

			case <-ticker.C:
				machine.mutex.Lock()
				if !machine.IsPaused && machine.FuelLevel > 0 {
					// Add Brownian motion to machine GPS location
					machine.Location.Lat += (2.0 * (rand.Float64() - 0.5)) * machine.brownian.stepSizeLatLon
					machine.Location.Lon += (2.0 * (rand.Float64() - 0.5)) * machine.brownian.stepSizeLatLon
					machine.Location.Alt += float32((2.0 * (rand.Float64() - 0.5)) * machine.brownian.stepSizeAlt)

					// Fuel drain with movement
					machine.FuelLevel -= machine.brownian.fuelDrainRate
					if machine.FuelLevel < 0 {
						machine.FuelLevel = 0
						machine.IsPaused = true
					}
				}
				machine.mutex.Unlock()

			}
		}
	}()
}

// gRPC method to pause machine
func (mm *MachineManager) Pause(ctx context.Context, req *pb.Machine) (*pb.Machine, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	machine, exists := mm.machines[req.Id]
	if !exists {
		return nil, nil
	}

	machine.mutex.Lock()
	machine.IsPaused = true
	machine.mutex.Unlock()

	return mm.machineToProto(machine), nil
}

// gRPC method to unpause machine
func (mm *MachineManager) UnPause(ctx context.Context, req *pb.Machine) (*pb.Machine, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	machine, exists := mm.machines[req.Id]
	if !exists {
		return nil, nil
	}
	machine.mutex.Lock()
	machine.IsPaused = false
	machine.mutex.Unlock()

	return mm.machineToProto(machine), nil
}

// gRPC method implementation (same as from .proto). Instantiate machine and stream it as protobuf
func (mm *MachineManager) MachineStream(req *pb.MachineStreamRequest, stream pb.MachineMap_MachineStreamServer) error {
	machine := mm.createMachine()

	// Start Brownian Motion of machine
	mm.startMachineMovement(machine)

	// Stream updates until client (WebSocket) disconnects
	for {
		select {
		case <-stream.Context().Done():
			mm.removeMachine(machine.ID)
			return nil
		case <-time.After(mm.updateRate):
			// Send current machine state as protobuf
			if err := stream.Send(mm.machineToProto(machine)); err != nil {
				mm.removeMachine(machine.ID)
				return err
			}
		}
	}
}

func (mm *MachineManager) removeMachine(id uint32) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if stopChan, exists := mm.stopChans[id]; exists {
		close(stopChan)
		delete(mm.stopChans, id)
	}
	delete(mm.machines, id)
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failted to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	machineManager := NewMachineManager()

	pb.RegisterMachineMapServer(grpcServer, machineManager) // Connect the MachineMapServer interface in machineManager to the gRPC server

	// Channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Run gRPC server in a separate goroutine so main can listen for sigterm

	go func() {
		log.Printf("Server starting on port 50051...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-sigChan
	log.Println("Server gracefully shutting down...")

	grpcServer.GracefulStop()
	log.Println("Server stopped.")

}
