package tto

import (
	"github.com/pelletier/go-toml"
)

var g *RegistryReader // global registry

type Registry struct {
	tree *toml.Tree
}

// custom registry
type RegistryReader struct {
	Data map[string]interface{}
	Keys []string
}

func LoadRegistry(path string) (*Registry, error) {
	tree, err := toml.LoadFile(path)
	if err == nil {
		re := &Registry{tree: tree}
		g = initRegistry(re)
		return re, err
	}
	return nil, err
}
func (r Registry) Contains(key string) bool {
	return r.tree.Has(key)
}
func (r Registry) GetString(key string) string {
	if r.Contains(key) {
		return r.Get(key).(string)
	}
	return ""
}

func (r Registry) Get(key string) interface{} {
	return r.tree.Get(key)
}
func (r Registry) GetBoolean(key string) bool {
	if r.Contains(key) {
		return r.Get(key).(bool)
	}
	return false
}

func GetRegistry() *RegistryReader {
	if g == nil {
		g = &RegistryReader{
			Keys: []string{},
			Data: make(map[string]interface{}, 0),
		}
	}
	return g
}

// init registry from config
func initRegistry(r *Registry) *RegistryReader {
	mp := map[string]interface{}{}
	keys := make([]string, 0)
	if n := r.tree.Get("global"); n != nil {
		if tree := n.(*toml.Tree); tree != nil {
			for _, k := range tree.Keys() {
				mp[k] = tree.Get(k)
				keys = append(keys, k)
			}
		}
	}
	return &RegistryReader{
		Data: mp,
		Keys: keys,
	}
}
