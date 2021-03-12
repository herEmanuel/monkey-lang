package object

func NewExtendedEnvironment(outer *Environment) *Environment {
	newEnv := NewEnvironment()
	newEnv.outer = outer
	return newEnv
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) Object {
	if value, ok := e.store[name]; ok {
		return value
	} else if e.outer != nil {
		value := e.outer.Get(name)
		return value
	}
	return nil
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
