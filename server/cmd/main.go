package main

import (
	"context"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/tonymontanapaffpaff/golang-training-university-grpc/proto/go_proto"
	"github.com/tonymontanapaffpaff/golang-training-university-grpc/server/pkg/api"
	"github.com/tonymontanapaffpaff/golang-training-university-grpc/server/pkg/data"
	"github.com/tonymontanapaffpaff/golang-training-university-grpc/server/pkg/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	serverEndpoint = os.Getenv("SERVER_ENDPOINT")

	host     = os.Getenv("DB_USERS_HOST")
	port     = os.Getenv("DB_USERS_PORT")
	user     = os.Getenv("DB_USERS_USER")
	dbname   = os.Getenv("DB_USERS_DBNAME")
	password = os.Getenv("DB_USERS_PASSWORD")
	sslmode  = os.Getenv("DB_USERS_SSL")
)

func init() {
	if serverEndpoint == "" {
		serverEndpoint = "localhost:8080"
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if dbname == "" {
		dbname = "university"
	}
	if password == "" {
		password = "postgres"
	}
	if sslmode == "" {
		sslmode = "disable"
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx, cancel = context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	conn, err := connectWithTimeout(ctx)
	if err != nil {
		log.Fatalf("can't connect to database, error: %v", err)
	}

	listener, err := net.Listen("tcp", serverEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	reflection.Register(server)

	pb.RegisterCourseServiceServer(server, api.NewCourseServer(data.CourseData{Db: conn}))

	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func connectWithTimeout(ctx context.Context) (*gorm.DB, error) {
	for {
		time.Sleep(2 * time.Second)
		conn, err := db.GetConnection(host, port, user, dbname, password, sslmode)
		if err == nil {
			return conn, nil
		}
		select {
		case <-ctx.Done():
			return nil, err
		default:
			continue
		}
	}
}
