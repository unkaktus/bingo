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
	"os/exec"
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

func pixelate(src image.Image, name string, n_pixels int) (string, error) {
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y

	if width != height {
		return "", fmt.Errorf("non-square image")
	}

	if n_pixels == 0 {
		n_pixels = width
	}

	g := gift.New(
		gift.Pixelate(width / n_pixels),
	)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	filename := fmt.Sprintf("%s_pixelated_%04d.png", name, n_pixels)
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := png.Encode(f, dst); err != nil {
		return "", fmt.Errorf("encode png: %w", err)
	}
	return filename, nil
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

	for i := 1; i < 12; i++ {
		pixel_config = append(pixel_config, int(math.Pow(1.5, float64(i))))
	}
	// Append the original
	pixel_config = append(pixel_config, 0)

	imageFilenames := []string{}

	for _, n_pixels := range pixel_config {
		imageFilename, err := pixelate(src, name, n_pixels)
		if err != nil {
			return fmt.Errorf("pixelate: %w", err)
		}
		imageFilenames = append(imageFilenames, imageFilename)
	}

	pdfFilename := fmt.Sprintf("%s_pixelated.pdf", name)
	cmd := exec.Command("magick", append(imageFilenames, pdfFilename)...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run magick command: %w", err)
	}

	for _, imageFilename := range imageFilenames {
		if err := os.Remove(imageFilename); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
