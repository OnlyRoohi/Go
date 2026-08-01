/*
 * TgMusicBot - Thumbnail utility
 * Minimal modern thumbnail system: fetches remote thumbnail if available
 * or generates a simple placeholder image. Caching and advanced drawing
 * can be added later.
 */

package core

import (
	"crypto/sha1"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// GenerateThumbnail tries to download a thumbnail from thumbURL and saves it to downloads/thumbnails.
// If thumbURL is empty or download fails, it generates a simple placeholder image.
// Returns the path to the saved thumbnail file.
func GenerateThumbnail(title, duration, requester, channel, thumbURL string) (string, error) {
	cacheDir := filepath.Join("database", "thumbnails")
	_ = os.MkdirAll(cacheDir, 0755)

	// create a deterministic filename based on title+channel
	h := sha1.New()
	_, _ = h.Write([]byte(title + "|" + channel))
	name := fmt.Sprintf("thumb_%x.jpg", h.Sum(nil))
	outPath := filepath.Join(cacheDir, name)

	// If already exists and not older than 24h, return it
	if fi, err := os.Stat(outPath); err == nil {
		if time.Since(fi.ModTime()) < 24*time.Hour {
			return outPath, nil
		}
	}

	// Try to download if thumbURL provided
	if thumbURL != "" {
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Get(thumbURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			f, err := os.Create(outPath)
			if err == nil {
				_, _ = io.Copy(f, resp.Body)
				f.Close()
				return outPath, nil
			}
		}
	}

	// Fallback: generate simple image
	const W, H = 1280, 720
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	// gradient background
	for y := 0; y < H; y++ {
		c := uint8(30 + (y*60)/H)
		col := color.RGBA{c, uint8(20 + (y*30)/H), uint8(60 + (y*80)/H), 255}
		for x := 0; x < W; x++ {
			img.Set(x, y, col)
		}
	}

	// simple overlay box
	over := image.NewRGBA(image.Rect(W-420, H-160, W-40, H-40))
	draw.Draw(img, over.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 120}}, image.Point{}, draw.Over)

	f, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	opts := &jpeg.Options{Quality: 80}
	if err := jpeg.Encode(f, img, opts); err != nil {
		return "", err
	}

	return outPath, nil
}
