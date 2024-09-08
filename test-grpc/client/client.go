package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"time"

	"grpcserver/routeguide"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func printFeature(client routeguide.RouteGuideClient, point *routeguide.Point) {
	log.Printf("Getting feature point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("client.GetFeature failed: %v", err)
	}
	log.Println(feature)
}

func printFeatures(client routeguide.RouteGuideClient, rect *routeguide.Rectangle) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("client.ListFeatures failed: %v", err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.ListFeatures failed: %v", err)
		}
		log.Printf("Feature: name: %q, point:(%v, %v)", feature.GetName(),
			feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}
}

func randomPoint(r *rand.Rand) *routeguide.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &routeguide.Point{Latitude: lat, Longitude: long}
}

func runRecordRoute(client routeguide.RouteGuideClient) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	pointCount := int(r.Int31n(100)) + 2
	var points []*routeguide.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("Traversing %d points.", len(points))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("client.RecordRoute failed: %v", err)
	}

	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("client.RecordRoute: stream.Send(%v) failed: %v", point, err)
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("client.RecordRoute failed: %v", err)
	}
	log.Printf("Route summary: %v", reply)
}

func runRouteChat(client routeguide.RouteGuideClient) {
	notes := []*routeguide.RouteNote{
		{Location: &routeguide.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}
			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", note, err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func main() {
	flag.Parse()

	altsTC := alts.NewClientCreds(alts.DefaultClientOptions())

	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(altsTC))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := routeguide.NewRouteGuideClient(conn)

	// valid feature
	printFeature(client, &routeguide.Point{Latitude: 409146138, Longitude: -746188906})

	// Feature missing.
	printFeature(client, &routeguide.Point{Latitude: 0, Longitude: 0})

	// Looking for features between 40, -75 and 42, -73.
	printFeatures(client, &routeguide.Rectangle{
		Lo: &routeguide.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &routeguide.Point{Latitude: 420000000, Longitude: -730000000},
	})

	// RecordRoute
	runRecordRoute(client)

	// RouteChat
	runRouteChat(client)
}
