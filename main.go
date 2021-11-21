package main

import (
	"crypto/rand"
	"crypto/rsa"
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
	// Generate RSA key
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	// Instanciate our Signer and Unsigner with our keypair
	priv, _ := newSignerFromKey(key)
	pub, _ := newUnsignerFromKey(key.Public())
	// Sign message
	msg := []byte("Hello, world!")
	signature, err := priv.Sign(msg)
	if err != nil {
		panic(err)
	}
	logrus.Info("Signed signature: ", signature)
	// Verify signature
	err = pub.Unsign(msg, signature)
	if err != nil {
		logrus.Warn("Unverified signature")
		panic(err)
	}
	logrus.Info("Verified signature")

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
