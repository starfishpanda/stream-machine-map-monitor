package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	pb "stream-machine-map-monitor/proto"
	"sync"

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
	var closeOnce sync.Once
	closeConn := func() {
		closeOnce.Do(func(){
			conn.Close()
		})
	}
	defer closeConn()

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-sigCtx.Done()
		closeConn() // This will trigger the read error and exit your loop
	}()

	// Create gRPC stream
	ctx, cancel := context.WithCancel(sigCtx)
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

// Handle incoming WebSocket messages (disconnects or pause/unpause requests)
for {
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("WebSocket read error: %v", err)
		break
	}

	var request struct {
		Type string `json:"type"`
		ID   uint32 `json:"id"`
	}

	if err := json.Unmarshal(message, &request); err != nil {
		log.Printf("Failed to unmarshal request: %v", err)
		continue
	}

	var response *pb.Machine
	switch request.Type {
	case "pause":
		response, err = s.grpcClient.Pause(ctx, &pb.Machine{Id: request.ID})
	case "unpause":
		response, err = s.grpcClient.UnPause(ctx, &pb.Machine{Id: request.ID})
	default:
		log.Printf("Unknown request type: %s", request.Type)
		continue
	}

	if err != nil {
		log.Printf("Failed to pause/unpause machine: %v", err)
		continue
	}

	// Send confirmation back to client
	responseJSON, _ := json.Marshal(response)
	if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		log.Printf("failed to write response: %v", err)
		cancel()
		break
	}
}
log.Println("Client disconnected, cleaning up")
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