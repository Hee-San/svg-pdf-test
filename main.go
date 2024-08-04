package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	imageData, err := os.ReadFile("test_image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// 元の画像サイズを取得
	originalImageSize := float64(len(imageData)) / (1024 * 1024) // MB単位

	// テストする画像の呼び出し回数
	imageCounts := []int{1, 2, 5, 10, 20, 50, 100}

	fmt.Printf("Original Image Size: %.4f MB\n\n", originalImageSize)
	fmt.Println("Image Count | PDF Size (MB)")
	fmt.Println("------------|---------------")

	for _, count := range imageCounts {
		pdfSize, err := generatePDF(base64Image, count)
		if err != nil {
			log.Printf("Error generating PDF for %d images: %v", count, err)
			continue
		}
		pdfSizeMB := float64(pdfSize) / (1024 * 1024) // Convert to MB
		fmt.Printf("%-12d| %.4f\n", count, pdfSizeMB)
	}
}

func generatePDF(base64Image string, imageCount int) (int64, error) {
	chromePath := os.Getenv("CHROME_PATH")
	if chromePath == "" {
		return 0, fmt.Errorf("CHROME_PATH environment variable is not set")
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	html := generateHTML(base64Image, imageCount)

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			return err
		}),
	); err != nil {
		return 0, err
	}

	filename := fmt.Sprintf("output_%d.pdf", imageCount)
	if err := os.WriteFile(filename, buf, 0644); err != nil {
		return 0, err
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

func generateHTML(base64Image string, count int) string {
	svgElements := ""
	for i := 0; i < count; i++ {
		svgElements += `
    <svg width="100" height="100">
      <use href="#testImage" />
    </svg>`
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>SVG Image Test</title>
</head>
<body>
    <svg width="0" height="0" style="display: none;">
      <defs>
        <image id="testImage" width="100" height="100" href="data:image/jpeg;base64,%s"/>
      </defs>
    </svg>

    <h1>SVG Image Test (%d images)</h1>

    %s
</body>
</html>
`, base64Image, count, svgElements)
}
