package main

import(	
	"log"
	"strconv"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
)

//Globals
var mosaicRatio float64 = 0.05
var ext string = ".png"
var window fyne.Window
var borderLine *canvas.Line
var inputImage, outputImage *canvas.Image

//Init
func init() {
	borderLine = canvas.NewLine(WHITE)
	
	//Input Image
	inputImage = canvas.NewImageFromImage(getDefaultImage())
	inputImage.FillMode = canvas.ImageFillContain
	inputImage.SetMinSize(fyne.NewSize(200, 200))
	
	//OutputImage
	outputImage = canvas.NewImageFromImage(getDefaultImage())
	outputImage.FillMode = canvas.ImageFillContain
	outputImage.SetMinSize(fyne.NewSize(200, 200))
	mosaicEffect()
}

//Main
func main() {
	myApp := app.New()
	window = myApp.NewWindow("Mosaic Maker v0.1")
	
	
	content := mainScreen()
	
	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 400))
	window.ShowAndRun()
}

//Renders
func mainScreen() *fyne.Container {
	//Main Content
	input := inputScreen()
	settings := settingsScreen()
	output := outputScreen()
	
	//Container
	content := container.NewHBox(input, settings, output)
	
	//Title
	title := widget.NewLabel("Mosaic Maker")
	titleContent := container.NewVBox(title, borderLine)
	
	return container.NewBorder(titleContent, nil, nil, nil, content)
}

func inputScreen() *fyne.Container {
	label := widget.NewLabel("Input Image")
	
	button := widget.NewButton("Select Image", func() {
		dialog.ShowFileOpen(onSelectImage , window)
	})
	
	content := container.NewVBox(label, button, inputImage)
	
	return container.NewHBox(content)
}

func settingsScreen() *fyne.Container {
	label := widget.NewLabel("Mosaic Settings")
	
	//Slider
	slider := widget.NewSlider(0.0001, 0.25)
	slider.Step = 0.00001
	slider.Value = mosaicRatio
	
	//Entry
	entry :=  widget.NewEntry()
	entry.Text = strconv.FormatFloat(mosaicRatio, 'f', -1, 64)
	entry.Refresh()
	
	//Onchange
	slider.OnChanged = func(newRatio float64) {
		mosaicRatio = newRatio
		entry.Text = strconv.FormatFloat(mosaicRatio, 'f', -1, 64)
		entry.Refresh()
		mosaicEffect()
	}
	entry.OnChanged = func(newStringRatio string) {
		newRatio, err := strconv.ParseFloat(newStringRatio, 64)
		if err != nil {
			return
		}
		mosaicRatio := newRatio
		slider.Value = mosaicRatio
		mosaicEffect()
	}
	
	content := container.NewVBox(label, slider, entry)
	
	return container.NewHBox(content)
}

func outputScreen() *fyne.Container {
	label := widget.NewLabel("Output Image")
	
	button := widget.NewButton("Save Image", openSaveFile)
	
	content := container.NewVBox(label, button, outputImage)
	
	return container.NewHBox(content)
}

//Actions
func onSelectImage(reader fyne.URIReadCloser, err error) {
	//Error checking
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	if reader == nil {
		return
	}
	
	//Change input image
	ext = getExt(reader.URI().Path())
	inputImage.Image = getImage(reader.URI().Path())
	inputImage.Refresh()
	
	//Mosaic and set Output
	mosaicEffect()
}

func mosaicEffect() {
	outputImage.Image = createMosaic(mosaicRatio, inputImage.Image)
	outputImage.Refresh()
}

func openSaveFile() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if writer == nil {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		defer writer.Close()
		
		ext := getExt(writer.URI().Extension())
		if ext != "png" || ext != "jpg" || ext != "jpeg" {
			info := dialog.NewInformation("Invalid Extension", "Saved file must use a valid image extension. (Example: png, jpeg, jpg)", window)
			info.SetOnClosed(openSaveFile)
			info.Show()
			return
		} else {
			err = saveImage(writer, outputImage.Image)	
			if err != nil {
				log.Println(err)
			}
		}
	}, window)
}