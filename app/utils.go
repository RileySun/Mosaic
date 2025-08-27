package main

import(
	"log"
	"math"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	
	_ "embed"
	
	"image"
	"image/png"
	"image/jpeg"
	"image/color"
	
	"github.com/disintegration/imaging"
)

//go:embed DEFAULT_IMAGE.png
var DEFAULT_IMG []byte

var WHITE color.NRGBA = color.NRGBA{R: 155, G: 155, B: 155, A: 255}

func getExt(path string) string {
	return filepath.Ext(path)
}

func getDefaultImage() image.Image {
	img, _, imgErr := image.Decode(bytes.NewReader(DEFAULT_IMG))
	if imgErr != nil {
		log.Fatalln(imgErr)
	}
	
	return img
}

func getImage(path string) image.Image {
	imgByte, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	
	img, _, imgErr := image.Decode(bytes.NewReader(imgByte))
	if imgErr != nil {
		log.Fatalln(imgErr)
	}
	
	return img
}

func previewImage(img image.Image) image.Image {
	//Get Size
	size := img.Bounds()
	w, h := float64(size.Max.X), float64(size.Max.Y)
	
	//Get Ratio
	rawRatio = 1.0
	if w > 300 || h > 300 {
		wRatio := 1.0
		if w > 300 {
			wRatio = w/300
		}
		
		hRatio := 1.0
		if h/wRatio > 300 {
			hRatio = h/rawRatio/300
		}
		
		rawRatio = wRatio*hRatio
	}
	
	//log.Println("Preview Image", w/rawRatio, h/rawRatio, rawRatio)
	
	//Resize
	return imaging.Resize(img, int(w/rawRatio), int(h/rawRatio), mosaicMethod)
}

func createMosaic(ratio float64, img image.Image) image.Image {
	//Get Size
	size := img.Bounds()
	w, h := float64(size.Max.X), float64(size.Max.Y)
	
	//Mosaic
	down := imaging.Resize(img, int(w*ratio), int(h*ratio), mosaicMethod)
	return imaging.Resize(down, int(w), int(h), mosaicMethod)
}

func saveImage(w io.Writer, img image.Image) error {
	switch ext {
		case ".jpeg":
			return jpeg.Encode(w, img, nil)
		case ".jpg":
			return jpeg.Encode(w, img, nil)
		case ".png":
			return png.Encode(w, img)
		default:
			return errors.New("Invalid extension: " + ext)
	}
}

func roundFloat(input float64) float64 {
	return math.Floor(input * 100000)/100000
}

//Mosaic Method
var mosaicMethod imaging.ResampleFilter = imaging.NearestNeighbor
