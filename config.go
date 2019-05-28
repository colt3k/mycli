package mycli

import (
	"log"

	"github.com/pelletier/go-toml"
)

func LoadToml(path string) (*toml.Tree, error) {
	tree, err := toml.LoadFile(path)
	if err != nil {
		log.Fatalf("issue loading toml file\n%+v",err)
		return &toml.Tree{}, err
	}

	return tree, nil
}
