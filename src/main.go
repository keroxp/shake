package main

import (
	"fmt"
	"bytes"
	"io/ioutil"
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
	for i = 0; i < len(text); {
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
	ret := ReadUntil(text, indexPtr, ':')
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

func main() {
	text, err := ioutil.ReadFile("./Shakefile")
	if err != nil {
		panic(err)
	}
	tasks := Parse(string(text))
	for i := range tasks {
		val := tasks[i]
		fmt.Printf("name: %s\n", val.name)
		fmt.Printf("deps: %s\n", val.dependencies)
		fmt.Printf("script:\n%s\n", val.script)
		out, _ := exec.Command("/bin/bash", "-c", val.script).Output()
		fmt.Println(string(out))
	}
}
