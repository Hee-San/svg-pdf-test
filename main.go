package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const htmlTemplate = `
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

    <h1>SVG Image Test</h1>

    <svg width="100" height="100">
      <use href="#testImage" />
    </svg>

    <svg width="100" height="100">
      <use href="#testImage" />
    </svg>
</body>
</html>
`

func main() {
	chromePath := os.Getenv("CHROME_PATH")
	if chromePath == "" {
		log.Fatal("CHROME_PATH environment variable is not set")
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// 画像ファイルを読み込み、base64エンコード
	imageData, err := ioutil.ReadFile("test_image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// HTMLを生成
	html := fmt.Sprintf(htmlTemplate, base64Image)

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
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			return ioutil.WriteFile("output.pdf", buf, 0644)
		}),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println("PDF generated successfully: output.pdf")
}
