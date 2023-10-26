package migrate

import (
	"fmt"
	"testing"
)

func TestMigrateLxdops(t *testing.T) {
	data, err := testFS.ReadFile("test/lxdops.yaml")
	var expected []byte
	if err == nil {
		expected, err = testFS.ReadFile("test/lxdops_migrated.yaml")
	}
	if err == nil {
		data, err = MigrateLxdops(data)
	}
	if err == nil {
		if !EqualYaml(data, expected) {
			fmt.Println(string(data))
			t.Fatalf("not expected")
			return
		}
	}
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
}
