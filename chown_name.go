package pcommon

import (
	"os"
)

func ChownName(name string, username, groupname string) error {
	uid, gid, err := convertIds(username, groupname)
	if err != nil {
		return err
	}

	return os.Chown(name, int(uid), int(gid))
}
