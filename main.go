package main

// GCS Connection Tester
// This application is designed to verify connectivity and authentication
// to Google Cloud Storage (GCS) by listing objects in a specified bucket.
// It uses a service account JSON key for authentication.

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

// config holds the application configuration derived from command-line flags.
type config struct {
	credFile     string // Path to the Google Cloud service account JSON key file
	bucketName   string // Name of the GCS bucket to test
	bucketPrefix string // Optional prefix to filter objects in the bucket
	projectID    string // Google Cloud Project ID
	maxObjects   int    // Maximum number of objects to list during the test
	showVersion  bool   // If true, print the version and exit
}

func main() {
	// Parse command-line arguments into the config structure.
	cfg := parseFlags()

	// Execute the core application logic.
	if err := run(cfg); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// run executes the main application logic:
// 1. Handles version display request.
// 2. Validates required configuration parameters.
// 3. Initializes the GCS client with provided credentials.
// 4. Performs the listing operation to verify connectivity.
func run(cfg config) error {
	// Section: Version Information
	if cfg.showVersion {
		fmt.Println("GCS Connection Tester - Version:", version)
		return nil
	}

	// Section: Configuration Validation
	if cfg.credFile == "" || cfg.bucketName == "" || cfg.projectID == "" {
		return fmt.Errorf("credentials file, bucket name, and project ID are required")
	}

	// Section: Authentication & Client Initialization
	// Read the service account JSON key file.
	data, err := os.ReadFile(cfg.credFile)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %v", err)
	}

	// Initialize the GCS client with explicit authentication.
	ctx := context.Background()
	client, err := storage.NewClient(ctx,
		option.WithAuthCredentialsJSON(option.ServiceAccount, data),
	)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}
	// Ensure the client is closed properly when the function returns.
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Warning: failed to close GCS client: %v", err)
		}
	}()

	// Section: Connectivity Test Execution
	// Attempt to list objects in the specified bucket.
	if err := listGCSObjects(ctx, client, cfg); err != nil {
		return fmt.Errorf("failed to list GCS objects: %w", err)
	}

	fmt.Println("GCS connectivity test successful.")
	return nil
}

// listGCSObjects attempts to retrieve and list objects from the target bucket.
// This serves as the primary verification step for GCS connectivity and permissions.
func listGCSObjects(ctx context.Context, client *storage.Client, cfg config) error {
	bucket := client.Bucket(cfg.bucketName)

	// Verify bucket existence and access by fetching attributes.
	if _, err := bucket.Attrs(ctx); err != nil {
		return fmt.Errorf("failed to get bucket attributes: %w", err)
	}

	// Inform the user about the listing parameters.
	if cfg.bucketPrefix == "" {
		fmt.Printf("Listing up to %d objects in bucket %s\n", cfg.maxObjects, cfg.bucketName)
	} else {
		fmt.Printf("Listing up to %d objects in bucket %s with prefix %s:\n", cfg.maxObjects, cfg.bucketName, cfg.bucketPrefix)
	}

	// Iterate through objects in the bucket, applying prefix filtering if provided.
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

		// Print the name of the object found.
		fmt.Println(attrs.Name)
		count++

		// Limit the number of objects returned to prevent excessive output.
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
