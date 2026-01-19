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
	PresetBlur    = ImagePreset{Name: "blur", Width: 20, Quality: 50}
	PresetAvatar  = ImagePreset{Name: "avatar", Width: 120, Quality: 80}
	PresetXSmall  = ImagePreset{Name: "xs", Width: 320, Quality: 75}
	PresetSmall   = ImagePreset{Name: "sm", Width: 640, Quality: 75}
	PresetMedium  = ImagePreset{Name: "md", Width: 768, Quality: 80}
	PresetLarge   = ImagePreset{Name: "lg", Width: 1024, Quality: 80}
	PresetXLarge  = ImagePreset{Name: "xl", Width: 1280, Quality: 85}
	Preset2XLarge = ImagePreset{Name: "2xl", Width: 1536, Quality: 85}
	PresetFullHD  = ImagePreset{Name: "fhd", Width: 1920, Quality: 90}
)

var AllPresets = []ImagePreset{
	PresetBlur,
	PresetAvatar,
	PresetXSmall,
	PresetSmall,
	PresetMedium,
	PresetLarge,
	PresetXLarge,
	Preset2XLarge,
	PresetFullHD,
}

type ProcessCallback func(presetName string, data []byte) error

func ProcessAndHandleImage(r io.Reader, presets []ImagePreset, onProcessed ProcessCallback) error {
	srcImage, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	originalBounds := srcImage.Bounds()

	var buf bytes.Buffer
	for _, p := range presets {
		targetWidth := p.Width
		if originalBounds.Dx() < targetWidth {
			targetWidth = originalBounds.Dx()
		}

		finalImage := imaging.Resize(srcImage, targetWidth, 0, imaging.Lanczos)

		buf.Reset()
		err = webp.Encode(&buf, finalImage, &webp.Options{
			Lossless: false,
			Quality:  p.Quality,
		})

		if err != nil {
			return err
		}

		if uploadErr := onProcessed(p.Name, buf.Bytes()); uploadErr != nil {
			return uploadErr
		}

		finalImage = nil
	}

	return nil
}
