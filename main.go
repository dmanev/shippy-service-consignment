// shippy-service-consignment/main.go
package main

import (

	"context"
	"fmt"
    "log"
    "os"

	// Import the generated protobuf code
	pb "github.com/dmanev/shippy-service-consignment/proto/consignment"
    vesselProto "github.com/dmanev/shippy-service-vessel/proto/vessel"
    "github.com/micro/go-micro"
)

const (
    port = ":50051"
    defaultHost = "datastore:27017"
)

func main() {

    // Set-up micro instance
    srv := micro.NewService(

        // This name must match the package name given in your protobug definition
        micro.Name("shippy.service.consignment"),
    )

    // Init will parse the command line flags.
    srv.Init()

    uri := os.Getenv("DB_HOST")
    if uri == "" {
        uri = defaultHost
    }

    client, err := CreateClient(uri)
    if err != nil {
        log.Panic(err)
    }

    defer client.Disconnect(context.TODO())

    consignmentCollection := client.Database("shippy").Collection("consignments")

    repository := &MongoRepository{consignmentCollection}

    vesselClient := vesselProto.NewVesselServiceClient("shippy.service.vessel", srv.Client())

    h := &handler{repository, vesselClient}

    // Register handlers
    pb.RegisterShippingServiceHandler(srv.Server(), h)

    // Run the server
    if err := srv.Run(); err != nil {
        fmt.Println(err)
    }
}
