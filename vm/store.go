package vm

import (
	"github.com/Zac-Garby/pluto/object"
)

type item struct {
	name  string
	value object.Object

	// Whether the item is defined locally
	local bool
}

// Store is an evaluation scope: it stores
// defined names, and their corresponding data
type Store struct {
	Names    []string
	Patterns []string
	Data     []*item

	*FunctionStore
}

// NewStore creates an empty store
func NewStore() *Store {
	return &Store{
		Names:         make([]string, 0),
		Patterns:      make([]string, 0),
		Data:          make([]*item, 0),
		FunctionStore: &FunctionStore{Functions: make([]object.Function, 0)},
	}
}

// Define name in the store, and returns its name index
func (s *Store) Define(name string, val object.Object, local bool) rune {
	s.Names = append(s.Names, name)

	for i, item := range s.Data {
		if item.name == name {
			item.value = val
			item.local = local

			return rune(i)
		}
	}

	s.Data = append(s.Data, &item{
		local: local,
		name:  name,
		value: val,
	})

	return rune(len(s.Names) - 1)
}

// GetName searches the store for data named 'name'
func (s *Store) GetName(name string) object.Object {
	for _, item := range s.Data {
		if item.name == name {
			return item.value
		}
	}

	return nil
}

// GetID searches the store for the data whose id is 'id'
func (s *Store) GetID(id rune) (string, object.Object) {
	name := s.Names[id]
	return name, s.GetName(name)
}

// ImportModule imports a module: If _module is
// defined in the module store, its contents are
// copied into a new map inside this store, assigned
// to other["title"].
func (s *Store) ImportModule(other *Store, name string) {
	var (
		module     *object.Map
		moduleName string
	)

	if mod := other.GetName("_module"); mod != nil {
		// _module exists in the imported package

		var (
			title object.Object
			ok    bool
		)

		module, ok = mod.(*object.Map)
		if !ok {
			goto invalid
		}

		title = module.Get(&object.String{Value: "title"})
		if title == nil {
			goto invalid
		}

		moduleName = title.String()
	}

	goto after_fallback
invalid:
	// _module doesn't exist
	// importing from file name

	module = &object.Map{
		Keys:   make(map[string]object.Object),
		Values: make(map[string]object.Object),
	}

after_fallback:

	functions := &object.Array{
		Value: make([]object.Object, len(other.Functions)),
	}

	for i, fn := range other.Functions {
		functions.Value[i] = &object.Function{
			Body:      fn.Body,
			Constants: fn.Constants,
			Names:     fn.Names,
			Pattern:   fn.Pattern,
			Patterns:  fn.Patterns,
		}
	}

	var (
		mkey       = &object.String{Value: "_methods"}
		oldMethods = module.Get(mkey)
		appended   = addMethods(oldMethods, functions.Elements())
	)

	module.Set(mkey, &object.Array{Value: appended})

	for _, item := range other.Data {
		if item.name == "_module" || !item.local {
			continue
		}

		module.Set(&object.String{Value: item.name}, item.value)
	}

	s.Define(moduleName, module, false)
}

func addMethods(old object.Object, methods []object.Object) []object.Object {
	if old == nil || old.Equals(object.NullObj) {
		return methods
	}

	oldArr, ok := old.(*object.Array)
	if !ok {
		return methods
	}

	return append(methods, oldArr.Elements()...)
}
