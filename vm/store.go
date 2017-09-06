package vm

import "github.com/Zac-Garby/pluto/object"

// Store is an evaluation scope: it stores
// defined names, and their corresponding data
type Store struct {
	Names     []string
	Patterns  []string
	Data      map[string]object.Object
	Functions FunctionStore
}

// NewStore creates an empty store
func NewStore() Store {
	return Store{
		Names:     make([]string, 0),
		Patterns:  make([]string, 0),
		Data:      make(map[string]object.Object),
		Functions: FunctionStore{Functions: make([]object.Function, 8)},
	}
}

// Define name in the store, and returns its name index
func (s *Store) Define(name string, val object.Object) rune {
	s.Names = append(s.Names, name)
	s.Data[name] = val

	return rune(len(s.Names) - 1)
}

// SearchName searches the store for data named 'name'
func (s *Store) SearchName(name string) object.Object {
	return s.Data[name]
}

// SearchID searches the store for the data whose id is 'id'
func (s *Store) SearchID(id rune) (string, object.Object) {
	name := s.Names[id]
	return name, s.Data[name]
}
