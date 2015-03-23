// theme_files
package main

import (
	"github.com/boltdb/bolt"
)

//go:generate go run buildassets.go

func InitThemes(t *bolt.Tx) {
	for name, datafunc := range ThemeTars {
		data, _ := datafunc()
		InsertTheme(t, name, data)
	}
}
