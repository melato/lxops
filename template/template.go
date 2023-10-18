package template

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type KeyExpression string

const (
	// Use keys of the form ${key}.  Compatible with sh and ant scripts.
	Ant = KeyExpression(`\$\{([^{}]+)\}`)

	// Use keys of the form (key)
	Paren = KeyExpression(`\(([^()]+)\)`)
)

// A simple template substitution where keys are substituted by provided values.
// Kes are recognized via the Key regular expression.
type Template struct {
	template string
	// The keys that we found in the template, in order
	keys []string
	// The parts of the template between the keys.  for n keys, we have n + 1 pieces
	pieces []string
}

// ErrorHandler allows a function to report errors, without returning them directly
// It is used for convenience in order to minimize checking error return values from function calls
type ErrorHandler interface {
	// Add an error.  If err is nil, Handle should do nothing.
	Add(err error)
	// Return true if processing should continue.  Should return true if there were no errors.
	// May also return true if there were prior errors, but processing should continue in order to check for more errors.
	Continue() bool
}

// Example:  Template{"(key1)/(key2)"}.Apply("key1", "value1", "key2", "value2") -> "value1/value2"
func (keyTemplate KeyExpression) NewTemplate(template string) (*Template, error) {
	tpl := &Template{template: template}
	if err := tpl.compile(keyTemplate); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (t *Template) compile(keyTemplate KeyExpression) error {
	re, err := regexp.Compile(string(keyTemplate))
	if err != nil {
		return err
	}
	matches := re.FindAllStringSubmatchIndex(t.template, -1)
	t.keys = make([]string, len(matches))
	t.pieces = make([]string, len(matches)+1)
	start := 0
	for i, m := range matches {
		t.pieces[i] = t.template[start:m[0]]
		if len(m) != 4 {
			return errors.New(fmt.Sprintf("expected 4 indexes from key template.  got %d", len(m)))
		}
		t.keys[i] = t.template[m[2]:m[3]]
		start = m[1]
	}
	t.pieces[len(t.pieces)-1] = t.template[start:]
	return nil
}

func (t *Template) find(args []string, key string) (string, error) {
	n := len(args)
	var found bool
	var value string
	for i := 0; i < n; i += 2 {
		if key == args[i] {
			if found {
				return "", errors.New("duplicate key: " + key)
			}
			found = true
			value = args[i+1]
		}
	}
	if !found {
		keys := make([]string, n/2)
		for i := 0; i < n; i += 2 {
			keys[i/2] = args[i]
		}
		return "", errors.New(fmt.Sprintf("missing key: %s. provided keys: %v", key, keys))
	}
	return value, nil
}

// Applyf substitutes the template keys with values provided by a function.
func (t *Template) Applyf(f func(key string) (string, error)) (string, error) {
	result := make([]string, 2*len(t.keys)+1)
	for i, s := range t.pieces {
		result[i*2] = s
	}
	for i, key := range t.keys {
		value, err := f(key)
		if err != nil {
			return "", err
		}
		result[i*2+1] = value
	}
	return strings.Join(result, ""), nil
}

// Apply substitutes the given key/value pairs to the template.
// The keys and values are alternate arguments, so for even i, arg[i] is a key, and arg[i+1] is the corresponding value.
// It is an error if a template key is not present in the args.
func (t *Template) Apply(arg ...string) (string, error) {
	if len(arg)%2 != 0 {
		return "", errors.New("must have even args")
	}
	return t.Applyf(func(key string) (string, error) { return t.find(arg, key) })
}

// Apply creates and applies a template in a single call, using this key expression
func (t KeyExpression) Apply(template string, arg ...string) (string, error) {
	tpl, err := t.NewTemplate(template)
	if err != nil {
		return "", err
	}
	return tpl.Apply(arg...)
}

// ApplyE is like Apply, but it stores its error in the given *error argument.  If there is an error there already, do nothing.
// Use to easily make a series of template substitutions and worry about any errors at the end.
func (t KeyExpression) ApplyE(errors ErrorHandler, template string, arg ...string) string {
	if !errors.Continue() {
		return ""
	}
	result, err := t.Apply(template, arg...)
	if err != nil {
		errors.Add(err)
		return ""
	}
	return result
}

func (t KeyExpression) ReplaceE(errors ErrorHandler, text *string, arg ...string) {
	if !errors.Continue() {
		return
	}
	result, err := t.Apply(*text, arg...)
	if err != nil {
		errors.Add(err)
	}
	*text = result
}
