package main

import (
	"log"
	"net"
	pb "stream-machine-map-monitor/proto"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// MachineManager handles all machines through goroutines
type MachineManager struct {
	pb.UnimplementedMachineMapServer
	machines map[uint32]*Machine // map of machines id to pointers to individual machines
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
		updateRate: 100 * time.Millisecond,
	}
}

// Machine represents a robot and its state
type Machine struct {
	ID uint32
	Location *pb.GPS
	IsPaused bool
	mutex sync.RWMutex
	FuelLevel float32
}

// createMachine creates a new machine with initial position
func (mm *MachineManager) createMachine() *Machine {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	machine := &Machine{
		ID: mm.nextID,
		Location: &pb.GPS{
			Lat: 47.695185 + (0.1 * (float64(mm.nextID%5) - 2)), // Start machine around Sammamish Valley
			Lon: -122.145161 + (0.1 * (float64(mm.nextID%5) - 2)),
			Alt: float32(100 + mm.nextID%50), // Different starting altitudes
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


// gRPC method implementation (same as from .proto). Instantiate machine and stream it as protobuf
func (mm *MachineManager) MachineStream(req *pb.MachineStreamRequest, stream pb.MachineMap_MachineStreamServer) error {
	machine := mm.createMachine()

	// Stream updates until client (WebSocket) disconnects
	for {
		select {
		case <-time.After(mm.updateRate):
			// Send current machine state as protobuf
			if err := stream.Send(mm.machineToProto(machine)); err != nil {
				return err
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failted to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	machineManager := NewMachineManager()

	pb.RegisterMachineMapServer(grpcServer, machineManager) // Connect the MachineMapServer interface in machineManager to the gRPC server

	log.Printf("Server starting on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
