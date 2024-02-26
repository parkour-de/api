package photo

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"pkv/api/src/domain"
	"pkv/api/src/repository/dpv"
	"slices"
	"strconv"
	"strings"
	"time"
)

var PyVipsFiles = []string{
	".gif", ".jpg", ".jpe", ".jpeg", ".jfif",
	".jxl", ".png", ".webp", ".pdf", ".svg", ".ppm",
	".tif", ".tiff", ".heic", ".heif", ".avif",
	".mat", ".v", ".vips", ".img", ".hdr",
	".pbm", ".pgm", ".ppm", ".pfm", ".pnm",
	".svg", ".svgz", ".svg.gz",
	".j2k", ".jp2", ".jpt", ".j2c", ".jpc",
	".fits", ".fit", ".fts",
	".exr", ".svs", ".vms", ".vmu", ".ndpi", ".scn", ".mrxs", ".svslide", ".bif",
	".bpg", ".bmp", ".dib", ".dcm", ".emf",
}

// Upload takes binary data and a filename and stores it as a new file in the tmp folder
// The tmp folder is found via dpv.ConfigInstance.Server.TmpPath
// Files should be named randomBase64String.o.jxl
// The original filename is used to determine the original file format
// It would call s.Encode(...) with the binary data and the original extension to convert it into a JXL file
// Encode will also return the width and height
// Upload function will save the encoded file in the right location and thus returns a domain.Photo struct
func (s *Service) Upload(data []byte, filename string, ctx context.Context) (domain.Photo, error) {
	ext := filepath.Ext(filename)
	if !slices.Contains(PyVipsFiles, ext) {
		return domain.Photo{}, fmt.Errorf("unsupported image format: %s", ext)
	}
	// save original file in a more temporary temp location
	tmpFile, err := os.CreateTemp(os.TempDir(), "upload-*"+ext)
	if err != nil {
		return domain.Photo{}, fmt.Errorf("could not create temporary file before conversion: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		return domain.Photo{}, fmt.Errorf("could not save uploaded file before conversion: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return domain.Photo{}, fmt.Errorf("could not close uploaded file before conversion: %w", err)
	}
	randomFilename := RandomString() + ".o.jxl"
	photo, err := PythonConvert(tmpFile.Name(), dpv.ConfigInstance.Server.TmpPath+randomFilename)
	if err != nil {
		return domain.Photo{}, fmt.Errorf("could not convert uploaded file: %w", err)
	}
	return photo, nil
}

type PythonInput struct {
	InputFile   string `json:"input_file"`
	OutputFile  string `json:"output_file"`
	ModTime     int64  `json:"mod_time"`
	Orientation int    `json:"orientation"`
}

type PythonOutput struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Color  string `json:"color"`
}

func PythonConvert(inFile string, outFile string) (domain.Photo, error) {
	img, err := Info(inFile)
	if err != nil {
		return domain.Photo{}, fmt.Errorf("empty image information: %w", err)
	}
	if img.Orientation > 4 && img.ExifH > img.ExifW {
		img.Orientation = 1
	}

	input := PythonInput{
		InputFile:   inFile,
		OutputFile:  outFile,
		ModTime:     img.Date,
		Orientation: img.Orientation,
	}
	jsonData, err := json.Marshal(input)
	if err != nil {
		return domain.Photo{}, fmt.Errorf("marshaling image info for image \"%v\" failed: %w", inFile, err)
	}
	cmd := exec.Command("python", "image_processor.py")
	cmd.Stdin = bytes.NewReader(jsonData)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return domain.Photo{}, fmt.Errorf("could not start python process for image \"%v\": %w", inFile, err)
	}
	if err := cmd.Wait(); err != nil {
		return domain.Photo{}, fmt.Errorf("python process exited with error for image \"%v\": %w", inFile, err)
	}
	var result PythonOutput
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		return domain.Photo{}, fmt.Errorf("error decoding python result for image \"%v\": %w", inFile, err)
	}
	img.Width = result.Width
	img.Height = result.Height
	img.Color = result.Color
	return domain.Photo{
		Src:   outFile,
		W:     result.Width,
		H:     result.Height,
		Lat:   img.Lat,
		Lon:   img.Lon,
		Color: result.Color,
	}, nil
}

type ImgInfo struct {
	Width       int     `json:"w"`           // width
	Height      int     `json:"h"`           // height
	Date        int64   `json:"d,omitempty"` // EXIF date
	Color       string  `json:"c,omitempty"` // string representing a 4x4 low-res mesh gradient
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Orientation int     `json:"-"`
	ExifW       int     `json:"-"`
	ExifH       int     `json:"-"`
}

func Info(filename string) (img ImgInfo, err error) {
	img.Date = time.Now().Unix()
	cmd := exec.Command("exiftool", "-T", "-datetimeoriginal", "-orientation", "-gps:GPSLatitude", "-gps:GPSLongitude", "-imagewidth", "-imageheight", "-n", filename)
	out, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("creating pipe for \"exiftool\" with \"%v\" failed: %w", filename, err)
		return
	}
	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("executing \"exiftool\" with \"%v\" failed: %w", filename, err)
		return
	}
	defer out.Close()
	b, err := io.ReadAll(out)
	defer cmd.Wait()
	if err != nil {
		err = fmt.Errorf("reading from pipe of \"exiftool\" with \"%v\" failed: %w", filename, err)
		return
	}
	data := strings.Split(strings.Trim(string(b), " \r\n"), "\t")
	layout := "2006:01:02 15:04:05"
	date, err := time.ParseInLocation(layout, data[0], time.Local)
	if err == nil {
		img.Date = date.Unix()
	}
	if len(data) < 2 { // Assume error "No matching files"
		err = fmt.Errorf("no matching files")
		return
	}

	img.Orientation = orientation(data[1])
	img.Lat = ddStr(data[2])
	img.Lon = ddStr(data[3])
	img.ExifW, _ = strconv.Atoi(data[4])
	img.ExifH, _ = strconv.Atoi(data[5])
	return
}

func orientation(in string) int {
	val, err := strconv.Atoi(in)
	if err != nil {
		return 1
	}
	return val
}
func ddStr(in string) float64 {
	val, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0
	}
	return val
}

func RandomString() string {
	buff := make([]byte, 6)
	_, err := rand.Read(buff)
	if err != nil {
		println(fmt.Errorf("random number generation failed: %w", err).Error())
	}
	return base64.RawURLEncoding.EncodeToString(buff)
}
