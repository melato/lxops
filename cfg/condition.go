package cfg

type HasProperty func(string) bool

type Condition struct {
	// Properties is a list of property names
	// The condition passes if all properties with these names exist
	Properties []string `yaml:"if,omitempty"`
}

func (t *Condition) Eval(hasProperty HasProperty) bool {
	for _, name := range t.Properties {
		if !hasProperty(name) {
			return false
		}
	}
	return true
}
