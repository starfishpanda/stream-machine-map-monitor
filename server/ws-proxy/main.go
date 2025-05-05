package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	pb "stream-machine-map-monitor/proto"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// A gRPC client stub that connects to 
type ProxyServer struct {
	// Use client stub for "local" function calls
	grpcClient pb.MachineMapClient
}

func NewProxyServer() (*ProxyServer, error){
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewMachineMapClient(conn)
	return &ProxyServer{grpcClient: client},nil
	
}

// Handle WebSocket machine connections
func (s *ProxyServer) handleMachine(w http.ResponseWriter, r *http.Request){
	// Initialize connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Create gRPC stream
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := s.grpcClient.MachineStream(ctx, &pb.MachineStreamRequest{})
	if err != nil {
		log.Printf("Failed to start machine stream: %v", err)
		return
	}

	// Forward gRPC stream to WebSocket
	go func() {
		for {
			machine, err := stream.Recv()
			if err != nil {
				log.Printf("Stream ended: %v", err)
				cancel()
				return
			}


	// Convert to JSON and send to WebSocket
	machineJSON, err := json.Marshal(machine)
	if err != nil {
		log.Printf("Failed to marshal machine: %v", err)
		continue
	}

	if err := conn.WriteMessage(websocket.TextMessage, machineJSON); err != nil {
		log.Printf("Failed to write message: %v", err)
		cancel()
		return
	}
}
}()


}
func main() {
	proxy, err := NewProxyServer()
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	http.HandleFunc("/machine", proxy.handleMachine)
	
	log.Println("WebSocket proxy server listening on port :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}