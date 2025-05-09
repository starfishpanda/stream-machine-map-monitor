package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	pb "stream-machine-map-monitor/proto"
	"syscall"
	"time"

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
	conn *grpc.ClientConn
}

func NewProxyServer() (*ProxyServer, error){
	// Get gRPC server address from environment variable or use default
	grpcServerAddr := os.Getenv("GRPC_SERVER")
	if grpcServerAddr == "" {
			grpcServerAddr = "localhost:50051" // Default if not specified
	}
	
	log.Printf("Connecting to gRPC server at %s", grpcServerAddr)
	conn, err := grpc.NewClient(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewMachineMapClient(conn)
	return &ProxyServer{grpcClient: client, conn: conn},nil
	
}

// Method for proxy server to close gRPC Client connection 
func (s *ProxyServer) Close() {
	if s.conn != nil {
		s.conn.Close()
		log.Println("gRPC connection closed")
	}
}

// Handle incoming Websocket connection requests from browser client
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

	// Goroutine to forward gRPC stream to WebSocket
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
		// log.Printf("WS Proxy: sending response in handleMachine goroutine before websocket.TextMessage: \n%s",machineJSON)

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

// Handle incoming WebSocket messages (disconnects, or pause and unpause requests)
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
	
	// log.Printf("WS Proxy: sending response to incoming command in handleMachine before websocket.TextMessage: \n%s",responseJSON)
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

	// CORS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Write([]byte("WebSocket proxy server running"))
	})

	// Create server instance
	server := &http.Server{
		Addr: ":3001",
		Handler: nil,
	}

	// Create channel to listen for SIGTERM
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	
	// Run server in goroutine so main thread can listen for SIGINT or SIGTERM
	go func() {
		log.Println("WebSocket proxy server listening on port :3001")
		if err:= server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}

	}()
	// Listening on channel for SIGINT or SIGTERM
	<-stopChan
	log.Println("Shutting down server...")

	// Create timeout for shutdown so connections can gracefully close
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}