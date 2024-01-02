package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

type ShiftIds struct {
	UidShift int  `name:"u" usage:"uid shift"`
	GidShift int  `name:"g" usage:"gid shift"`
	Verbose  bool `name:"v"`
}

type shiftDir struct {
	ShiftIds
	Dir string
}

func (t *shiftDir) shift(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if (fs.ModeSymlink & d.Type()) != 0 {
		// ignore symbolic links
		return nil
	}
	info, err := d.Info()
	if err != nil {
		return err
	}
	s := info.Sys().(*syscall.Stat_t)
	uid := int(s.Uid)
	gid := int(s.Gid)
	if uid < t.UidShift {
		uid += t.UidShift
	} else {
		uid = -1
	}
	if gid < t.GidShift {
		gid += t.GidShift
	} else {
		gid = -1
	}
	if uid != -1 || gid != -1 {
		fullpath := filepath.Join(t.Dir, path)
		if t.Verbose {
			fmt.Printf("chown %d:%d %s\n", uid, gid, fullpath)
		}
		return os.Chown(fullpath, uid, gid)
	}
	return nil
}

func (t *ShiftIds) Run(dirs ...string) error {
	var shift shiftDir
	shift.ShiftIds = *t
	for _, dir := range dirs {
		shift.Dir = dir
		err := fs.WalkDir(os.DirFS(dir), ".", shift.shift)
		if err != nil {
			return err
		}
	}
	return nil
}
