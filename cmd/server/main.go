package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aapldev00/synthetic_gen_gaas/internal/server"
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// main serves as the bootstrap entry point for the GaaS Engine.
// It orchestrates the network listener setup, gRPC service registration,
// and implements a graceful shutdown mechanism to ensure resource integrity.
func main() {
	// TCP listener initialization on the standard gRPC port.
	// This provides the underlying transport layer for the HTTP/2 stream.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("bootstrap: failed to listen on port 50051: %v", err)
	}

	// gRPC server instantiation.
	// The configuration allows for future injection of interceptors
	// for telemetry, authentication, and rate limiting.
	grpcServer := grpc.NewServer()

	// High-performance generator service registration.
	// Connects the gRPC networking layer with the internal planning and execution logic.
	genServer := server.NewServer()
	genproto.RegisterGeneratorServiceServer(grpcServer, genServer)

	// gRPC Reflection enables schema discovery for debugging and CLI tools.
	// This is critical for interoperability testing without pre-compiled clients.
	reflection.Register(grpcServer)

	// Asynchronous service execution.
	// Running the server in a separate goroutine allows the main thread
	// to monitor system signals for lifecycle management.
	go func() {
		log.Printf("GaaS Engine operational on %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("bootstrap: server execution failure: %v", err)
		}
	}()

	// Signal interception for graceful termination.
	// Listens for SIGINT and SIGTERM to ensure the engine flushes active
	// streams and releases memory pools before exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Initiating graceful shutdown sequence...")

	// GracefulStop ensures all active gRPC streams are completed
	// and underlying connections are closed properly.
	grpcServer.GracefulStop()

	log.Println("GaaS Engine successfully decommissioned.")
}
