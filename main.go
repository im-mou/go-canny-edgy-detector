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
    "strings"
    "github.com/montanaflynn/stats"
)

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

    fmt.Println("< Done");

    return imRgb
}

func rgbToGreyscale(img image.Image) image.Image {

    fmt.Println("> Converting gray scale...");
    ok := true

    // Converting image to grayscale
    grayImg := image.NewGray(img.Bounds())

    for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
        for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
            // this line automatically converts RGBA -> Gray
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

    fmt.Println("< Done");
    return grayImg
}

func imageToTensor(img image.Image) (*[][]color.Gray, image.Point) {

    fmt.Println("> Image to Tensor...");

    size := img.Bounds().Size()
    var pixels [][]color.Gray
    //put pixels into two dimensional array
    for j:=0; j<size.Y;j++{
        var y []color.Gray
        for i:=0; i<size.X;i++{
            p := color.GrayModel.Convert(img.At(j,i)).(color.Gray)
            y = append(y, p)
        }
        pixels = append(pixels,y)
    }

    fmt.Println("< Done");

    return &pixels, size
}

func tensorToImage(pixels [][]color.Gray) image.Image {

    fmt.Println("> Tensor to Image...");

    rect := image.Rect(0,0,len(pixels),len(pixels[0]))
    newImage := image.NewGray(rect)

    for x:=0; x<len(pixels); x++{
        for y:=0; y<len(pixels[0]); y++ {
            original,ok := color.GrayModel.Convert(pixels[x][y]).(color.Gray)
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

func applyGaussuianFilter(size image.Point, oldImg [][]color.Gray, kernel *[][]uint32) *[][]color.Gray {

    fmt.Println("> Applying gaussian filter...");

    nKer := *kernel
    var k_upper int = (len(nKer)/2) + 1
    var k_lower int = len(nKer)/2

    newImg := make([][]color.Gray, size.Y)
    for i := range newImg {
        newImg[i] = make([]color.Gray, size.X)
    }

    // iterate over imge
    for y := k_lower; y < len(newImg) - k_lower; y++ {
        for x := k_lower; x < len(newImg[y]) - k_lower; x++ {

            var sum int = 0
            newPixelValue := color.Gray{}

            // iterate over kernel
            for i := -k_lower; i < k_upper; i++ {
                for j := -k_lower; j < k_upper; j++ {
                    pixel := oldImg[y + i][x + j].Y
                    sum += int(pixel) * int(nKer[i + k_lower][j + k_lower])
                }
            }

            newPixelValue.Y = uint8(sum/273)
            newImg[y][x] = newPixelValue
        }
    }

    fmt.Println("< Done");

    return &newImg

}

func getGradients() {}

func nonMaximumSuppression() {}

func doubleThreadhold() {}

func applyHistersis() {}


func exportImage(img image.Image, dest string, filename string, encoding string) {

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
    input_filename := arg[0]
    output := arg[1]

    allowed_outputs := [3]string{"jpg", "jpeg", "png"}

    // Split output string into array
    output_arg := strings.Split(output, ".")
    extension := output_arg[len(output_arg)-1]

    // get file name witout .ext
    output_filename_arr := output_arg[:len(output_arg)-1]
    output_filename := strings.Join(output_filename_arr[:], ".")

    valid_ext := false
    for _, ext := range allowed_outputs {
        if ext == extension {
            valid_ext = true
        }
    }

    if !valid_ext {
        panic("Output file extension not allowed, use one of the following -> [\"jpg\", \"png\"]")
    }

    fmt.Println("--- Script initialized! ---");

    rgb_image := loadImage(input_filename)
    gray_image := rgbToGreyscale(rgb_image)
    tensor, size := imageToTensor(gray_image)
    kenrel := getGaussianKernel(5, 2.5)
    filtered := applyGaussuianFilter(size, *tensor, &kenrel)
    blurred := tensorToImage(*filtered)

    exportImage(blurred, "output", output_filename, extension);

    fmt.Printf("\n");
    fmt.Println(">>> Script executed successfully <<<");
    fmt.Println(">> Input file:", input_filename);
    fmt.Printf(">> Ouput file: output/%v.%v\n", output_filename,extension);

}
