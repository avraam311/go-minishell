package integrationtests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const binName = "minishell"

func buildBinary(t *testing.T) string {
	t.Helper()
	binPath := filepath.Join(t.TempDir(), binName)
	cmd := exec.Command("go", "build", "-o", binPath, "../cmd/minishell")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, string(out))
	}
	return binPath
}

func runShell(t *testing.T, bin string, input string) string {
	t.Helper()

	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	for _, line := range strings.Split(input, "\n") {
		_, _ = stdin.Write([]byte(line + "\n"))
	}
	_ = stdin.Close()

	if err = cmd.Wait(); err != nil {
		_ = ""
	}

	return out.String()
}

func TestCd(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	out := runShell(t, bin, readFile(t, "./test_files/cd.txt"))
	res := strings.Fields(out)[1]
	expected := readFile(t, "./expected_files/cd.txt")
	if res != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", res, expected)
	}
}

func TestPwd(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	out := runShell(t, bin, readFile(t, "./test_files/pwd.txt"))
	res := strings.Fields(out)[1]
	expected := readFile(t, "./expected_files/pwd.txt")
	if res != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", res, expected)
	}
}

func TestEcho(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	out := runShell(t, bin, readFile(t, "./test_files/echo.txt"))
	res := strings.Join(strings.Fields(out)[1:4], " ")
	expected := readFile(t, "./expected_files/echo.txt")
	if res != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", res, expected)
	}
}

func TestKill(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	pid := cmd.Process.Pid
	bin := buildBinary(t)
	out := runShell(t, bin, readFile(t, "./test_files/kill.txt")+strconv.Itoa(pid))
	res := strings.Fields(out)[1]
	expected := readFile(t, "./expected_files/kill.txt")
	if res != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", res, expected)
	}
}

func TestPs(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	res := runShell(t, bin, readFile(t, "./test_files/ps.txt"))
	expected := readFile(t, "./expected_files/ps.txt")
	if !strings.Contains(res, expected) {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", "wrong result", expected)
	}
}

func TestExecPipes(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("touch", "test.txt")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	bin := buildBinary(t)
	res := runShell(t, bin, readFile(t, "./test_files/exec_pipes.txt"))
	expected := readFile(t, "./expected_files/exec_pipes.txt")
	if strings.Contains(res, expected) {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", res, "no test.txt")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(data)
}
