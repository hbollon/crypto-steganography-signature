package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	img, err := openImage(os.Args[1])
	if err != nil {
		panic(err)
	}
	imgWithMsg := EncodeLSBSteganography(imageToNRGBA(img), []byte(os.Args[2]))
	err = saveToImg("test.png", imgWithMsg)
	if err != nil {
		panic(err)
	}
	decodedMdg := DecodeLSBSteganography(&imgWithMsg, len([]byte(os.Args[2]))*8)
	fmt.Println(string(decodedMdg))
}

func openImage(path string) (image.Image, error) {
	f, err := os.Open(os.Args[1])
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return src, nil
}

func saveToImg(path string, img image.NRGBA) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, &img)
}

// imageToNRGBA converts image.Image to image.NRGBA
func imageToNRGBA(src image.Image) *image.NRGBA {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	m := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), src, bounds.Min, draw.Src)
	return m
}
