package cfg

type HasProperty func(string) bool

type Condition struct {
	// The condition passes if the property exists or is empty)
	Property string `yaml:"if,omitempty"`
}

func (t *Condition) Eval(hasProperty HasProperty) bool {
	return t.Property == "" || hasProperty(t.Property)
}
