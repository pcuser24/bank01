package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/user2410/simplebank/api"
	db "github.com/user2410/simplebank/db/sqlc"
	"github.com/user2410/simplebank/storage"
	"github.com/user2410/simplebank/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Establish new connection to database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer func() {
		ch := make(chan bool, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go func() {
			conn.Close()
			ch <- true
		}()
		select {
		case <-ctx.Done():
			log.Println("Database forced to disconnect")
			return
		case <-ch:
			log.Println("Disconnected from database")
		}
	}()

	store := db.NewStore(conn)

	s3Client, err := storage.NewS3Client(
		config.AWSRegion,
		config.AWSAccessKeyID,
		config.AWSSecretAccessKey,
		config.AWSS3Endpoint,
	)
	if err != nil {
		log.Fatal("Error while initializing AWS S3 client", err)
	}
	storage := storage.NewS3Storage(s3Client, config.AWSS3Bucket)

	// Start the server
	server, err := api.NewServer(config, store, storage)
	if err != nil {
		log.Fatal("Failed to create server:", err.Error())
	}
	server.Start(config.ServerAddress)
}
