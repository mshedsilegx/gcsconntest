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
var version string

// config struct to hold command-line flags
type config struct {
	credFile     string
	bucketName   string
	bucketPrefix string
	projectID    string
	maxObjects   int
	showVersion  bool
}

func main() {
	cfg := parseFlags()
	if err := run(cfg); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// run executes the main application logic.
func run(cfg config) error {
	if cfg.showVersion {
		fmt.Println("GCS Connection Tester - Version:", version)
		return nil
	}

	if cfg.credFile == "" || cfg.bucketName == "" || cfg.projectID == "" {
		return fmt.Errorf("credentials file, bucket name, and project ID are required")
	}

	data, err := os.ReadFile(cfg.credFile)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %v", err)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx,
		option.WithAuthCredentialsJSON(option.ServiceAccount, data),
	)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Warning: failed to close GCS client: %v", err)
		}
	}()

	if err := listGCSObjects(ctx, client, cfg); err != nil {
		return fmt.Errorf("failed to list GCS objects: %w", err)
	}

	fmt.Println("GCS connectivity test successful.")
	return nil
}

// listGCSObjects lists objects in a GCS bucket.
func listGCSObjects(ctx context.Context, client *storage.Client, cfg config) error {
	bucket := client.Bucket(cfg.bucketName)

	if _, err := bucket.Attrs(ctx); err != nil {
		return fmt.Errorf("failed to get bucket attributes: %w", err)
	}

	if cfg.bucketPrefix == "" {
		fmt.Printf("Listing up to %d objects in bucket %s\n", cfg.maxObjects, cfg.bucketName)
	} else {
		fmt.Printf("Listing up to %d objects in bucket %s with prefix %s:\n", cfg.maxObjects, cfg.bucketName, cfg.bucketPrefix)
	}

	it := bucket.Objects(ctx, &storage.Query{Prefix: cfg.bucketPrefix})
	count := 0
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list objects in bucket: %w", err)
		}
		fmt.Println(attrs.Name)
		count++
		if count >= cfg.maxObjects {
			break
		}
	}

	return nil
}

// parseFlags parses the command-line flags and returns them in a config struct.
func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.credFile, "credentials", "", "Path to JSON credential file")
	flag.StringVar(&cfg.bucketName, "bucket", "", "GCS bucket name")
	flag.StringVar(&cfg.bucketPrefix, "prefix", "", "GCS bucket prefix")
	flag.StringVar(&cfg.projectID, "project", "", "GCP project ID")
	flag.IntVar(&cfg.maxObjects, "max", 10, "Maximum number of objects to list")
	flag.BoolVar(&cfg.showVersion, "version", false, "Display application version")
	flag.Parse()
	return cfg
}
