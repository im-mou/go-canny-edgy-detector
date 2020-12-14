# Canny edge detector Go implementation

A simple script to apply canny edge detector to png or jpeg images.
It uses gaussian blur to remove the noise and converts the image to Grayscale before applying the gaussian filter.

To run the script, execute the following commands into you CLI.
```bash
cd ~/go-canny-edgy-detector
go run ./main.go "images/input.<png|jpg|jpeg>" "output.<png|jpg|jpeg>"
```

The script takes a string argument to take the input filename with the path included and it outputs the resulting image into ./output folder.