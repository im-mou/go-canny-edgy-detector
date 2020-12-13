// Canny edge detector implementation in golang

/**
* To-do:
* - Load image into the the script and create a 2d array
* - Convert the image color to grey scale
* - apply gaussian filter to smooth the image
* - Find the intensity gradient of the image
* - Edge thinning -> non-maximum suppression
* - Remove weak gradients -> Double threshold
* - Apply histeresis
* - Export the image with the edges
*/

package main

import (
    "fmt"
    "os"
    "image"
    "log"
    "math"
    "image/color"
    "errors"
    "image/png"
    "image/jpeg"
)

// Pixel struct example
type Pixel struct {
    V int
}

func loadImage(url string) image.Image {

    fmt.Println("> Loading Image");

	//Decode the JPEG data. If reading from file, create a reader with
	reader, err := os.Open(url)
	if err != nil {
	    log.Fatal(err)
	}
    defer reader.Close()

    imRgb, _, err := image.Decode(reader)
	if err != nil {
        panic(err)
    }

    return imRgb
}

func rgbToGreyscale(img image.Image) image.Image {

    fmt.Println("> Converting gray scale...");
    ok := true

    // Converting image to grayscale
    grayImg := image.NewGray(img.Bounds())

    for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
        for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
            grayImg.Set(x, y, img.At(x, y))
        }
    }

    // check if grey image was created successfully
    if grayImg.ColorModel() != color.GrayModel {
        ok = false
    }

    if !ok {
        panic("Image was not converted to grayscale")
    }

    fmt.Println("< Conversion Successfull!");
    return grayImg
}

func imageTotensor(img image.Image) *[][]color.Color {
    size := img.Bounds().Size()
    var pixels [][]color.Color
    //put pixels into two dimensional array
    for i:=0; i<size.X;i++{
      var y []color.Color
      for j:=0; j<size.Y;j++{
         y = append(y,img.At(i,j))
      }
      pixels = append(pixels,y)
    }

    return &pixels
}

func convolution(img image.Image, encoding string) {

}

func applyGaussuian(pixels *[][]color.Color, sigma float64) {

    guassian_sum := 0.0

    // gaussian filter with 5 x 5 kernel
    k_size := int((4 * sigma * 0.5) + 1)
    // create filter empty slice
    gaussian_filter := make([][]float32, k_size)
    for i := range gaussian_filter {
        gaussian_filter[i] = make([]float32, k_size)
    }

    sig2 := 1.0
    sig22 := 2 * sig2
    x1 := (1/ (math.Pi * sig22))

    // gaussian continous sample matrix
    for y := 1; y < k_size+1; y++ {
        for x := 1; x < k_size+1; x++ {
            sum := math.Pow(float64(x - k_size+2), 2) + math.Pow(float64(y - k_size+2), 2)
            x2 := math.Exp(-(sum) / sig22)
            gaussian_filter[y-1][x-1] = float32(x1 * x2)
            guassian_sum += x1 * x2
        }
    }

    // Discretization-> gaussian matri
    for y := 0; y < k_size; y++ {
        for x := 0; x < k_size; x++ {
            gaussian_filter[y][x] /= float32(guassian_sum)
        }
    }

    // new image


}

func getGradients() {

}

func nonMaximumSuppression() {

}

func doubleThreadhold() {
}

func applyHistersis( ){

}

func exportImage(img image.Image, encoding string, dest string, filename string) {

    fmt.Println("> Generating image...");

    newImage, err := os.Create(dest+"/"+filename+"."+encoding)
    if err != nil {
        // handle error
        log.Fatal(err)
    }
    defer newImage.Close()

    var img_err error

    if encoding == "png" {
        img_err = png.Encode(newImage, img)
    } else if (encoding == "jpg" || encoding == "jpeg") {
        options := jpeg.Options{
            Quality: 100,
        }
        img_err = jpeg.Encode(newImage, img, &options)
    } else {
        img_err = errors.New("Image was not converted to grayscale")
    }

    if img_err != nil {
        panic("Image was not converted to grayscale")
    }

    fmt.Println("< Image genetarted Successfully!");
}


func main(){

    // parse arguments
    arg := os.Args[1:]
    filename := arg[0]


    fmt.Println("--- Script initialized! ---");

    rgb_image := loadImage(filename)
    gray_image := rgbToGreyscale(rgb_image)
    tensor := imageTotensor(gray_image)

    applyGaussuian(tensor, 2)


    // exportImage(gray_image, "jpdg", "output", "grey-image");


}
