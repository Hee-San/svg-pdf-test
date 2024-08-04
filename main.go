package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

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
	// 画像ファイルを読み込み、base64エンコード
	imageData, err := ioutil.ReadFile("/app/test_image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// HTMLを生成
	html := fmt.Sprintf(htmlTemplate, base64Image)

	// ChromeDPのコンテキストを設定
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// PDFのバイトスライスを格納する変数
	var pdfBuffer []byte

	// ChromeDPを使用してHTMLをPDFに変換
	err = chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+html),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.PDF(&pdfBuffer, chromedp.PDFPrintBackground(true)),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 生成されたPDFをファイルに保存
	err = ioutil.WriteFile("/app/output.pdf", pdfBuffer, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PDF generated successfully: /app/output.pdf")
}
