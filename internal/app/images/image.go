package images

import (
	e "error"
	"fmt"
	"generateId"
	"github.com/disintegration/imaging"
	"image"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"user"
)

func init() {
	//Ensure our goroutines run across all cores
	runtime.GOMAXPROCS(runtime.NumCPU())

}

const imageIDLength = 10

type Image struct {
	ID          string
	UserID      string
	Name        string
	Location    string
	Size        int64
	CreatedAt   time.Time
	Description string
}

func NewImage(user *user.User) *Image {
	return &Image{
		ID:        generateId.GenerateID("img", imageIDLength),
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
}

// A map of accepted mime types and their file extension
var mimeExtensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

func (image *Image) CreateFromURL(imageURL string) error {
	// Get the response from the URL
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return e.ErrImageURLInvalid
	}

	defer response.Body.Close()

	// Ascertain the type of file we downloaded
	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return e.ErrInvalidImageType
	}

	// Get an extension for the file
	ext, valid := mimeExtensions[mimeType]
	fmt.Println("ext_CreateFromURl_image.go: ", ext)
	if !valid {
		return e.ErrInvalidImageType
	}

	// Get a name from the URL
	image.Name = filepath.Base(imageURL)
	image.Location = image.ID + ext

	// Open a file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()

	// Copy the entire response to the output file
	size, err := io.Copy(savedFile, response.Body)
	if err != nil {
		return err
	}

	// The returned value from io.Copy is the number of bytes copied
	image.Size = size

	// Create the various resizes of the images
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}

	// Save our image to the store
	return GlobalImageStore.Save(image)
}

func (image *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	// Move our file to an appropriate place, with an appropriate name
	image.Name = headers.Filename
	image.Location = image.ID + filepath.Ext(image.Name)

	// Open a file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	defer savedFile.Close()
	// Copy the uploaded file to the target location
	size, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}
	image.Size = size

	// Create the various resizes of the images
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}

	// Save the image to the database
	return GlobalImageStore.Save(image)
}

func (image *Image) StaticRoute() string {
	return "/im/" + image.Location
}

func (image *Image) ShowRoute() string {
	return "/image/" + image.ID
}
func (image *Image) StaticThumbnailRoute() string {
	return "/im/thumbnail/" + image.Location
}

func (image *Image) StaticPreviewRoute() string {
	return "/im/preview/" + image.Location
}
func (image *Image) CreateResizedImages() error {
	//Generate an image from file
	srcImage, err := imaging.Open("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	//Create a channel to receive errors on
	errorChan := make(chan error)

	//Process each size
	// Process each size
	go image.resizePreview(errorChan, srcImage)
	go image.resizeThumbnail(errorChan, srcImage)
	//Wait for images to finish resizing
	for i := 0; i < 2; i++ {
		err := <- errorChan
		if err != nil {
			return err
		}
	}

	return nil
}

var widthThumbnail = 400

func (image *Image) resizeThumbnail(errorChan chan error, srcImage image.Image)  {
	dstImage := imaging.Thumbnail(srcImage, widthThumbnail, widthThumbnail, imaging.Lanczos)

	destination := "./data/images/thumbnail/" + image.Location
	errorChan <- imaging.Save(dstImage, destination)
}

var widthPreview = 800

func (image *Image) resizePreview(errorChan chan error, srcImage image.Image) {
	size := srcImage.Bounds().Size()
	ratio := float64(size.Y) / float64(size.X)
	targetHeight := int(float64(widthPreview) * ratio)

	dstImage := imaging.Resize(srcImage, widthPreview, targetHeight, imaging.Lanczos)

	destination := "./data/images/preview/" + image.Location
	errorChan <- imaging.Save(dstImage, destination)
}
