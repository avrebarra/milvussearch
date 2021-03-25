package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Entity struct {
	ID     string     `json:"id"`
	Traits [100]Trait `json:"traits"`
}

func (e *Entity) FromFile(filename string) (err error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(bytes, e)
	if err != nil {
		return
	}

	return
}

type Trait struct {
	Name     string  `json:"name"`
	Tendency float32 `json:"tendency"`
}

// ***

func SHA1(text string) string {
	algorithm := sha1.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}
