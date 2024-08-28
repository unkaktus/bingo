package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/gift"
)

func readImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func pixelate(src image.Image, name string, n_pixels int) error {
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y

	if width != height {
		return fmt.Errorf("non-square image")
	}

	g := gift.New(
		gift.Pixelate(width / n_pixels),
	)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	filename := fmt.Sprintf("%s_pixelated_%04d.png", name, n_pixels)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, dst); err != nil {
		return fmt.Errorf("encode png: %w", err)
	}
	return nil
}

func run() error {
	flag.Parse()
	if len(flag.Args()) == 0 {
		return fmt.Errorf("image path is not specified")
	}
	filename := flag.Args()[0]
	name := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	src, err := readImage(filename)
	if err != nil {
		return fmt.Errorf("read image: %w", err)
	}

	pixel_config := []int{}

	for i := 0; i < 12; i++ {
		pixel_config = append(pixel_config, int(math.Pow(1.5, float64(i))))
	}

	for _, n_pixels := range pixel_config {
		if err := pixelate(src, name, n_pixels); err != nil {
			return fmt.Errorf("pixelate: %w", err)
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
