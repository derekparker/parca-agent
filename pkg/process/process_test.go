package process

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestDetectInvalidPid(t *testing.T) {
	_, err := Detect(0)
	if err == nil {
		t.Error("Expected error for invalid pid")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected ErrNotExist error, got %T", err)
	}
}

// TestDetectNonSpecialBinary tests that Detect() returns essentially a
// noop Process for non-special binaries, e.g. those which don't require
// anything special to happen to be able to profile them.
func TestDetectNonSpecialBinary(t *testing.T) {
	p, err := Detect(os.Getpid())
	if err != nil {
		t.Error("Expected error for non-special binary")
	}
	if _, ok := p.(*NonSpecialProcess); !ok {
		t.Errorf("Expected NonSpecialProcess, got %T", p)
	}
}

func TestDetectJava(t *testing.T) {
	// We must compile and run a Java program so skip the test
	// if we don't have a JDK.
	if _, err := exec.LookPath("java"); err != nil {
		t.Skip("java not found")
	}

	if output, err := exec.Command("javac", "testdata/DarkRoast.java").CombinedOutput(); err != nil {
		t.Fatalf("javac failed: %v\n%s", err, output)
	}

	cmd := exec.Command("java", "testdata/DarkRoast")

	if err := cmd.Start(); err != nil {
		t.Fatalf("running java process failed: %v", err)
	}

	defer func(t *testing.T, cmd *exec.Cmd) {
		if err := cmd.Process.Kill(); err != nil {
			t.Errorf("killing java process failed: %v", err)
		}
		cmd.Wait()

		if err := os.Remove("testdata/DarkRoast.class"); err != nil {
			t.Errorf("removing DarkRoast.class failed: %v", err)
		}
	}(t, cmd)

	p, err := Detect(cmd.Process.Pid)
	if err != nil {
		t.Error("Expected error for Java binary")
	}
	if _, ok := p.(*JavaProcess); !ok {
		t.Errorf("Expected JavaProcess, got %T", p)
	}
}
