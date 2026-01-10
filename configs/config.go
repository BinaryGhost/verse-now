package config

import (
	"log"
	"os"
	"path/filepath"
)

func personal_babij_source() string {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not use home-directory")
	}

	return filepath.Join(home_dir, "Projects", "json-bibles", "bible_collection_eo")
}

/*
The path, where it the babij_ressources are located.
For myself, i use this function, but you might specify a different path.

It is advised to clone from https://github.com/BinaryGhost/Json-Bibles
*/
var BABIJ_SOURCE = personal_babij_source()
