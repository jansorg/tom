package util

import (
	"io"
	"os"
)

func CopyFile(source, target string, overwrite bool) error {
	from, err := os.OpenFile(source, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer from.Close()

	flags := os.O_WRONLY
	if overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_CREATE
	}

	to, err := os.OpenFile(target, flags, 0600)
	if err != nil {
		return err
	}
	defer to.Close()

	if _, err := io.Copy(to, from); err != nil {
		return err
	}

	if err := to.Sync(); err != nil {
		return err
	}

	srcInfo, err := os.Stat(target)
	if err != nil {
		return err
	}

	if err := os.Chmod(target, srcInfo.Mode()); err != nil {
		return err
	}

	return nil
}
