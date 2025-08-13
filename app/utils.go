package main

import(
	"os"
	"log"
	"bytes"
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

func createMosaic(ratio float64, img image.Image) image.Image {
	//Get Size
	size := img.Bounds()
	w, h := float64(size.Max.X), float64(size.Max.Y)
	
	//Mosaic
	down := imaging.Resize(img, int(w*ratio), int(h*ratio), imaging.NearestNeighbor)
	return imaging.Resize(down, int(w), int(h), imaging.NearestNeighbor)
}

func saveImage(ext string, img image.Image) {
	//Ready File for save
	f, err := os.Create("./mosaic" + ext)
	if err != nil {
    	log.Fatal(err)
	}
	defer f.Close()
	
	//Save 
	var saveErr error
	if ext == "jpeg" {
		saveErr = jpeg.Encode(f, img, nil);
	} else {
		saveErr = png.Encode(f, img);
	}
	if saveErr != nil {
		log.Fatal(saveErr)
	}
}