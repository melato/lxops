package yaml

import (
	"testing"
)

func TestComment(t *testing.T) {
	data := []byte("#one\ntwo")
	if FirstLineComment(data) != "#one" {
		t.Fatalf("expected #one")
	}
	if !FirstLineIs(data, "#one") {
		t.Fatalf("expected #one")
	}
}
