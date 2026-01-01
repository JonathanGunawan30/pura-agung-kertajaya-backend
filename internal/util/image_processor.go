package util

import (
	"bytes"
	"image"
	"io"
	"sync"

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

type ProcessedImages map[string][]byte

func ProcessImage(r io.Reader, presets []ImagePreset) (ProcessedImages, error) {
	srcImage, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	results := make(ProcessedImages)
	originalBounds := srcImage.Bounds()

	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, 3)

	var firstErr error
	var errOnce sync.Once

	for _, preset := range presets {
		wg.Add(1)

		p := preset

		go func() {
			defer wg.Done()

			sem <- struct{}{}

			defer func() { <-sem }()

			if firstErr != nil {
				return
			}

			targetWidth := p.Width

			if originalBounds.Dx() < targetWidth {
				targetWidth = originalBounds.Dx()
			}

			finalImage := imaging.Resize(srcImage, targetWidth, 0, imaging.Lanczos)

			var buf bytes.Buffer
			err = webp.Encode(&buf, finalImage, &webp.Options{
				Lossless: false,
				Quality:  p.Quality,
			})

			if err != nil {
				errOnce.Do(func() {
					firstErr = err
				})
				return
			}

			mu.Lock()
			results[p.Name] = buf.Bytes()
			mu.Unlock()

		}()
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return results, nil
}
