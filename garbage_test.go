
//this file created automatically
package main
import (
	"bytes"
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func FuzzHex(f *testing.F) {
	for _, seed := range [][]byte{{}, {0}, {9}, {0xa}, {0xf}, {1, 2, 3, 4}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		enc := hex.EncodeToString(in)
		out, err := hex.DecodeString(enc)
		if err != nil {
			t.Fatalf("%v: decode: %v", in, err)
		}
		if !bytes.Equal(in, out) {
			t.Fatalf("%v: not equal after round trip: %v", in, out)
		}
	})
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{255, 0, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x), fixed.Int26_6(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func draw_the_dependances(all_func_code_block []function) {
	// rectangle := "dependances.png"

	// rectImage := image.NewRGBA(image.Rect(0, 0, 1080, 1080))
	// white := color.RGBA{255, 255, 255, 255}

	// draw.Draw(rectImage, rectImage.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	// file, err := os.Create(rectangle)
	// if err != nil {
	// 	log.Fatalf("failed create file: %s", err)
	// }
	// png.Encode(file, rectImage)
	img := image.NewRGBA(image.Rect(0, 0, 1080, 1080))
	white := color.RGBA{255, 255, 255, 255}
	addLabel(img, 540, 540, "Hello Go")
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	f, err := os.Create("dependances.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
