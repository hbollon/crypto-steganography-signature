package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Img is a composite type for an image.NRGBA
// Which allow us to keep original image.NRGBA methods and add new ones
type Img struct {
	*image.NRGBA
}

// GenerateCustomDiplome generates a custom diplome to generated_diplome.png
func GenerateCustomDiplome(nom, moyenne, photo string) {
	// Add diplome's custom text
	img := LoadImage("template.png")
	img.AddCenteredText("Diplôme", 62.0, color.RGBA{255, 65, 65, 255}, 150)
	img.AddCenteredText("Master Informatique", 42.0, color.RGBA{255, 65, 65, 255}, 240)
	img.AddCenteredText(fmt.Sprintf("Désservi à %s", nom), 24.0, color.RGBA{255, 65, 65, 255}, 320)
	img.AddCenteredText(fmt.Sprintf("Avec une moyenne de %s", moyenne), 18.0, color.RGBA{255, 65, 65, 255}, 400)
	img.Save("generated_diplome.png")

	// Hide photo using LSB
	diplome, err := openImage("generated_diplome.png")
	if err != nil {
		logrus.Fatal(err)
	}
	f, err := ioutil.ReadFile(photo)
	if err != nil {
		logrus.Fatal(err)
	}
	diplomeNRGBA := EncodeLSBSteganography(imageToNRGBA(diplome), f)
	err = saveToImg("generated_diplome.png", diplomeNRGBA)
	if err != nil {
		logrus.Fatal(err)
	}

	// Add signature
	f, err = ioutil.ReadFile(photo)
	if err != nil {
		logrus.Fatal(err)
	}
	signer, _, err := InitializeKeyPair(4096, true)
	if err != nil {
		logrus.Fatal(err)
	}
	_, err = signer.Sign(f)
	if err != nil {
		logrus.Fatal(err)
	}
	err = saveToImg("generated_diplome.png", diplomeNRGBA)
	if err != nil {
		logrus.Fatal(err)
	}

}

// LoadImage loads an image from a file and return a Img instance
// path must be a valid path to an png image file
func LoadImage(path string) Img {
	file, err := os.Open(path)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		logrus.Fatal(err)
	}

	return Img{img.(*image.NRGBA)}
}

// AddCenteredText adds text to the image centered on the image horitzontally
// Takes an imput text, a font size, a color and an y position as arguments
// Returns modified Img instance
func (img Img) AddCenteredText(text string, size float64, color color.RGBA, y int) Img {
	col := color
	fontRaw, err := ioutil.ReadFile("sans.ttf")
	if err != nil {
		logrus.Fatal(err)
	}
	fontType, _ := truetype.Parse(fontRaw)
	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(col),
		Face: truetype.NewFace(fontType, &truetype.Options{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
	}

	y = y + int(math.Ceil(size*72/72))
	d.Dot = fixed.Point26_6{
		X: (fixed.I(img.Bounds().Dx()) - d.MeasureString(text)) / 2,
		Y: fixed.I(y),
	}

	d.DrawString(text)
	return img
}

// Save saves the image to a file
func (img Img) Save(path string) {
	file, err := os.Create(path)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	if err := png.Encode(file, img.NRGBA); err != nil {
		logrus.Fatal(err)
	}
}
