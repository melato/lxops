package lxops

import (
	"testing"

	"melato.org/lxops/util"
)

func TestConfigProfiles(t *testing.T) {
	var config Config
	config.ProfilesConfig = []string{"config"}
	config.ProfilesRun = []string{"run"}
	profiles := []string{"a", "run", "b"}
	profiles = config.GetProfilesConfig(profiles)
	if !util.StringSlice(profiles).Equals([]string{"config", "a", "b"}) {
		t.Fatalf("%v", profiles)
	}
}
