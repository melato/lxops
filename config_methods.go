package lxops

import (
	"fmt"
	"os"
	"path/filepath"

	"melato.org/lxops/util"
)

func (t *OS) String() string {
	return t.Name
}

func (t *OS) Equals(x *OS) bool {
	return t.Name == x.Name
}

/** Check that the requirements are met */
func (t *Config) VerifyFileExists(file HostPath) bool {
	if file != "" && !util.FileExists(string(file)) {
		fmt.Fprintf(os.Stderr, "file does not exist: %s\n", string(file))
		return false
	}
	return true
}

func (config *Config) verifyDevices() bool {
	valid := true
	devicePaths := make(map[string]bool)
	for _, d := range config.Devices {
		if config.Filesystems[d.Filesystem] == nil {
			valid = false
			fmt.Fprintf(os.Stderr, "unknown filesystem id: %s\n", d.Filesystem)
		}
		if devicePaths[d.Path] {
			valid = false
			fmt.Fprintf(os.Stderr, "duplicate device path: %s\n", d.Path)
		}
		devicePaths[d.Path] = true
	}
	return valid
}

func (config *Config) Verify() bool {
	valid := true
	for _, file := range config.CloudConfigFiles {
		if !config.VerifyFileExists(file) {
			valid = false
		}
	}
	if !config.VerifyFileExists(config.SourceConfig) {
		valid = false
	}
	if !config.verifyDevices() {
		valid = false
	}

	duplicates := config.getDuplicates(config.Profiles)
	if len(duplicates) > 0 {
		valid = false
		fmt.Fprintf(os.Stderr, "duplicate profiles: %v\n", duplicates)
	}

	return valid
}

func (t *Config) Print() error {
	return PrintConfigYaml(t)
}

/** Read config without includes */
func ReadRawConfig(file string) (*Config, error) {
	return ReadConfigYaml(file)
}

func (path HostPath) Resolve(dir string) HostPath {
	if path == "" {
		return ""
	}
	if filepath.IsAbs(string(path)) {
		return path
	}
	return HostPath(filepath.Join(dir, string(path)))
}

func (t *Config) ResolvePaths(dir string) {
	for i, f := range t.Include {
		t.Include[i] = f.Resolve(dir)
	}
	t.SourceConfig = t.SourceConfig.Resolve(dir)
	for i, path := range t.CloudConfigFiles {
		t.CloudConfigFiles[i] = path.Resolve(dir)

	}
}

// Return the filesystem for the given id, or nil if it doesn't exist.
func (t *Config) Filesystem(id string) *Filesystem {
	return t.Filesystems[id]
}

func (t *Config) getDuplicates(lists ...[]string) []string {
	var duplicates []string
	set := make(util.Set[string])
	for _, list := range lists {
		for _, s := range list {
			if set.Contains(s) {
				duplicates = append(duplicates, s)
			}
			set.Put(s)
		}
	}
	return duplicates
}

func (t *Config) HasProfilesConfig() bool {
	return len(t.ProfilesConfig)+len(t.ProfilesRun) > 0
}

func (t *Config) GetProfilesConfig(profiles []string) []string {
	if !t.HasProfilesConfig() {
		return profiles
	}
	profiles = util.StringSlice(profiles).Diff(t.ProfilesRun)
	profiles = util.StringSlice(t.ProfilesConfig).Union(profiles)
	return profiles
}

func (t *Config) AbsFilename() string {
	filename := t.File
	if !filepath.IsAbs(filename) {
		dir, err := os.Getwd()
		if err == nil {
			filename = filepath.Join(dir, filename)
		}
	}
	return filename
}
