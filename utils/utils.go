package utils

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/simon0191/slack-visitor/model"
	s "github.com/simon0191/slack-visitor/shared"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var digits = []rune("0123456789")

func RandString(size int) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandDigits(size int) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

func LoadConfig(fileName string) (*model.Config, error) {

	fileName, err := filepath.Abs(fileName)

	if err != nil {
		return nil, s.NewError("utils.load_config.file_path_error", err, s.Options{"Filename": fileName})
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, s.NewError("utils.load_config.opening_file_error", err, s.Options{"Filename": fileName})
	}

	decoder := json.NewDecoder(file)
	config := model.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, s.NewError("utils.load_config.decoding_file_error", err, s.Options{"Filename": fileName})
	}

	return &config, nil
}
