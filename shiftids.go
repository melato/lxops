package lxops

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"melato.org/script"
)

type shiftIds struct {
	Owner      string
	Uid        string
	Gid        string
	Executable string
}

func (t *shiftIds) ShiftDir(s *script.Script, dir string) {
	if t.Uid != "" || t.Gid != "" {
		s.Run("sudo", t.Executable, "shiftids", "-u", t.Uid, "-g", t.Gid, "-v", dir)
	}
}

func newShiftIds(owner string) (*shiftIds, error) {
	if owner == "" {
		return &shiftIds{}, nil
	}
	uid, gid, ok := parseOwner(owner)
	if !ok {
		return nil, fmt.Errorf("owner should have the form uid:gid (%s)", owner)
	}
	var perm shiftIds
	perm.Uid = strconv.Itoa(uid)
	perm.Gid = strconv.Itoa(gid)
	var err error
	perm.Executable, err = os.Executable()
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func parseOwner(owner string) (int, int, bool) {
	parts := strings.Split(owner, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	ids := make([]int, len(parts))
	for i, s := range parts {
		var err error
		ids[i], err = strconv.Atoi(s)
		if err != nil {
			return 0, 0, false
		}
	}
	return ids[0], ids[1], true
}
