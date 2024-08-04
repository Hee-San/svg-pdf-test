# SVG PDF テスト

このプロジェクトは、SVG画像参照の数とPDFファイルサイズの関係を分析し、同時に`chromedp/chromedp`ライブラリの挙動をテストします。

## 実験結果

`chromedp/chromedp`を使用して、同じ画像をSVG参照で複数回埋め込む際のPDFファイルサイズを測定する実験を行いました。

元の画像サイズ: 10.2654 MB

| 画像の呼び出し回数 | PDFサイズ (MB) |
|------------------:|---------------:|
|                 1 |        10.2848 |
|                 2 |        10.2934 |
|                 5 |        10.2935 |
|                10 |        10.2937 |
|                20 |        10.2940 |
|                50 |        10.2944 |
|               100 |        10.2952 |

## 分析

1. 元の画像サイズは約10.27MBです。
2. 1回の画像参照で生成されたPDFは元の画像よりわずかに大きく、約10.28MBです。
3. 画像の参照回数を1から100に増やしても、PDFのサイズは10.2848MBから10.2952MBへとわずか0.0104MB（約10.6KB）しか増加していません。
4. これは、SVGの<use>要素を使用して画像を参照する方法が非常に効率的であることを示しています。
5. 画像データ自体は1回しか保存されず、追加の参照はファイルサイズにほとんど影響を与えません。

これらの結果は、同じ画像を多数回使用する必要がある大規模な文書や、ファイルサイズの制限が厳しい環境で、このアプローチが特に有用であることを示しています。

また、この実験を通じて`chromedp/chromedp`ライブラリがPDF生成においてSVG参照を効率的に処理できることが確認されました。

## SVGの`<use>`要素について

このプロジェクトでは、SVGの`<use>`要素を活用して画像の効率的な再利用を実現しています。

### `<use>`要素とは

SVGの`<use>`要素は、既に定義されたSVG要素を再利用するための機能です。これにより、同じ画像やグラフィック要素を複数回描画する際に、データを重複させることなく参照できます。

### 本プロジェクトでの使用例

```xml
<svg width="0" height="0" style="display: none;">
  <defs>
    <image id="testImage" width="100" height="100" href="data:image/jpeg;base64,..."/>
  </defs>
</svg>

<h1>SVG Image Test (5 images)</h1>

<svg width="100" height="100">
  <use href="#testImage" />
</svg>

<svg width="100" height="100">
  <use href="#testImage" />
</svg>

<svg width="100" height="100">
  <use href="#testImage" />
</svg>

<svg width="100" height="100">
  <use href="#testImage" />
</svg>

<svg width="100" height="100">
  <use href="#testImage" />
</svg>
```

この例では、`<defs>`内で画像を一度定義し、その後`<use>`要素で5回参照しています。実際のテストでは、この参照回数を1回から100回まで変化させて実験を行いました。
`<use>`要素を使用することで、画像データは1回だけ保存され、後はその参照のみが追加されます。そのため、参照回数が増えてもファイルサイズの増加を最小限に抑えることができます。
実験結果が示すように、この方法は画像の参照回数が大幅に増加しても、PDFサイズにはほとんど影響を与えません。

## セットアップと使用方法

1. Go 1.x以降をインストールします。
2. Chrome または Chromium ブラウザをインストールします。
3. このリポジトリをクローンします：
    ```
    git clone https://github.com/Hee-San/svg-pdf-test.git
    cd svg-pdf-test
    ```
4. 依存関係をインストールします：
    ```
    go mod tidy
    ```
5. 環境変数 `CHROME_PATH` を設定します。例：
    ```
    export CHROME_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
    ```
7. プログラムを実行します：
    ```
    go run main.go
    ```
7. コンソールに出力される結果を確認します。
8. 生成されたPDFファイルを確認します。

## 依存関係

- Go 1.x
- github.com/chromedp/chromedp
- Chrome または Chromium ブラウザ
