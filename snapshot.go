package virtualbox

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type Snapshot struct {
	ID   string
	Name string
}

var errNoSnapshots = "This machine does not have any snapshots"

// ListSnapshots returns the list the machine's snapshots.
// If id is blank, use current machine.
func (m *Machine) ListSnapshots(id string) ([]Snapshot, error) {
	var (
		out string
		err error
	)
	if id == "" {
		out, _, err = Manage().run("snapshot", m.Name, "list")
	} else {
		out, _, err = Manage().run("snapshot", id, "list")
	}
	if err != nil && !strings.Contains(out, errNoSnapshots) {
		return nil, err
	}
	return parseSnapshots(out)
}

// DeleteSnapshot deletes a snapshot given snapshot name
func (m *Machine) DeleteSnapshot(name string) error {
	snapshots, err := m.ListSnapshots("")
	if err != nil {
		return err
	}

	found := false
	for _, s := range snapshots {
		if s.Name == name {
			found = true
			break
		}
	}
	if !found {
		return errors.New("snapshot not found")
	}

	_, _, err = Manage().run("snapshot", m.Name, "delete", name)
	return err
}

// CreateSnapshot creates a snapshot given snapshot name
func (m *Machine) CreateSnapshot(name string) error {
	snapshots, err := m.ListSnapshots("")
	if err != nil {
		return err
	}

	for _, s := range snapshots {
		if s.Name == name {
			return errors.New("snapshot already exists")
		}
	}

	out, _, err := Manage().run("snapshot", m.Name, "take", name)

	s := bufio.NewScanner(strings.NewReader(out))
	for s.Scan() {
		line := s.Text()
		fmt.Println(line)
		if line == "" {
			continue
		}
		match := reSnapshotFailed.FindStringSubmatch(line)
		if len(match) == 2 {
			return errors.New(match[1])
		}
	}

	return err
}

// RestoreSnapshot restores a snapshot given snapshot name
func (m *Machine) RestoreSnapshot(name string) error {
	snapshots, err := m.ListSnapshots("")
	if err != nil {
		return err
	}

	found := false
	for _, s := range snapshots {
		if s.Name == name {
			found = true
			break
		}
	}
	if !found {
		return errors.New("snapshot not found")
	}

	_, _, err = Manage().run("snapshot", m.Name, "restore", name)
	return err
}

func parseSnapshots(out string) ([]Snapshot, error) {
	var (
		snapshots []Snapshot
	)

	s := bufio.NewScanner(strings.NewReader(out))
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		match := reSnapshotList.FindStringSubmatch(line)
		if len(match) == 3 {
			snapshots = append(snapshots, Snapshot{
				Name: match[1],
				ID:   match[2],
			})
		}
	}

	return snapshots, nil
}
