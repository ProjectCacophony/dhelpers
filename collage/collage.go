package collage

import (
	"bytes"
	"image"
	"strings"

	"context"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/opentracing/opentracing-go"
	"github.com/ungerik/go-cairo"
	"gitlab.com/Cacophony/dhelpers/cache"
	"gitlab.com/Cacophony/dhelpers/net"
)

// FromUrls creates a Collage PNG Image from internet image urls (PNG or JPEG).
// imageUrls		: a slice with all image URLs. Empty strings will create an empty space in the collage.
// descriptions		: a slice with text that will be written on each tile. Can be empty.
// width			: the width of the result collage image.
// height			: the height of the result collage image.
// tileWidth		: the width of each tile image.
// tileHeight		: the height of each tile image.
// backgroundColour	: the background colour as a hex string.
func FromUrls(ctx context.Context, imageUrls, descriptions []string, width, height, tileWidth, tileHeight int, backgroundColour string) (collageBytes []byte) {
	// start tracing span
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "dhelpers.collage.FromUrls")
	defer span.Finish()

	imageDataArray := make([][]byte, 0)
	// download images
	for _, imageURL := range imageUrls {
		if imageURL == "" {
			imageDataArray = append(imageDataArray, nil)
			continue
		}
		imageData, err := net.Get(imageURL)
		if err == nil {
			imageDataArray = append(imageDataArray, imageData)
		} else {
			cache.GetLogger().WithField("module", "collage").Errorln("error downloading", imageURL, err.Error())
			imageDataArray = append(imageDataArray, nil)
		}
	}

	return FromBytes(ctx, imageDataArray, descriptions, width, height, tileWidth, tileHeight, backgroundColour)
}

// FromBytes creates a Collage PNG Image from image []byte (PNG or JPEG).
// imageDataArray   : a slice of all image []byte data
// descriptions		: a slice with text that will be written on each tile. Can be empty.
// width			: the width of the result collage image.
// height			: the height of the result collage image.
// tileWidth		: the width of each tile image.
// tileHeight		: the height of each tile image.
// backgroundColour	: the background colour as a hex string.
func FromBytes(ctx context.Context, imageDataArray [][]byte, descriptions []string, width, height, tileWidth, tileHeight int, backgroundColour string) (collageBytes []byte) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "dhelpers.collage.FromBytes")
	defer span.Finish()

	// create surface with given background colour
	backgroundColourRGB, _ := colorful.Hex(backgroundColour) // nolint: errcheck, gas
	cairoSurface := cairo.NewSurface(cairo.FORMAT_RGB24, width, height)
	cairoSurface.SetSourceRGB(backgroundColourRGB.R, backgroundColourRGB.G, backgroundColourRGB.B)
	cairoSurface.Paint()

	var posX, posY int

	for i, imageData := range imageDataArray {
		// switch tile to new line if required
		if posX > 0 && posX+tileWidth > width {
			posY += tileHeight
			posX = 0
		}
		// draw image on tile if image exists
		if len(imageData) > 0 {
			tileImage, _, err := image.Decode(bytes.NewReader(imageData))
			if err != nil {
				cache.GetLogger().WithField("module", "collage").Errorln("error decoding tile image", err.Error())
			}
			if err == nil {
				tileSurface := cairo.NewSurfaceFromImage(tileImage)
				cairoSurface.SetSourceSurface(tileSurface, float64(posX), float64(posY))
				cairoSurface.Paint()
			}
		}
		// draw description on tile if description exists
		if len(descriptions) > i {
			// setup font and variables
			cairoSurface.SelectFontFace("UnDotum", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
			var offset, fontSize int
			// split description in lines
			lines := strings.Split(descriptions[i], "\n")
			for _, line := range lines {
				// clean line
				line = strings.TrimSpace(line)
				// reset font size
				fontSize = 28
				// adjust font size to fit tile
				for {
					// gather dimensions of line with current font size
					cairoSurface.SetFontSize(float64(fontSize))
					extend := cairoSurface.TextExtents(line)
					// break if line fits into tile, or font size is <= 10
					if extend.Width < float64(tileWidth)-6-6 || fontSize <= 10 {
						break
					}
					// try a smaller font
					fontSize--
				}
				// draw text
				cairoSurface.SetSourceRGB(1, 1, 1) // white
				cairoSurface.MoveTo(float64(posX+6), float64(posY+6+fontSize+offset))
				cairoSurface.ShowText(line)
				// draw white outline to improve readability
				cairoSurface.MoveTo(float64(posX+6), float64(posY+6+fontSize+offset))
				cairoSurface.TextPath(line)
				cairoSurface.SetSourceRGB(0, 0, 0) // black
				cairoSurface.SetLineWidth(4.5)
				cairoSurface.Stroke()
				// draw black outline to make text bold
				cairoSurface.MoveTo(float64(posX+6), float64(posY+6+fontSize+offset))
				cairoSurface.TextPath(line)
				cairoSurface.SetSourceRGB(1, 1, 1) // white
				cairoSurface.SetLineWidth(2.5)
				cairoSurface.Stroke()
				// switch to new line
				offset += fontSize + 6
			}
		}
		// switch to next tile
		posX += tileWidth
	}

	// write surface to byte slice and return it
	bytesData, _ := cairoSurface.WriteToPNGStream()
	return bytesData
}
