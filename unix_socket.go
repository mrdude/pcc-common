package pcommon

import (
	"net"
	"os"
)

func BindUnixServerSocket(socketPath string, mode os.FileMode) (net.Listener, error) {
	// delete the socket, if it exists
	if _, err := os.Stat(socketPath); err == nil {
		err = os.Remove(socketPath)
		if err != nil {
			return nil, err
		}
	}

	// start listening
	sock, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}

	err = os.Chmod(socketPath, mode)
	if err != nil {
		sock.Close()
		return nil, err
	}

	return sock, nil
}
