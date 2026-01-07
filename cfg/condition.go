package cfg

/*
EvalCondition is used to evaluate conditionals when loading Config files
It is a function that returns true if the corresponding condition is true.

This is called only with a non-empty argument.

Empty conditions are assumed to be true.

The condition is meant to be a property name and it evaluates to true,
if the property exists.

In the future this could be augmented to a more complex condition,
without changing the config file format or requiring a migration,
as long as the condition format is backward compatible.
*/
type EvalCondition func(condition string) bool

type Condition struct {
	// A string condition determines if the enclosing Config should be used or not.
	// The interpretation of the condition is up to EvalCondition
	// lxops currently evaluates conditions by checking if the condition string exists
	// as a global or command-line property.
	Condition string `yaml:"condition,omitempty"`
}

func (t *Condition) Eval(eval EvalCondition) bool {
	return t.Condition == "" || eval(t.Condition)
}
