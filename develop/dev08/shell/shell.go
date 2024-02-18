package shell

import (
	"bufio"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===

# Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Shell struct {
	wd string
}

func NewShell() *Shell {
	wd, _ := os.Getwd()
	return &Shell{
		wd: wd,
	}
}

func (s *Shell) readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSuffix(input, "\r\n")
	return input, nil
}

func (s *Shell) RunShell() {
	for {
		fmt.Printf("%s > ", s.wd)
		input, err := s.readInput()
		commands := strings.Split(input, "|")
		if err != nil {
			fmt.Print("Input err: %s\n", err.Error())
			continue
		}
		if len(commands) > 1 {
			fmt.Println("commands case")
		} else {
			if err = s.runCommand(input); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (s *Shell) runCommand(command string) error {
	c := strings.Split(command, " ")
	if c[0] == "" {
		return nil
	}
	switch c[0] {
	case "cd":
		c := strings.Split(command, " ")
		if len(c) < 2 {
			return &NoArgumentsPassed{}
		}
		path := c[1]
		if err := s.dirChange(path); err != nil {
			return err
		}
	case "ls":
		if err := s.ls(); err != nil {
			return err
		}
	case "pwd":
		fmt.Println("- " + s.wd)
	case "echo":
		for i, text := range c {
			if i != 0 {
				fmt.Print(text)
			}
		}
		fmt.Println("\n")
	case "kill":
		c := strings.Split(command, " ")
		if len(c) < 2 {
			return &NoArgumentsPassed{}
		}
		proc := c[1]
		if err := s.killPoc(proc); err != nil {
			fmt.Println(&PorcKillError{}, proc)
		} else {
			fmt.Println(proc + " killed")
		}
	case "ps":
		if pl, err := s.ps(); err != nil {
			return err
		} else {
			for _, p := range pl {
				name, _ := p.Name()
				pid := p.Pid
				ppid := p.Ppid
				username, _ := p.Username()
				fmt.Printf("Имя: %s, PID: %d, PPID: %d, Пользователь: %s\n", name, pid, ppid, username)
			}
		}
	default:
		return &UnknownCommand{}
	}
	return nil
}

func (s *Shell) ls() error {
	fl, err := os.ReadDir(".")
	if err != nil {
		return err
	}
	for _, f := range fl {
		fmt.Println(f)
	}
	return nil
}

func (s *Shell) ps() ([]*process.Process, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	return processes, nil
}

func (s *Shell) killPoc(PIDorName string) error {
	var cmd *exec.Cmd
	pid, err := strconv.Atoi(PIDorName)
	if err == nil {
		PIDorName = strconv.Itoa(pid)
	}
	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/F", "/IM", PIDorName)
		_, err = cmd.CombinedOutput()
		if err != nil {
			return err
		}
	} else {
		pid, _ = strconv.Atoi(PIDorName)
		proc, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		if err = proc.Kill(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Shell) dirChange(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	s.wd, _ = os.Getwd()
	return nil
}
