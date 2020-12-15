// Canny edge detector implementation in golang
// https://en.wikipedia.org/wiki/Canny_edge_detector

/**
* To-do:
* - Load image into the the script and create a 2d array
* - Convert the image color to grey scale
* - apply gaussian filter to smooth the image
* - Find the intensity gradient of the image
*   - Sobel operator
* - Edge thinning -> non-maximum suppression
* - Remove weak gradients -> Double threshold
* - Apply histeresis
* - Export the image with the edges
 */

package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
)

func loadImage(url string) image.Image {

	fmt.Println("> Loading Image")

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

	fmt.Println("< Done")

	return imRgb
}

func rgbToGreyscale(img image.Image) image.Image {

	fmt.Println("> Converting gray scale...")
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

	fmt.Println("< Done")
	return grayImg
}

func imageToTensor(img image.Image) (*[][]color.Gray, image.Point) {

	fmt.Println("> Image to Tensor...")

	size := img.Bounds().Size()
	var pixels [][]color.Gray

	//put pixels into two dimensional array
	for j := 0; j < size.Y; j++ {
		var y []color.Gray
		for i := 0; i < size.X; i++ {
			p := color.GrayModel.Convert(img.At(j, i)).(color.Gray)
			y = append(y, p)
		}
		pixels = append(pixels, y)
	}

	fmt.Println("< Done")

	return &pixels, size
}

func tensorToImage(pixels [][]color.Gray) image.Image {

	fmt.Println("> Tensor to Image...")

	rect := image.Rect(0, 0, len(pixels), len(pixels[0]))
	newImage := image.NewGray(rect)

	for x := 0; x < len(pixels); x++ {
		for y := 0; y < len(pixels[0]); y++ {
			original, ok := color.GrayModel.Convert(pixels[x][y]).(color.Gray)
			if ok {
				newImage.Set(x, y, original)
			}
		}
	}

	fmt.Println("< Done")
	return newImage
}

func convolveKernel(img [][]color.Gray, kernel [][]float64) *[][]color.Gray {

	ker := kernel
	rows := len(img)
	columns := len(img[0])

	filtered := make([][]color.Gray, rows)
	for i := 0; i < columns; i++ {
		filtered[i] = make([]color.Gray, columns)
	}

	var newValue float64

	// convolve kernel over image

	min := (len(kernel) / 2)
	max := (len(kernel) / 2) + 1

	for y := min; y < rows-max; y++ {
		for x := min; x < columns-max; x++ {
			newValue = 0.
			for i := -len(kernel) / 2; i < max; i++ {
				for j := -len(kernel) / 2; j < max; j++ {
					pixel := img[y-i][x-j].Y
					newValue += float64(pixel) * ker[i+min][j+min]
				}
			}

			// save new pixel
			filtered[y+min][x+min].Y = uint8(newValue)
		}
	}

	return &filtered

}

func getGaussianKernel(size int, sigma float64) ([][]uint32, float64) {

	fmt.Println("> Generating gaussian filter...")

	// https://homepages.inf.ed.ac.uk/rbf/HIPR2/gsmooth.htm
	k_size := size

	kern1d := make([]float64, k_size)
	kern2d := make([][]float64, k_size)
	gaussian_filter := make([][]uint32, k_size)

	// initialize matrices
	for i := range kern2d {
		kern2d[i] = make([]float64, k_size)
		gaussian_filter[i] = make([]uint32, k_size)
	}

	// Calculate 1-D Gaussian distribution
	two_sigma_sq := 2 * math.Pow(sigma, 2)
	calc1 := 1.0 / (math.Sqrt(2*math.Pi) * sigma)

	for i := -size / 2; i < (size/2)+1; i++ {
		numerator := math.Pow(float64(i), 2)
		kern1d[i+(size/2)] = calc1 * math.Exp(-(numerator / two_sigma_sq))
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
	scalar := 1.0 / kern2d[0][0]
	for i := range kern1d {
		for j := range kern1d {
			gaussian_filter[i][j] = uint32(math.Floor((kern2d[i][j] / kern2d_csum) * scalar))
		}
	}

	fmt.Println("< Done")

	return gaussian_filter, scalar
}

func applyGaussuianFilter(size image.Point, oldImg [][]color.Gray, kernel *[][]uint32, k_scalar float64) *[][]color.Gray {

	fmt.Println("> Applying gaussian filter...")

	nKer := *kernel
	var k_upper int = (len(nKer) / 2) + 1
	var k_lower int = len(nKer) / 2

	newImg := make([][]color.Gray, size.Y)
	for i := range newImg {
		newImg[i] = make([]color.Gray, size.X)
	}

	// Convolve filter mask over image
	for y := k_lower; y < len(newImg)-k_lower; y++ {
		for x := k_lower; x < len(newImg[y])-k_lower; x++ {

			var sum int = 0
			newPixelValue := color.Gray{}

			// iterate over kernel
			for i := -k_lower; i < k_upper; i++ {
				for j := -k_lower; j < k_upper; j++ {
					pixel := oldImg[y+i][x+j].Y
					sum += int(pixel) * int(nKer[i+k_lower][j+k_lower])
				}
			}
			// calculate sum average
			newPixelValue.Y = uint8(sum / int(k_scalar))
			newImg[y][x] = newPixelValue
		}
	}

	fmt.Println("< Done")

	return &newImg

}

func applySobelGradients(img [][]color.Gray) *[][]color.Gray {

	fmt.Println("> Applying Sobel Gradient...")

	gx := [3][3]float64{{-1., 0., 1.}, {-2., 0., 2.}, {-1., 0., 1.}}
	gy := [3][3]float64{{-1., -2., -1.}, {0., 0., 0.}, {1., 2., 1.}}
	sharpen_kernel := [][]float64{{0, -1, 0}, {-1, 5, -1}, {0, -1, 0}}

	// sharpen the image for better gradient results
	fmt.Println(">> Applying Sharpening Kernel...")
	img = *convolveKernel(img, sharpen_kernel)
	fmt.Println("<< Done")

	threshold := 255 * 0.3
	rows := len(img)
	columns := len(img[0])

	magnitude := make([][]color.Gray, rows)
	for i := 0; i < columns; i++ {
		magnitude[i] = make([]color.Gray, columns)
	}

	var sx, sy float64

	// convolve sobel kernel over image
	for y := 1; y < rows-2; y++ {
		for x := 1; x < columns-2; x++ {

			sx = 0.
			sy = 0.

			// convolve kernel over image slice
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {

					pixel := img[y-i][x-j].Y

					sx += float64(pixel) * gx[i+1][j+1]
					sy += float64(pixel) * gy[i+1][j+1]

				}
			}

			// remove neg. values
			sx = math.Abs(sx)
			sy = math.Abs(sy)

			gradient := math.Ceil(math.Sqrt(math.Pow(sx, 2) + math.Pow(sy, 2)))

			// apply threshold
			if math.Max(gradient, threshold) == threshold {
				gradient = 0.
			}

			// save new pixel
			magnitude[y+1][x+1].Y = uint8(gradient)
		}
	}

	fmt.Println("< Done")

	return &magnitude

}

func nonMaximumSuppression() {}

func doubleThreadhold() {}

func applyHistersis() {}

func exportImage(img image.Image, dest string, filename string, encoding string) {

	fmt.Println("> Generating image...")

	newImage, err := os.Create(dest + "/" + filename + "." + encoding)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer newImage.Close()

	var img_err error

	if encoding == "png" {
		img_err = png.Encode(newImage, img)
	} else if encoding == "jpg" || encoding == "jpeg" {
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

	fmt.Println("< Image genetarted Successfully!")
}

func main() {

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

	fmt.Println("--- Script initialized! ---")

	rgb_image := loadImage(input_filename)
	gray_image := rgbToGreyscale(rgb_image)
	tensor, size := imageToTensor(gray_image)
	kenrel, k_scalar := getGaussianKernel(5, 2.5)
	filtered := applyGaussuianFilter(size, *tensor, &kenrel, k_scalar)
	sobel := applySobelGradients(*filtered)

	new_image := tensorToImage(*sobel)

	exportImage(new_image, "output", output_filename, extension)

	fmt.Println("====================================")
	fmt.Println(">>> Script executed successfully <<<")
	fmt.Println(">> Input file:", input_filename)
	fmt.Printf(">> Ouput file: output/%v.%v\n", output_filename, extension)
	fmt.Println("====================================")

}
