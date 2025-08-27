# GCS Connection Tester

## Overview and Objectives

This is a simple command-line tool to test connectivity to Google Cloud Storage (GCS). It allows you to verify your credentials and list objects within a specified GCS bucket. The main objective of this tool is to provide a quick and easy way to ensure that your application can connect to GCS and access the required buckets and objects.

## Command-Line Syntax

The tool is configured using command-line flags:

| Flag          | Description                                | Required | Default |
|---------------|--------------------------------------------|----------|---------|
| `-credentials`| Path to the JSON credential file.          | Yes      |         |
| `-bucket`     | The name of the GCS bucket.                | Yes      |         |
| `-project`    | The ID of your GCP project.                | Yes      |         |
| `-prefix`     | An optional prefix to filter objects.      | No       |         |
| `-max`        | The maximum number of objects to list.     | No       | `10`    |
| `-version`    | Displays the application version.          | No       | `false` |

## Examples of Usage

### Build the application

First, build the application using the following command:

```bash
go build
```

### Basic Usage

To list the first 10 objects in a bucket, provide your credentials file, bucket name, and project ID:

```bash
./gcsconntest -credentials <path-to-credentials.json> -bucket <your-bucket-name> -project <your-project-id>
```

### List objects with a prefix

You can filter the objects by providing a prefix. This is useful for listing objects in a specific "folder":

```bash
./gcsconntest -credentials <path-to-credentials.json> -bucket <your-bucket-name> -project <your-project-id> -prefix "my-folder/"
```

### Limit the number of objects

To change the maximum number of objects listed, use the `-max` flag:

```bash
./gcsconntest -credentials <path-to-credentials.json> -bucket <your-bucket-name> -project <your-project-id> -max 5
```

### Display the version

To see the version of the tool, use the `-version` flag:

```bash
./gcsconntest -version
```
