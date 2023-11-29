package doc

type TypeDescriptor struct {
	Description string            `yaml:"description"`
	Fields      map[string]string `yaml:"fields"`
}
