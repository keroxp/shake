package main

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"os/exec"
	"github.com/urfave/cli"
	"os"
	"html/template"
	"io"
	"errors"
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
func Parse(text string) []Task {
	var i int
	var ret []Task
	for ; i < len(text); {
		name := ReadTask(text, &i)
		deps := ReadDependencies(text, &i)
		script := ReadScript(text, &i)
		ret = append(ret, Task{
			name:         name,
			dependencies: deps,
			script:       script,
		})
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

func Action(c *cli.Context) error {
	text, err := ioutil.ReadFile("./Shakefile")
	if err != nil {
		panic(err)
	}
	tasks := Parse(string(text))
	mapp := map[string]Task{}
	for _, v := range tasks {
		mapp[v.name] = v
	}
	var cmds []string
	for i := 0; i < c.NArg(); i++ {
		cmd := c.Args().Get(i)
		task, ok := mapp[cmd]
		if !ok {
			panic(fmt.Sprintf("\"%s\" was not defined", cmd))
		}
		cmds = append(cmds, task.name)
	}
	for i := range tasks {
		val := tasks[i]
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "shake"
	app.Usage = "make by shell"
	app.Action = Action
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
