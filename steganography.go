package main

import (
	"image"
	"strconv"
)

// EncodeLSBSteganography encodes a message in an image using LSB steganography method
// Returns the image with the message hidden
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

// DecodeLSBSteganography decodes a message hidden in an image using LSB steganography method
// Returns the message
func DecodeLSBSteganography(img *image.NRGBA, nbBits int) []byte {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	msgSize := nbBits
	// Real message size is msgSize/8
	lenBytes := msgSize/8 + 1
	msg := make([]byte, lenBytes)
	var bitIdx int
	var bitBuffer string
	var bitBufSize int
	var idx int
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			pix := img.NRGBAAt(j, i)
			bit := getLSB(pix.R)
			bitIdx++
			if bitBufSize == 8 {
				// ParseInt interprets a string s in the given base (0, 2 to 36)
				// and bit size (0 to 64) and returns the corresponding value i.
				tmp, err := strconv.ParseInt(bitBuffer, 2, 64)
				if err != nil {
					panic(err)
				}
				msg[idx] = byte(tmp)
				idx++
				bitBuffer, bitBufSize = "", 0
			}
			bitBuffer = bitBuffer + strconv.Itoa(int(bit))
			bitBufSize++
			if bitIdx == msgSize {
				if bitBufSize != 0 {
					for i := 0; i < 8-bitBufSize; i++ {
						bitBuffer = bitBuffer + strconv.Itoa(0)
					}
					tmp, _ := strconv.ParseInt(bitBuffer, 2, 64)
					msg[idx] = byte(tmp)
					idx++
				}
				// Truncate the message to the actual size
				msg = msg[:idx]
				return msg
			}
		}
	}
	return msg
}

// getBitFromByte returns the bit at the given index in the given byte
func getBitFromByte(b byte, indexInByte int) byte {
	b = b << uint(indexInByte)
	var mask byte = 0x80

	bit := mask & b
	if bit == 128 {
		return 1
	}
	return 0
}

// getLSB returns the LSB of the given byte
func getLSB(b byte) byte {
	if b%2 == 0 {
		return 0
	}
	return 1
}

// setLSB sets the LSB of the given byte to the given value
func setLSB(b *byte, bit byte) {
	if bit == 1 {
		*b = *b | 1
	} else if bit == 0 {
		var mask byte = 0xFE
		*b = *b & mask
	}
}
