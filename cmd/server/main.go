package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/aapldev00/synthetic_gen_gaas/internal/generator/providers"
	"github.com/aapldev00/synthetic_gen_gaas/internal/metrics"
	"github.com/aapldev00/synthetic_gen_gaas/internal/server"
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// main serves as the bootstrap entry point for the GaaS Engine.
// It orchestrates the initialization of the network transport layer,
// telemetry services, and gRPC service registration, while managing
// the application lifecycle through OS signal interception.
func main() {
	// TCP listener initialization for the primary gRPC transport layer.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("bootstrap: failed to initialize TCP listener on port 50051: %v", err)
	}

	// gRPC server instantiation and service orchestration.
	// This instance serves as the execution context for the streaming handler.
	grpcServer := grpc.NewServer()

	// High-performance generator service registration.
	// Links the networking layer with internal execution planning and concurrency logic.
	genServer := server.NewServer()
	genproto.RegisterGeneratorServiceServer(grpcServer, genServer)

	// gRPC Reflection enables runtime service discovery.
	// Critical for infrastructure interoperability and diagnostic tooling.
	reflection.Register(grpcServer)

	// Prometheus telemetry server initialization.
	// Runs on a dedicated port to isolate monitoring traffic from the primary data plane.
	go func() {
		log.Printf("Telemetry server operational on :9090")
		if err := metrics.StartMetricsServer(":9090"); err != nil {
			log.Printf("bootstrap: metrics server failure: %v", err)
		}
	}()

	// Asynchronous gRPC service execution.
	// Non-blocking execution allows the main thread to remain responsive
	// to lifecycle management signals.
	go func() {
		log.Printf("GaaS Engine operational on %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("bootstrap: server execution failure: %v", err)
		}
	}()

	// Lifecycle synchronization and graceful decommissioning.
	// Intercepts SIGINT and SIGTERM to ensure active streams are flushed
	// and memory pools are released before process exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Initiating graceful shutdown sequence...")

	// GracefulStop facilitates a controlled shutdown, ensuring resource
	// integrity across all active connections.
	grpcServer.GracefulStop()

	log.Println("GaaS Engine successfully decommissioned.")
}
