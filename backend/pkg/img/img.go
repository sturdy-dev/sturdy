package img

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"github.com/disintegration/imaging"
)

func Thumbnail(size int, fp io.ReadSeeker, output io.Writer) error {
	buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	if _, err := fp.Read(buff); err != nil {
		return err
	}

	// Seek to start
	_, err := fp.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	var img image.Image
	switch ct := http.DetectContentType(buff); ct {
	case "image/jpeg":
		img, err = jpeg.Decode(fp)
	case "image/png":
		img, err = png.Decode(fp)
	default:
		return fmt.Errorf("unexpected content type %s", ct)
	}
	if err != nil {
		return err
	}

	thumb := imaging.Fill(img, size, size, imaging.Center, imaging.Lanczos)
	err = png.Encode(output, thumb)
	if err != nil {
		return err
	}

	return nil
}
