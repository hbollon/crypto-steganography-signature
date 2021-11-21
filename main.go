package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	img, err := openImage(os.Args[1])
	if err != nil {
		panic(err)
	}

	logrus.Info("LSB Steganography")
	imgWithMsg := EncodeLSBSteganography(imageToNRGBA(img), []byte(os.Args[2]))
	err = saveToImg("test.png", imgWithMsg)
	if err != nil {
		panic(err)
	}
	decodedMdg := DecodeLSBSteganography(&imgWithMsg, len([]byte(os.Args[2]))*8)
	fmt.Println(string(decodedMdg))

	logrus.Info("Sign file with random key")
	priv, pub, err := InitializeKeyPair(4096)
	if err != nil {
		panic(err)
	}

	// Sign message
	msg := []byte("Hello, world!")
	signature, err := priv.Sign(msg)
	if err != nil {
		panic(err)
	}
	logrus.Info("Signed signature: ", signature)

	// Verify signature
	err = pub.Verify(msg, signature)
	if err != nil {
		logrus.Warn("Unverified signature")
		panic(err)
	}
	logrus.Info("Verified signature")

	// Verify signature with wrong one
	err = pub.Verify(msg, []byte("wrong signature"))
	if err == nil {
		logrus.Warn("Verified signature")
		panic(nil)
	}
	logrus.Info("Unverified signature")

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
