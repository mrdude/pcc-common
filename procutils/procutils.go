package procutils

import (
	"os/exec"
	"syscall"
)

// KillChildrenOnProcessDeath configures an *exec.Cmd created with exec.CommandContext.
// When this process' context is done, this process's children will be killed along with the process itself.
//
// This allows us to cancel running ansible-playbook runs without leaving orphaned mitogen connections around.
func KillChildrenOnProcessDeath(cmd *exec.Cmd) {
	// https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// when creating this process, set its ProcessGroupId to -PID
		// this allows us to kill the process and all of its subprocesses with syscall.Kill(-ChildPID, syscall.SIGKILL)
		Setpgid: true,
	}

	cmd.Cancel = func() error {
		// send SIGKILL to this process, as well as all of it's children
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}
