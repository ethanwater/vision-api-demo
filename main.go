package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

var chunkSize = 2048

func IsValidImage(imagePath string) bool {
	switch {
  case strings.HasPrefix(imagePath, "."):
    return false
	case
		strings.HasSuffix(imagePath, ".jpg"),
		strings.HasSuffix(imagePath, ".jpeg"),
		strings.HasSuffix(imagePath, ".png"),
		strings.HasSuffix(imagePath, ".gif"),
		strings.HasSuffix(imagePath, ".bmp"):
		return true
  default:
    return false
	}
}

func ObtainImageLabels(ctx context.Context, imagePath string) {
	f, err := os.Open(fmt.Sprintf("bucket/%s", imagePath))
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer f.Close()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer client.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	entityAnnotations, err := client.DetectLabels(ctx, image, &pb.ImageContext{}, 10)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	for _, entity := range entityAnnotations {
		fmt.Println(entity)
	}
}

func main() {
	ctx := context.Background()

	bucket, err := os.ReadDir("bucket")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	for _, object := range bucket {
    if IsValidImage(object.Name()) {
		  fmt.Println(object.Name())
		  ObtainImageLabels(ctx, object.Name())
    }
		fmt.Println()
	}

}
