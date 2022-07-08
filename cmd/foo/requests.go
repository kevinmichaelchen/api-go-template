package main

import (
	"context"
	"github.com/kevinmichaelchen/api-go-template/internal/idl/coop/drivers/foo/v1beta1"
	"github.com/spf13/cobra"
	"log"
)

func getFoo(cmd *cobra.Command, args []string) {
	// Create request
	req := &v1beta1.GetFooRequest{
		Id: id,
	}

	// Log request
	s, err := marshalProto(req)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}
	log.Println(s)

	// Execute request
	client := v1beta1.NewFooServiceClient(conn)
	res, err := client.GetFoo(context.Background(), req)
	if err != nil {
		log.Fatalf("gRPC request failed: %v", err)
	}

	// Print response
	s, err = marshalProto(res)
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}
	log.Println(s)
}

func createFoo(cmd *cobra.Command, args []string) {
	// Create request
	req := &v1beta1.CreateFooRequest{
		Name: name,
	}

	// Log request
	s, err := marshalProto(req)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}
	log.Println(s)

	// Execute request
	client := v1beta1.NewFooServiceClient(conn)
	res, err := client.CreateFoo(context.Background(), req)
	if err != nil {
		log.Fatalf("gRPC request failed: %v", err)
	}

	// Print response
	s, err = marshalProto(res)
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}
	log.Println(s)
}
