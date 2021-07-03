package localStorage

import (
	"fmt"
	"github.com/AliceDiNunno/go-image-database/src/config"
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type localStorage struct {
	config config.StorageConfig
}

func (l localStorage) makePath(id string) string {
	return fmt.Sprintf("%s/%s", l.config.Folder, id)
}

func (l localStorage) UploadFile(id string, file io.Reader) error {
	path := l.makePath(id)
	out, err := os.Create(path)
	if err != nil {
		logrus.Errorln(err)
		return domain.ErrUnableToUploadFile
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		logrus.Errorln(err)
		return domain.ErrUnableToUploadFile
	}
	return nil
}

func (l localStorage) DeleteFile(id string) error {
	path := l.makePath(id)
	err := os.Remove(path)

	if err != nil {
		return domain.ErrFileNotFound
	}

	return nil
}

func (l localStorage) GetFile(id string) (io.Reader, error) {
	path := l.makePath(id)
	file, err := os.Open(path)

	if err != nil {
		return nil, domain.ErrFileNotFound
	}

	return file, nil
}

func NewLocalStorage(config config.StorageConfig) localStorage {
	err := os.MkdirAll(config.Folder, os.ModePerm)

	if err != nil {
		logrus.Fatalln(err)
	}

	return localStorage{
		config: config,
	}
}
