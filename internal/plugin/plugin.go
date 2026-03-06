package plugin

// HookPoint identifies when a plugin hook fires.
type HookPoint int

const (
	HookBeforeRender HookPoint = iota
	HookAfterRender
	HookBeforeIndex
	HookAfterIndex
	HookOnNavigate
)

// Hook is a function called at a specific point in processing.
type Hook struct {
	Name  string
	Point HookPoint
	Fn    func(ctx *HookContext) error
}

// HookContext provides data to hook functions.
type HookContext struct {
	FilePath string
	Content  string
	Metadata map[string]interface{}
}

// Registry manages registered plugin hooks.
type Registry struct {
	hooks map[HookPoint][]Hook
}

// NewRegistry creates an empty plugin registry.
func NewRegistry() *Registry {
	return &Registry{
		hooks: make(map[HookPoint][]Hook),
	}
}

// Register adds a hook at the specified point.
func (r *Registry) Register(hook Hook) {
	r.hooks[hook.Point] = append(r.hooks[hook.Point], hook)
}

// Run executes all hooks registered at the given point.
func (r *Registry) Run(point HookPoint, ctx *HookContext) error {
	for _, hook := range r.hooks[point] {
		if err := hook.Fn(ctx); err != nil {
			return err
		}
	}
	return nil
}
