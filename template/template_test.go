package template

import (
	"testing"
)

func testTemplateKey(t *testing.T, template string, keyTemplate KeyExpression, expected string) {
	s, err := keyTemplate.Apply(template, "a", "1", "b", "2")
	if err != nil {
		t.Error(err)
		return
	}
	if s != expected {
		t.Errorf("%s != %s", s, expected)
	}
}

func TestTemplate(t *testing.T) {
	testTemplateKey(t, "${a}/${b}", Ant, "1/2")
	testTemplateKey(t, "(a)/(b)", Paren, "1/2")
	testTemplateKey(t, "+(a)/(b)-", Paren, "+1/2-")
	if _, err := Paren.Apply("(a)/(b)", "a", "1"); err == nil {
		t.Error("no error for missing key")
	}
}
