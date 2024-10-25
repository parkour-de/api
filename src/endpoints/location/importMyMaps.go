package location

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"pkv/api/src/domain"
	"pkv/api/src/repository/t"
	"pkv/api/src/service/description"
	"strings"
	"time"
)

type KML struct {
	Document Document `xml:"Document"`
}

type Document struct {
	Name        string      `xml:"name"`
	Description string      `xml:"description"`
	Folders     []Folder    `xml:"Folder"`
	Placemarks  []Placemark `xml:"Placemark"`
}

type Folder struct {
	Name       string      `xml:"name"`
	Placemarks []Placemark `xml:"Placemark"`
	Folders    []Folder    `xml:"Folder"`
}

func (d Document) toFolder() Folder {
	return Folder{
		Name:       d.Name,
		Placemarks: d.Placemarks,
		Folders:    d.Folders,
	}
}

type Placemark struct {
	Name        string    `xml:"name"`
	Description string    `xml:"description"`
	Point       Point     `xml:"Point"`
	StyleURL    string    `xml:"styleUrl"`
	IconStyle   IconStyle `xml:"IconStyle"`
}

type Point struct {
	Coordinates string `xml:"coordinates"`
}

type IconStyle struct {
	IconHref string `xml:"Icon>href"`
}

func extractKML(data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	zipReader, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		// It's not a KMZ file, return the data as-is
		return data, nil
	}
	for _, file := range zipReader.File {
		if strings.HasSuffix(file.Name, ".kml") {
			kmlFile, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer kmlFile.Close()
			return io.ReadAll(kmlFile)
		}
	}
	return nil, t.Errorf("KML file not found in KMZ archive")
}

func processDocument(d Document, documentName string) {
	processFolder(d, d.toFolder(), documentName, []string{documentName})
}

func processFolder(d Document, folder Folder, documentName string, folderPath []string) {
	for _, p := range folder.Placemarks {
		location := processPlacemark(d, p, folderPath)
		fmt.Printf("Location: %+v\n", location)
	}

	for _, subfolder := range folder.Folders {
		processFolder(d, subfolder, documentName, append(folderPath, subfolder.Name))
	}
}

func processPlacemark(d Document, p Placemark, folderPath []string) domain.Location {
	coords := strings.Split(p.Point.Coordinates, ",")
	lat, lng := parseFloat(coords[1]), parseFloat(coords[0])
	placemarkDescription := domain.Descriptions{
		"de": {
			Title:  p.Name,
			Text:   p.Description,
			Render: description.Render([]byte(p.Description)),
		},
	}
	folderPathStr := ""
	if len(folderPath) > 0 {
		folderPathStr = strings.Join(folderPath, ";")
	}

	information := map[string]string{
		"importedFrom":                "mymaps",
		"importedId":                  p.Name + p.Point.Coordinates,
		"importedCategory":            folderPathStr,
		"importedStyle":               p.StyleURL,
		"importedDocumentName":        d.Name,
		"importedDocumentDescription": d.Description,
		"importedDocumentId":          "mymaps mid",
		"importedGeometry":            fmt.Sprintf("<Point>%s</Point>", p.Point.Coordinates),
	}

	return domain.Location{
		Entity: domain.Entity{
			Created:  time.Now().UTC(),
			Modified: time.Now().UTC(),
		},
		Lat:          lat,
		Lng:          lng,
		Type:         "spot",
		Information:  information,
		Descriptions: placemarkDescription,
		Photos:       processPhotos(p),
	}
}

func processPhotos(_ Placemark) domain.Photos {
	return domain.Photos{Photos: []domain.Photo{}}
}
