package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func displayHelp() {
	fmt.Println("Usage: go run . <COMMAND> <ARGS>")
	fmt.Println("Commands:")
	fmt.Println("\tgenerate-custom-diplome <NAME> <GRADE> <IMG_TO_HIDE_PATH> <KEYPAIR_BIT_SIZE>")
	fmt.Println("\textract-lsb-from-diplome <HID_IMG_PATH>")
}

func main() {
	if len(os.Args) < 2 {
		displayHelp()
		return
	}

	switch os.Args[1] {
	case "generate-custom-diplome":
		if len(os.Args) != 6 {
			displayHelp()
			return
		}
		bitSize, err := strconv.Atoi(os.Args[5])
		if err != nil {
			logrus.Fatal(err)
		}
		GenerateCustomDiplome(os.Args[2], os.Args[3], os.Args[4], bitSize)
		logrus.Info("Custom diplome successfully generated! Saved to generated_diplome.png")
	case "extract-lsb-from-diplome":
		if len(os.Args) != 3 {
			displayHelp()
			return
		}
		ExtractLSBFromDiplome(os.Args[2])
		logrus.Info("Hidden data successfully extracted! Saved to extracted_lsb.png")
	default:
		displayHelp()
	}

	// img, err := openImage(os.Args[1])
	// if err != nil {
	// 	panic(err)
	// }

	// logrus.Info("LSB Steganography")
	// imgWithMsg := EncodeLSBSteganography(imageToNRGBA(img), []byte(os.Args[2]))
	// err = saveToImg("test.png", imgWithMsg)
	// if err != nil {
	// 	panic(err)
	// }
	// decodedMdg := DecodeLSBSteganography(&imgWithMsg, len([]byte(os.Args[2]))*8)
	// fmt.Println(string(decodedMdg))

	// logrus.Info("Sign file with random key")
	// var priv Signer
	// var pub Unsigner

	// if len(os.Args) <= 3 {
	// 	logrus.Info("Generate key")
	// 	priv, pub, err = InitializeKeyPair(4096, true)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// } else {
	// 	logrus.Info("Load key")
	// 	priv, pub, err = InitializeKeyPairFromFiles(os.Args[3], os.Args[4])
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// // Sign message
	// msg := []byte("Hello, world!")
	// signature, err := priv.Sign(msg)
	// if err != nil {
	// 	panic(err)
	// }
	// logrus.Info("Signed signature: ", signature)

	// // Verify signature
	// err = pub.Verify(msg, signature)
	// if err != nil {
	// 	logrus.Warn("Unverified signature")
	// 	panic(err)
	// }
	// logrus.Info("Verified signature")

	// // Verify signature with wrong one
	// err = pub.Verify(msg, []byte("wrong signature"))
	// if err == nil {
	// 	logrus.Warn("Verified signature")
	// 	panic(nil)
	// }
	// logrus.Info("Unverified signature")

	// // Generate diplome
	// // Load template
	// logrus.Info("Load template")
	// template := LoadImage("template.png")

	// // Add coucou text
	// logrus.Info("Add coucou text")
	// template.AddCenteredText("Diplôme", 62.0, color.RGBA{255, 65, 65, 255}, 150)
	// template.AddCenteredText("Master Informatique", 42.0, color.RGBA{255, 65, 65, 255}, 240)
	// template.AddCenteredText("Désservi à M. Hugo Bollon", 24.0, color.RGBA{255, 65, 65, 255}, 320)
	// logrus.Info("Save coucou text")
	// template.Save("test_template.png")
}

func openImage(path string) (image.Image, error) {
	f, err := os.Open(path)
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
	file, err := os.Create(path)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	if err := png.Encode(file, &img); err != nil {
		logrus.Fatal(err)
	}

	return nil
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
