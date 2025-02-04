package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Compile-time version variable
var Version string

func main() {
	// Define command-line flags
	credFile := flag.String("credentials", "", "Path to JSON credential file")
	bucketName := flag.String("bucket", "", "GCS bucket name")
	bucketPrefix := flag.String("prefix", "", "GCS bucket prefix")
	projectID := flag.String("project", "", "GCP project ID")
	maxObjects := flag.Int("max", 10, "Maximum number of objects to list")
	showVersion := flag.Bool("version", false, "Display application version")

	flag.Parse()

	if *showVersion {
		fmt.Println("Application Version:", Version)
		return
	}

	// Validate required flags
	if *credFile == "" || *bucketName == "" || *projectID == "" {
		fmt.Println("Credentials file, bucket name, project ID and max objects are required.")
		flag.Usage()
		os.Exit(1)
	}

	// Create a new context
	ctx := context.Background()

	// Initialize the Google Cloud Storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(*credFile))
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	defer client.Close()

	// Get a handle for the bucket
	bucket := client.Bucket(*bucketName)

	// Check if the bucket exists
	_, err = bucket.Attrs(ctx)
	if err != nil {
		log.Fatalf("Failed to get bucket attributes: %v", err)
	}

	// List objects with the given prefix, up to the specified maximum number
	it := bucket.Objects(ctx, &storage.Query{Prefix: *bucketPrefix})
	if *bucketPrefix == "" {
		fmt.Printf("Listing up to %d objects in bucket %s\n", *maxObjects, *bucketName)
	} else {
		fmt.Printf("Listing up to %d objects in bucket %s with prefix %s:\n", *maxObjects, *bucketName, *bucketPrefix)
	}
	count := 0
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list objects in bucket: %v", err)
		}
		fmt.Println(attrs.Name)
		count++
		if count >= *maxObjects {
			break
		}
	}

	fmt.Println("GCS connectivity test successful.")
}
