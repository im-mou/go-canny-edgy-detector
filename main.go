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
    // "math"
    "image/color"
    "errors"
    "image/png"
    "image/jpeg"
    "github.com/montanaflynn/stats"
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

func imageTotensor(img image.Image) (*[][]color.Color, image.Point) {

    fmt.Println("> Image to Tensor...");

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

    fmt.Println("< Done");

    return &pixels, size
}

func tensorToImage(pixels [][]color.Color) image.Image {

    fmt.Println("> Tensor to Image...");

    rect := image.Rect(0,0,len(pixels),len(pixels[0]))
    newImage := image.NewRGBA(rect)

    for x:=0; x<len(pixels); x++{
        for y:=0; y<len(pixels[0]); y++ {
            q:=pixels[x]
            if q==nil{
                continue
            }
            p := pixels[x][y]
            if p==nil{
                continue
            }
            original,ok := color.RGBAModel.Convert(p).(color.RGBA)
            if ok{
                newImage.Set(x,y,original)
            }
        }
    }

    fmt.Println("< Done");
    return newImage
}

func getGaussianKernel(size int, sigma float64) [][]uint32 {

    fmt.Println("> Generating gaussian filter...");

    // https://stackoverflow.com/questions/29731726/how-to-calculate-a-gaussian-kernel-matrix-efficiently-in-numpy
    k_size := size
    delta := (sigma * 2) / float64(k_size)
    linspace := make([]float64, k_size+1)
    cfd := make([]float64, len(linspace))

    kern1d := make([]float64, len(linspace)-1)
    kern2d := make([][]float64, len(kern1d))
    gaussian_kernel := make([][]uint32, len(kern1d))

    for i := range kern2d {
        kern2d[i] = make([]float64, len(kern1d))
        gaussian_kernel[i] = make([]uint32, len(gaussian_kernel))
    }

    // Normal cumulative distribution function
    for i := range linspace {
        linspace[i] = -sigma + delta * float64(i)
        cfd[i] = stats.NormCdf(linspace[i], 0, delta)
    }

    for i := 0; i < len(linspace)-1; i++ {
        kern1d[i] = cfd[i + 1] - cfd[i]
    }

    //outer product
    kern2d_csum := 0.0
    for i := range kern1d {
        for j := range kern1d {
            mult := kern1d[i] * kern1d[j]
            kern2d[i][j] = mult
            kern2d_csum += mult
        }
    }

    // normalize
    for i := range kern1d {
        for j := range kern1d {
            gaussian_kernel[i][j] = uint32((kern2d[i][j] / kern2d_csum) * 273)
        }
    }

    fmt.Println("< Done");

    return gaussian_kernel
}

func applyGaussuianFilter(size image.Point, oldImg *[][]color.Color, kernel *[][]uint32) *[][]color.Color {

    fmt.Println("> Applying gaussian filter...");

    nKer := *kernel
    // pad := int(len(nKer) / 2)
    newImg := make([][]color.Color, size.Y)
    for i := range newImg {
        newImg[i] = make([]color.Color, size.X)
    }

    copy(newImg, *oldImg)

    // iterate over imge
    for y := len(nKer); y < len(newImg) - len(nKer); y++ {
        for x := len(nKer); x < len(newImg[y]) - len(nKer); x++ {

            newPixelValue := color.RGBA{}

            // iterate over kernel
            for i := range nKer {
                for j := range nKer {
                    r,g,b,a := newImg[y + i][x + j].RGBA()

                    newPixelValue.R += uint8(r * nKer[i][j])
                    newPixelValue.G += uint8(g * nKer[i][j])
                    newPixelValue.B += uint8(b * nKer[i][j])
                    newPixelValue.A += uint8(a * nKer[i][j])
                }
            }

            newImg[y][x] = newPixelValue

        }
    }

    fmt.Println("< filter applied");

    return &newImg

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
    tensor, size := imageTotensor(gray_image)
    kenrel := getGaussianKernel(5, 2.5)
    filtered := applyGaussuianFilter(size, tensor, &kenrel)
    blured := tensorToImage(*filtered)

    exportImage(blured, "jpg", "output", "blureed-image");

}
