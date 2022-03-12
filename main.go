package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/fogleman/gg"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pkg/errors"
)

func main() {
	const pdfFileName = "out/_Всі лейбли для коробок--AllBoxLabels.pdf"

	// Remove old file.
	os.Remove(pdfFileName)

	records, err := ReadCSV()
	if err != nil {
		panic(errors.Wrap(err, "failed to read CSV"))
	}

	var files []string
	for i, raw := range records {
		if len(raw) < 3 {
			continue
		}

		// Ignore big category.
		if strings.Contains(raw[0], "00") {
			continue
		}

		// Skip empty raws.
		if len(raw[0]) == 0 {
			continue
		}

		fmt.Printf("Generating %d/%d barcode\n", i, len(records))

		var (
			barcodeText = raw[0]
			uaText      = raw[1]
			deText      = raw[2]
		)

		fileName := DrawPNG(barcodeText, uaText, deText)

		defer func() {
			fmt.Println("Removing " + fileName)
			if err := os.Remove(fileName); err != nil {
				panic(err)
			}
		}()

		files = append(files, fileName)
	}

	fmt.Println("Concatenating PNG files. This might take up to few minutes")
	const batch = 100
	for i := 0; i < len(files); i += batch {
		fmt.Printf("Merging %d/%d PNG files\n", i, len(files))
		if err := pdfcpu.ImportImagesFile(files[i:minInt(i+batch, len(files))], pdfFileName, nil, nil); err != nil {
			panic(err)
		}
	}

	for i := 0; i < len(files); i++ {
		fmt.Printf("Generating %d/%d separate PDF files\n", i, len(files))

		name := strings.ReplaceAll(files[i], ".png", ".pdf")
		// Remove old file.
		os.Remove(name)
		if err := pdfcpu.ImportImagesFile([]string{files[i], files[i]}, name, nil, nil); err != nil {
			panic(err)
		}
	}
}

func ReadCSV() ([][]string, error) {
	file, err := os.Open("data/category_list.csv")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open CSV file")
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CSV records")
	}
	return records, nil
}

func DrawPNG(barcodeText, uaText, deText string) string {
	const W = 1600
	const H = 1000
	// const P = 100
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)

	const margin = 30

	uaTextP := 120.0
	if utf8.RuneCountInString(uaText) < 15 {
		uaTextP = 200.0
	}
	if err := dc.LoadFontFace("data/Arial Unicode.ttf", uaTextP); err != nil {
		panic(err)
	}
	dc.DrawStringWrapped(uaText, W/2, margin, 0.5, 0, W-100, 1.5, gg.AlignCenter)

	deTextP := 90.0
	if utf8.RuneCountInString(deText) < 15 {
		deTextP = 150.0
	}
	if err := dc.LoadFontFace("data/Arial Unicode.ttf", float64(deTextP)); err != nil {
		panic(err)
	}
	dc.DrawStringWrapped(deText, W/2, 350, 0.5, 0, W-100, 1.5, gg.AlignCenter)

	dc.DrawImage(GenerageBarCode(barcodeText), 300, 550)

	const barcodeTextP = 100
	if err := dc.LoadFontFace("data/Arial Unicode.ttf", barcodeTextP); err != nil {
		panic(err)
	}
	dc.DrawStringWrapped(barcodeText, W/2, 870, 0.5, 0, W-100, 1.5, gg.AlignCenter)

	fileName := fileName(barcodeText, uaText, deText)
	if err := dc.SavePNG(fileName); err != nil {
		panic(err)
	}
	return fileName
}

func GenerageBarCode(text string) image.Image {
	barCode, _ := code128.Encode(text)
	scaledBarCode, _ := barcode.Scale(barCode.(barcode.Barcode), 1000, 300)
	return scaledBarCode
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fileName(barcodeText, uaText, deText string) string {
	name := fmt.Sprintf("%s--%s--%s.png", barcodeText, uaText, deText)
	name = strings.ReplaceAll(name, "/", "")
	return fmt.Sprintf("out/%s", name)
}
