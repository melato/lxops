package cli

import (
	"fmt"
	"os"

	"melato.org/lxops/srv"
	"melato.org/lxops/util"
)

const Running = "Running"

// Assign persistent number to containers.
// Used to assign a unique ssh port to each container
type AssignNumbers struct {
	Client  srv.Client
	File    string `name:"f" usage:"numbers CSV file (container,number)"`
	First   int    `name:"first" usage:"first number"`
	Last    int    `name:"last" usage:"last number (optional)"`
	Project string `name:"project" usage:"add containers from LXD project"`
	Running bool   `name:"r" usage:"use only running containers from specified project"`
	Clean   bool   `name:"clean" usage:"remove numbers for containers that are not selected"`
}

func (t *AssignNumbers) Configured() error {
	if t.File == "" {
		return fmt.Errorf("missing file")
	}
	if t.Running && t.Project == "" {
		return fmt.Errorf("-r can be used only with -project")
	}
	return nil
}

func (t *AssignNumbers) selectContainers(names []string, f func(name string) error) error {
	selectedNames := make(map[string]bool)
	// add names from LXD project
	if t.Project != "" {
		server, err := t.Client.ProjectInstanceServer(t.Project)
		if err != nil {
			return err
		}

		instances, err := server.GetInstanceNames(t.Running)
		if err != nil {
			return err
		}
		for _, instance := range instances {
			selectedNames[instance] = true
			err = f(instance)
			if err != nil {
				return err
			}
		}
	}
	// add names from command line
	for _, name := range names {
		if !selectedNames[name] {
			selectedNames[name] = true
			err := f(name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func selectNumbers(numbers []*util.NamedNumber, names []string) []*util.NamedNumber {
	namesMap := make(map[string]bool)
	for _, name := range names {
		namesMap[name] = true
	}
	var result []*util.NamedNumber
	for _, num := range numbers {
		if namesMap[num.Name] {
			result = append(result, num)
		}
	}
	return result
}

func MinNumber(numbers []*util.NamedNumber) int {
	if len(numbers) == 0 {
		return 0
	}
	min := numbers[0].Value
	for _, num := range numbers[1:] {
		if num.Value < min {
			min = num.Value
		}
	}
	return min
}

func (t *AssignNumbers) AddNumbers(numbers []*util.NamedNumber, names []string) ([]*util.NamedNumber, error) {
	var containers []string
	err := t.selectContainers(names, func(name string) error {
		containers = append(containers, name)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if t.Clean {
		numbers = selectNumbers(numbers, containers)
	}
	usedNumbers := make(map[int]bool)
	numberedContainers := make(map[string]bool)
	nextNumber := t.First
	if nextNumber == 0 {
		if len(numbers) == 0 {
			return nil, fmt.Errorf("no numbers found.  please specify -first")
		}
		nextNumber = MinNumber(numbers)
	}

	for _, num := range numbers {
		usedNumbers[num.Value] = true
		numberedContainers[num.Name] = true
	}
	for _, name := range containers {
		if !numberedContainers[name] {
			for ; ; nextNumber++ {
				if t.Last != 0 && nextNumber > t.Last {
					return nil, fmt.Errorf("no numbers available between %d, %d", t.First, t.Last)
				}
				if !usedNumbers[nextNumber] {
					numbers = append(numbers, &util.NamedNumber{name, nextNumber})
					nextNumber++
					break
				}
			}
		}
	}
	return numbers, nil
}

func (t *AssignNumbers) Run(containers []string) error {
	var numbers []*util.NamedNumber
	_, err := os.Stat(t.File)
	if err == nil {
		numbers, err = util.ReadNumbers(t.File)
	} else if os.IsNotExist(err) {
		err = nil
		// the file does not exist, start with empty list
	}

	if err != nil {
		return err
	}
	numbers, err = t.AddNumbers(numbers, containers)
	if err != nil {
		return err
	}
	return util.WriteNumbers(numbers, t.File)
}
