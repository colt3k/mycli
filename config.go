package mycli

import (
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"strings"
	"sync"
)

var once sync.Once
var (
	instance *TomlWrapper
)

func Toml() *TomlWrapper {
	once.Do(func() { // <-- atomic, does not allow repeating
		instance = &TomlWrapper{} // <-- thread safe
	})

	return instance
}

type TomlWrapper struct {
	Map map[string]interface{}
}

func (t *TomlWrapper) LoadToml(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, &t.Map)
	if err != nil {
		log.Fatalf("issue loading toml file\n%+v", err)
	}
	return nil
}
func (t *TomlWrapper) Has(key string) bool {
	if key == "" {
		return false
	}
	return t.HasPath(strings.Split(key, "."))
	//if _, ok := t.Map[key]; ok {
	//	return true
	//}
	//return false
}
func (t *TomlWrapper) HasPath(keys []string) bool {
	return t.GetPath(keys) != nil
}

func (t *TomlWrapper) Get(key string) interface{} {
	if key == "" {
		return t
	}
	return t.GetPath(strings.Split(key, "."))
}
func (t *TomlWrapper) GetPath(keys []string) interface{} {
	if len(keys) == 0 {
		return t
	}
	subtree := t.Map
	for _, intermediateKey := range keys[:len(keys)-1] {
		value, exists := subtree[intermediateKey]
		if !exists {
			return nil
		}
		switch node := value.(type) {
		case map[string]interface{}:
			subtree = node
		case []interface{}:
			// go to most recent element
			if len(node) == 0 {
				return nil
			}
			idx := len(node) - 1
			subtree = node[idx].(map[string]interface{})
		default:
			return nil // cannot navigate through other node types
		}
	}
	// branch based on final node type
	switch node := subtree[keys[len(keys)-1]].(type) {
	default:
		return node
	}
}

//func LoadToml(path string) (*map[string]interface{}, error) {
//	data, err := ioutil.ReadFile(path)
//	if err != nil {
//		return nil, err
//	}
//	var props map[string]interface{}
//	err = toml.Unmarshal(data, &props)
//	if err != nil {
//		log.Fatalf("issue loading toml file\n%+v", err)
//	}
//	return &props, nil
//
//	//tree, err := toml.LoadFile(path)
//	//if err != nil {
//	//	log.Fatalf("issue loading toml file\n%+v", err)
//	//	return &toml.Tree{}, err
//	//}
//	//return tree, nil
//}
