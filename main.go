package main

import(
	"os"
	"log"
	"bytes"
	"strconv"
	"io/ioutil"
	"path/filepath"
	
	"image"
	"image/jpeg"
	"image/png"
	
	"github.com/disintegration/imaging"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("Invalid Syntax: ./mosaic mosaicRatio imagePath")
	}
	
	//Get Settings
	ratio, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		log.Fatal("mosaicRatio is not a valid float64")
	}
	src := getImage(os.Args[2])
	ext := filepath.Ext(os.Args[2])
	
	//Set up output
	size := src.Bounds()
	w, h := float64(size.Max.X), float64(size.Max.Y)
	
	//Mosaic
	down := imaging.Resize(src, int(w*ratio), int(h*ratio), imaging.NearestNeighbor)
	output := imaging.Resize(down, int(w), int(h), imaging.NearestNeighbor)
	
	//Ready File for save
	f, err := os.Create("./mosaic" + ext)
	if err != nil {
    	log.Fatal(err)
	}
	defer f.Close()
	
	//Save 
	var saveErr error
	if ext == "jpeg" {
		saveErr = jpeg.Encode(f, output, nil);
	} else {
		saveErr = png.Encode(f, output);
	}
	if saveErr != nil {
		log.Fatal(saveErr)
	}
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