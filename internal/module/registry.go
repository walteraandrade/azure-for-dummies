package module

type Registry struct {
	modules []Module
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(m Module) {
	r.modules = append(r.modules, m)
}

func (r *Registry) All() []Module {
	return r.modules
}

func (r *Registry) Get(name string) (Module, bool) {
	for _, m := range r.modules {
		if m.Name() == name {
			return m, true
		}
	}
	return nil, false
}
