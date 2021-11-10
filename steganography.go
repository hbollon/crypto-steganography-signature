package main

import (
	"image"
)

func EncodeLSBSteganography(img *image.NRGBA, msg []byte) image.NRGBA {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	output := image.NewNRGBA(image.Rect(0, 0, width, height))
	var msgIdx int
	var bitIdx int
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if msgIdx < len(msg) {
				pix := img.NRGBAAt(j, i)
				setLSB(&pix.R, getBitFromByte(msg[msgIdx], bitIdx))
				output.Set(j, i, pix)

				bitIdx++
				if bitIdx == 8 {
					bitIdx = 0
					msgIdx++
				}
			} else {
				output.Set(j, i, img.At(j, i))
			}
		}
	}
	return *output
}

func DecodeLSBSteganography(img *image.NRGBA, nbBits int) []byte {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	msgSize := nbBits
	var msg []byte
	var bitIdx int
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			pix := img.NRGBAAt(j, i)
			bit := getLSB(pix.R)
			msg = append(msg, bit)
			bitIdx++
			if bitIdx == msgSize {
				return msg
			}
		}
	}
	return msg
}

func getBitFromByte(b byte, indexInByte int) byte {
	b = b << uint(indexInByte)
	var mask byte = 0x80

	bit := mask & b
	if bit == 128 {
		return 1
	}
	return 0
}

func getLSB(b byte) byte {
	if b%2 == 0 {
		return 0
	}
	return 1
}

func setLSB(b *byte, bit byte) {
	if bit == 1 {
		*b = *b | 1
	} else if bit == 0 {
		var mask byte = 0xFE
		*b = *b & mask
	}
}
