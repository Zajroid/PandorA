package data

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

const (
	accountFile = "./data/account.bin"
)

func WriteAccountInfo(ecsID, password string) error {
	data := []byte(ecsID + ":" + password)

	file, err := os.Create(accountFile)
	if err != nil {
		return err
	}

	if err := binary.Write(file, binary.LittleEndian, rot47(data)); err != nil {
		return err
	}

	return nil
}

func ReadAccountInfo() (ecsID, password string, err error) {
	content := make([]byte, 0, 32)
	content, err = ioutil.ReadFile(accountFile)
	if err != nil {
		return
	}

	text := strings.Split(string(rot47(content)), ":")
	if len(text) != 2 {
		err = errors.New("Invalid format")
		return
	}

	ecsID, password = text[0], text[1]

	return
}

func rot47(target []byte) (result []byte) {
	result = make([]byte, len(target))

	for i, t := range target {
		result[i] = (t-33+47)%94 + 33
	}

	return
}