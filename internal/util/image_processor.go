package util

import (
	"bytes"
	"image"
	"io"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

type ImagePreset struct {
	Name    string
	Width   int
	Quality float32
}

var (
	PresetThumb   = ImagePreset{Name: "thumb", Width: 400, Quality: 75}
	PresetMobile  = ImagePreset{Name: "mobile", Width: 768, Quality: 80}
	PresetDesktop = ImagePreset{Name: "desktop", Width: 1200, Quality: 85}
	PresetLarge   = ImagePreset{Name: "large", Width: 1920, Quality: 90}
)

type ProcessedImages map[string][]byte

func ProcessImage(r io.Reader, presets []ImagePreset) (ProcessedImages, error) {
	srcImage, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	results := make(ProcessedImages)
	originalBounds := srcImage.Bounds()

	for _, preset := range presets {
		var finalImage image.Image

		targetWidth := preset.Width

		if originalBounds.Dx() < targetWidth {
			targetWidth = originalBounds.Dx()
		}

		finalImage = imaging.Resize(srcImage, targetWidth, 0, imaging.Lanczos)

		var buf bytes.Buffer
		err = webp.Encode(&buf, finalImage, &webp.Options{
			Lossless: false,
			Quality:  preset.Quality,
		})

		if err != nil {
			return nil, err
		}

		results[preset.Name] = buf.Bytes()
	}
	return results, nil
}
