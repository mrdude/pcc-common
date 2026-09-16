package pcommon

import (
	"errors"
	"os/user"
	"strconv"
	"syscall"
)

func DropPrivs(username, groupname string) error {
	uid, gid, err := convertIds(username, groupname)
	if err != nil {
		return err
	}

	if err := dropPrivs(uid, gid); err != nil {
		return err
	}

	// ensure that we can't regain our privs
	if err := syscall.Setgid(0); err == nil {
		return errors.New("failed to drop privs")
	}

	if err := syscall.Setuid(0); err == nil {
		return errors.New("failed to drop privs")
	}

	return nil
}

// converts a username and groupname into a uid and gid
func convertIds(username, groupname string) (uid int64, gid int64, err error) {
	var linuxUser *user.User
	linuxUser, err = user.Lookup(username)
	if err != nil {
		linuxUser, err = user.LookupId(username)
	}
	if err != nil {
		return
	}

	var linuxGroup *user.Group
	linuxGroup, err = user.LookupGroup(groupname)
	if err != nil {
		linuxGroup, err = user.LookupGroupId(groupname)
	}
	if err != nil {
		return
	}

	gid, err = strconv.ParseInt(linuxGroup.Gid, 10, 32)
	if err != nil {
		return
	}

	uid, err = strconv.ParseInt(linuxUser.Uid, 10, 32)
	if err != nil {
		return
	}

	return
}

func dropPrivs(uid, gid int64) error {
	if err := syscall.Setgid(int(gid)); err != nil {
		return err
	}

	if err := syscall.Setuid(int(uid)); err != nil {
		return err
	}

	return nil
}

func MustDropPrivs(username, groupname string) {
	if err := DropPrivs(username, groupname); err != nil {
		panic(err)
	}
}
