package pkg

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DialCollector creates a gRPC connection to the collector.
func DialCollector(host, port string) (*grpc.ClientConn, error) {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	defer cancel()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port),
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	log.Printf("Status: %v", conn.GetState())

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}
