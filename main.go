package main

import (
	"fmt"
	"image"
	"image/color"
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
	var priv Signer
	var pub Unsigner

	if len(os.Args) <= 3 {
		logrus.Info("Generate key")
		priv, pub, err = InitializeKeyPair(4096, true)
		if err != nil {
			panic(err)
		}
	} else {
		logrus.Info("Load key")
		priv, pub, err = InitializeKeyPairFromFiles(os.Args[3], os.Args[4])
		if err != nil {
			panic(err)
		}
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

	// Generate diplome
	// Load template
	logrus.Info("Load template")
	template := LoadImage("template.png")

	// Add coucou text
	logrus.Info("Add coucou text")
	template.AddCenteredText("coucou", 24.0, color.RGBA{0, 0, 0, 255}, template.Bounds().Dy()/2)
	logrus.Info("Save coucou text")
	template.Save("test_template.png")
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
