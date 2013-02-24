// TOML markup language parser.
package toml

import (
	"strings"
)

// Definition of a TomlTree.
// This is the result of the parsing of a TOML file.
type TomlTree map[string]interface{}

// Keys returns the keys of the toplevel tree.
// Warning: this is a costly operation.
func (t *TomlTree) Keys() []string {
	keys := make([]string, 0)
	mp := (map[string]interface{})(*t)
	for k, _ := range mp {
		keys = append(keys, k)
	}
	return keys
}

func (t *TomlTree) Get(key string) interface{} {
	subtree := t
	keys := strings.Split(key, ".")
	for _, intermediate_key := range keys[:len(keys)-1] {
		_, exists := (*subtree)[intermediate_key]
		if !exists {
			return nil
		}
		subtree = (*subtree)[intermediate_key].(*TomlTree)
	}
	return (*subtree)[keys[len(keys)-1]]
}

// Set an element in the tree.
// Creates all necessary intermediates trees, if needed.
func (t *TomlTree) Set(key string, value interface{}) {
	subtree := t
	keys := strings.Split(key, ".")
	for _, intermediate_key := range keys[:len(keys)-1] {
		_, exists := (*subtree)[intermediate_key]
		if !exists {
			var new_tree TomlTree = make(TomlTree)
			(*subtree)[intermediate_key] = &new_tree
		}
		subtree = (*subtree)[intermediate_key].(*TomlTree)
	}
	(*subtree)[keys[len(keys)-1]] = value
}

// createSubTree takes a tree and a key andcreate the necessary intermediate
// subtrees to create a subtree at that point. In-place.
//
// e.g. passing a.b.c will create (assuming tree is empty) tree[a], tree[a][b]
// and tree[a][b][c]
func (t *TomlTree) createSubTree(key string) {
	subtree := t
	for _, intermediate_key := range strings.Split(key, ".") {
		_, exists := (*subtree)[intermediate_key]
		if !exists {
			var new_tree TomlTree = make(TomlTree)
			(*subtree)[intermediate_key] = &new_tree
		}
		subtree = ((*subtree)[intermediate_key]).(*TomlTree)
	}
}

func Load(content string) *TomlTree {
	_, flow := lex(content)
	return parse(flow)
}
