package image

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

type Image struct {
	path string
}

func New(imagesPath string) (*Image, error) {
	const op = "storage.images.New"

	if _, err := os.Stat(imagesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: images dir does not exist: %w", op, err)
	}

	return &Image{path: imagesPath}, nil
}

func (i *Image) Save(img []byte, extension string) (string, error) {
	const op = "storage.images.Save"

	imagePath := i.generateImagePath(extension)
	fullPath := path.Join(i.path, imagePath)

	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("%s: can't create image: %w", op, err)
	}
	defer f.Close()

	_, err = f.Write(img)
	if err != nil {
		return "", fmt.Errorf("%s: can't write image byte: %w", op, err)
	}

	return imagePath, nil
}

func (i *Image) Get(imagePath string) ([]byte, error) {
	const op = "storage.images.Get"

	fullPath := path.Join(i.path, imagePath)
	b, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("%s: can't read file: %w", op, err)
	}

	return b, nil
}

func (i *Image) generateImagePath(extension string) string {
	now := time.Now()
	y, m, d := now.Date()
	year := strconv.Itoa(y)
	month := m.String()
	date := strconv.Itoa(d)

	dir := path.Join(year, month, date)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		i.createDir(dir)
	}

	name := generateRandomName(extension)

	return path.Join(dir, name)
}

func (i *Image) createDir(dir string) {
	_ = os.MkdirAll(path.Join(i.path, dir), os.ModePerm)
}

// TODO: Обработка расширения
func generateRandomName(extension string) string {
	name := ""
	for i := 0; i < 15; i++ {
		c := rune(rand.Intn('z'-'a') + 'a')
		name += string(c)
	}

	return name + "." + extension
}
