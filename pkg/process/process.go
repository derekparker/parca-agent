package process

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// ErrInvalidPid is returned when attempting to detect a process with an invalid pid.
type ErrInvalidPid struct {
	pid int
}

func (e ErrInvalidPid) Error() string {
	return fmt.Sprintf("invalid pid: %d", e.pid)
}

// Process represents an OS process.
// It can be used to perform certain actions on the process that
// are required to begin profiling it.
type Process interface {
	Pid() int
	PrepareToProfile() error
}

// NonSpecialProcess represents a process that does not require special
// handling to begin profiling it.
type NonSpecialProcess struct {
	pid  int
	path string
}

// Pid returns the pid of the process.
func (ns *NonSpecialProcess) Pid() int {
	return ns.pid
}

// PrepareToProfile takes a process and attempts to prepare it for
// profiling. It is a noop for NonSpecialProcess.
func (ns *NonSpecialProcess) PrepareToProfile() error {
	return nil
}

// JavaProcess represents a Java process.
type JavaProcess struct {
	NonSpecialProcess
}

// PrepareToProfile takes a process and attempts to prepare it for
// profiling. For Java processes it will attempt to auto attach the
// java agent to the process to export the JIT symbol map.
func (ns *JavaProcess) PrepareToProfile() error {
	// Command to attach the agent to the process:
	// (On linux we must be the same user/group as the process)
	// sudo -u \#$TARGET_UID -g \#$TARGET_GID $JAVA_HOME/bin/java -Xms32m -Xmx128m -cp $AGENT_JAR:$JAVA_HOME/lib/tools.jar net.virtualvoid.perf.AttachOnce $pid $opts
	var (
		javaHome = os.Getenv("JAVA_HOME")
		java     = ns.path
		agent    = "$AGENT_JAR"
		opts     = "unfoldall"
	)

	cmd := exec.Command(java,
		"-Xms32m",  // Minimum heap size.
		"-Xmx128m", // Maximum heap size.
		"-cp", fmt.Sprintf("%s:%s/lib/tools.jar", agent, javaHome),
		"net.virtualvoid.perf.AttachOnce", fmt.Sprintf("%d", ns.pid), opts)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(os.Getuid()),
			Gid: uint32(os.Getgid()),
		},
	}

	return cmd.Run()
}

// Detect takes a pid and attempts to determine the type of process
// that belongs to that pid.
func Detect(pid int) (Process, error) {
	path, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		return nil, err
	}

	var (
		p Process

		nsp = NonSpecialProcess{pid: pid, path: path}
	)

	switch {
	case strings.Contains(path, "java"):
		p = &JavaProcess{nsp}
	default:
		p = &nsp
	}

	return p, nil
}
