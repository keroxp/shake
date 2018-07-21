package main

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"github.com/urfave/cli"
	"os"
	"errors"
	"os/exec"
	"github.com/apex/log"
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
		if str[i] == ' ' {
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
func Parse(text string) map[string]Task {
	ret := map[string]Task{}
	for i := 0; i < len(text); {
		name := ReadTask(text, &i)
		deps := ReadDependencies(text, &i)
		script := ReadScript(text, &i)
		ret[name] = Task{
			name:         name,
			dependencies: deps,
			script:       script,
		}
	}
	return ret
}

func ReadTask(text string, indexPtr *int) string {
	if text[*indexPtr] == '\t' {
		panic("tab in root context")
	}
	ret, err := ReadUntilInThisLine(text, indexPtr, ':')
	if err != nil {
		panic("task definition must follow syntax. taskName: (deps1, deps2...)")
	}
	*indexPtr++
	return ret
}

func ReadDependencies(text string, indexPtr *int) []string {
	line := ReadUntil(text, indexPtr, '\n')
	*indexPtr++
	return TrimSpaces(line)
}

func ReadScript(text string, indexPtr *int) string {
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

func Includes(arr *[]string, str string) bool {
	for i := 0; i < len(*arr); i++ {
		if (*arr)[i] == str {
			return true
		}
	}
	return false
}
func BuildCommands(tasks map[string]Task, cmds []string, caller string, dep string) []string {
	task, ok := tasks[dep]
	if !ok {
		panic(fmt.Sprintf("\"%s\" was not defined\n", dep))
	}
	if Includes(&cmds, dep) {
		if len(caller) > 0 {
			fmt.Fprintln(
				os.Stderr,
				fmt.Sprintf("circular dependency %s <- %s was dropped", dep, caller),
			)
		}
		return cmds
	}
	cmds = append(cmds, dep)
	for _, v := range task.dependencies {
		cmds = BuildCommands(tasks, cmds, dep, v)
	}
	return cmds
}
func Action(c *cli.Context) error {
	file := c.String("f")
	text, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	tasks := Parse(string(text))
	var result []string
	for i := 0; i < c.NArg(); i++ {
		cmd := c.Args().Get(i)
		result = BuildCommands(tasks, result, "", cmd)
	}
	log.Debug(fmt.Sprintf("commands: %s", result))
	for i := len(result) - 1; i >= 0; i-- {
		task := result[i]
		script := tasks[task].script
		cmd := exec.Command("/usr/bin/env", "bash", "-c", script)
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
	app.Version = "0.0.1-alpha"
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
