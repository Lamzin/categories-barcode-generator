package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

const N = 1000

func main() {
	const pdfFileName = "out/Pallets.pdf"

	// Remove old file.
	os.Remove(pdfFileName)

	var files []string
	for i := 1; i <= N; i++ {
		fmt.Printf("Generating %d/%d barcode\n", i, N)

		text := fmt.Sprintf("Pallet-%d", i)
		DrawPNG(text)

		fileName := fmt.Sprintf("out/%s.png", text)
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
		if err := pdfcpu.ImportImagesFile(files[i:minInt(i+batch, len(files))], "out/Pallets.pdf", nil, nil); err != nil {
			panic(err)
		}
	}
}

func DrawPNG(text string) {
	const W = 1600
	const H = 1000
	const P = 150
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.DrawImage(ReadHeart(), 10, 10)
	dc.DrawImage(ReadHeart(), W-250-10, 10)
	dc.DrawImage(ReadHeart(), 10, H-250-10)
	dc.DrawImage(ReadHeart(), W-250-10, H-250-10)

	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("data/Arial Unicode.ttf", P); err != nil {
		panic(err)
	}
	dc.DrawStringWrapped(text, W/2, P, 0.5, 0, 0, 1.5, gg.AlignCenter)
	dc.DrawImage(GenerageBarCode(text), 300, 400)

	if err := dc.SavePNG(fmt.Sprintf("out/%s.png", text)); err != nil {
		panic(err)
	}
}

func GenerageBarCode(text string) image.Image {
	barCode, _ := code128.Encode(text)
	scaledBarCode, _ := barcode.Scale(barCode.(barcode.Barcode), 1000, 300)
	return scaledBarCode
}

func ReadHeart() image.Image {
	file, err := os.Open("data/pallet.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	scaled := resize.Resize(250, 0, img, resize.Lanczos3)
	return scaled
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
