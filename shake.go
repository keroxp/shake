package main

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"github.com/urfave/cli"
	"os"
	"errors"
	"github.com/apex/log"
	"os/exec"
)

type Task struct {
	name         string
	script       string
	dependencies []string
}

func TrimSpaces(str string) []string {
	var ret []string
	var buf bytes.Buffer
	for i := 0; i < len(str); i++ {
		if str[i] == ' ' || str[i] == '\t' {
			if buf.Len() > 0 {
				ret = append(ret, buf.String())
				buf.Reset()
			}
			continue
		}
		buf.WriteByte(str[i])
	}
	if buf.Len() > 0 {
		ret = append(ret, buf.String())
	}
	return ret
}

func ReadUntilInThisLine(str string, indexPtr *int, char byte) (string, error) {
	var buf bytes.Buffer
	for ; *indexPtr < len(str) && str[*indexPtr] != char; *indexPtr++ {
		c := str[*indexPtr]
		if c == '\n' {
			return "", errors.New(fmt.Sprintf("'%b' was not found in this line", char))
		}
		buf.WriteByte(c)
	}
	return buf.String(), nil
}
func ReadUntil(str string, indexPtr *int, char byte) string {
	var buf bytes.Buffer
	for ; *indexPtr < len(str) && str[*indexPtr] != char; *indexPtr++ {
		buf.WriteByte(str[*indexPtr])
	}
	return buf.String()
}
func ReadLine(str string, indexPtr *int) string {
	line := ReadUntil(str, indexPtr, '\n')
	*indexPtr++
	return line + "\n"
}
func ParseTasks(text string) map[string]Task {
	ret := map[string]Task{}
	for i := 0; i < len(text); {
		name, deps := ReadTask(text, &i)
		script := ReadScript(text, &i)
		ret[name] = Task{
			name:         name,
			dependencies: deps,
			script:       script,
		}
	}
	return ret
}

func SkipEmptyLineAndComment(text string, indexPtr *int) {
	if len(text) <= *indexPtr {
		return
	}
	if text[*indexPtr] == '\n' || text[*indexPtr] == '#' {
		ReadLine(text, indexPtr)
		SkipEmptyLineAndComment(text, indexPtr)
	}
}
func ReadTask(text string, indexPtr *int) (string, []string) {
	SkipEmptyLineAndComment(text, indexPtr)
	if text[*indexPtr] == '\t' {
		panic("tab in root context")
	}
	ret, err := ReadUntilInThisLine(text, indexPtr, ':')
	if err != nil {
		panic("task definition must follow syntax. taskName: (deps1, deps2...)")
	}
	*indexPtr++
	line := ReadUntil(text, indexPtr, '\n')
	*indexPtr++
	deps := TrimSpaces(line)
	return ret, deps
}

func ReadScript(text string, indexPtr *int) string {
	SkipEmptyLineAndComment(text, indexPtr)
	var buf bytes.Buffer
	for ; ; {
		if len(text) <= *indexPtr || text[*indexPtr] != '\t' {
			break
		}
		line := ReadLine(text, indexPtr)
		buf.WriteString(line)

	}
	return buf.String()
}

func Includes(arr []string, str string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return true
		}
	}
	return false
}
func BuildCommandsInternal(tasks *map[string]Task, cmds []string, caller string, deps ...string) []string {
	for _, dep := range deps {
		task, ok := (*tasks)[dep]
		if !ok {
			panic(fmt.Sprintf("\"%s\" was not defined", dep))
		}
		if Includes(cmds, dep) {
			if len(caller) > 0 {
				fmt.Fprintln(
					os.Stderr,
					fmt.Sprintf("circular dependency %s <- %s was dropped", dep, caller),
				)
			}
			continue
		}
		cmds = BuildCommandsInternal(tasks, cmds, dep, task.dependencies...)
		cmds = append(cmds, dep)
	}
	return cmds
}
func BuildCommands(tasks *map[string]Task, deps ...string) []string {
	var result []string
	return BuildCommandsInternal(tasks, result, "", deps...)
}

func Action(c *cli.Context) error {
	file := c.String("f")
	text, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	tasks := ParseTasks(string(text))
	var cmds []string
	for i := 0; i < c.NArg(); i++ {
		cmds = append(cmds, c.Args().Get(i))
	}
	result := BuildCommands(&tasks, cmds...)
	log.Debug(fmt.Sprintf("commands: %s", result))
	for i := 0; i < len(result); i++ {
		task := result[i]
		script := tasks[task].script
		cmd := exec.Command("/usr/bin/env", "bash", "-c", script)
		if err != nil {
			return err
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "shake"
	app.Usage = "make by shell"
	app.Version = "0.1.0-alpha"
	app.Action = Action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "f",
			Value: "Shakefile",
			Usage: "specify input Shakefile `FILE`",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s", err.Error()))
		os.Exit(1)
	}
}
