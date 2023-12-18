package main

import (
	"context"
	"fmt"
	"os"
	"strings"
  "sync"

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

func ObtainImageLabels(ctx context.Context, imagePath string, wg *sync.WaitGroup) {
  defer wg.Done()
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

  fmt.Printf("%s\n", imagePath)
	for _, entity := range entityAnnotations {
		fmt.Println(entity)
	}
  fmt.Println()
}

func main() {
	ctx := context.Background()
  var wg sync.WaitGroup

	bucket, err := os.ReadDir("bucket")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	for _, entry := range bucket {
    if IsValidImage(entry.Name()) {
      wg.Add(1)
      go ObtainImageLabels(ctx, entry.Name(), &wg)
    }
	}

  wg.Wait()
}
