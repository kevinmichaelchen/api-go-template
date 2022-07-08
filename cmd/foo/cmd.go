package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "foo",
	Short: "Foo is a tool to make gRPC requests",
	Long: `Foo is a tool to make gRPC requests built with
                love by The Drivers Coop.`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var conn *grpc.ClientConn

var id string
var name string

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(getCmd)

	createCmd.PersistentFlags().StringVarP(&name, "name", "", "", "The name of the new Foo")
	getCmd.PersistentFlags().StringVarP(&id, "id", "", "", "The ID of the Foo to get")

	log.Println("Initializing gRPC connection...")
	var err error
	conn, err = grpc.Dial("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial gRPC connection: %v", err)
	}
	log.Println("Initialized gRPC connection.")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Foo",
	Long:  `Create a new Foo`,
	Run:   createFoo,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a new Foo",
	Long:  `Get a new Foo`,
	Run:   getFoo,
}

func marshalProto(m proto.Message) (string, error) {
	b, err := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
