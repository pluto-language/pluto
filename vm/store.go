package vm

import "github.com/Zac-Garby/pluto/object"

// Store is an evaluation scope: it stores
// defined names, and their corresponding data
type Store struct {
	Names    []string
	Patterns []string
	Data     map[string]object.Object

	*FunctionStore
}

// NewStore creates an empty store
func NewStore() *Store {
	return &Store{
		Names:         make([]string, 0),
		Patterns:      make([]string, 0),
		Data:          make(map[string]object.Object),
		FunctionStore: &FunctionStore{Functions: make([]object.Function, 0)},
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

// Extend extends one store with the data from the other.
// If a variable named _module is found, it's renamed
// to x where x = _module["title"]. If _module["title"]
// is not defined, _module is copied over normally.
func (s *Store) Extend(other *Store) {
	for k, v := range other.Data {
		if k == "_module" {
			var (
				m                *object.Map
				title            object.Object
				newVar           string
				fnArr            object.Object
				generalFunctions []object.Object
			)

			m, ok := v.(*object.Map)
			if !ok {
				goto invalid
			}

			title = m.Get(&object.String{Value: "title"})
			if title == nil {
				goto invalid
			}

			generalFunctions = make([]object.Object, len(other.Functions))

			for i, fn := range other.Functions {
				generalFunctions[i] = &fn
			}

			fnArr = &object.Array{
				Value: generalFunctions,
			}

			m.Set(&object.String{Value: "_methods"}, fnArr)

			newVar = title.String()
			s.Data[newVar] = m

		invalid:
		}

		s.Data[k] = v
	}

	s.Functions = append(other.Functions, s.Functions...)
}
