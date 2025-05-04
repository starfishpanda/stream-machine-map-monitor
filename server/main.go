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
	machines map[uint32]*Machine // map of machines
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
	isPaused bool
	mutex sync.RWMutex
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
