package main

import "testing"

func TestReadUntil(t *testing.T) {
	var i = 0
	result := ReadUntil("abcdefg", &i, 'd')
	if result != "abc" {
		t.Fatalf("expected 'abc' but '%s'", result)
	}
	if i != 3 {
		t.Fatalf("expected i = 3, but %d", i)
	}
	i = 0
	if ReadUntil("abcdefg", &i, 'z') != "abcdefg" {
		t.Fatalf("")
	}
}

func TestReadLine(t *testing.T) {
	var i = 0
	result := ReadLine("abcde\nfghij\nk", &i)
	if result != "abcde\n" {
		t.Fatalf("expected 'abcde', but '%s'", result)
	}
	if i != 6 {
		t.Fatalf("expected i= 6, but %d", i)
	}
}

func TestTrimSpaces(t *testing.T) {
	result := TrimSpaces("  hoge foo  hey ")
	expected := []string{"hoge", "foo", "hey"}
	if len(result) != len(expected) {
		t.Fatalf("result count %d doesn't equal expected %d", len(result), len(expected))
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != result[i] {
			t.Fatalf("result[%d] != '%s'", i, expected[i])
		}
	}
	if len(TrimSpaces("    ")) != 0 {
		t.Fatalf("")
	}
}

func TestReadTask(t *testing.T) {
	const kFile = `
task: hoge
	way | way
`
	var i = 0
	name, deps := ReadTask(kFile, &i)
	if name != "task" {
		t.Fatalf("expected 'task', but %s", name)
	}
	if deps[0] != "hoge" {
		t.Fatalf("")
	}
	if i != 12 {
		t.Fatalf("")
	}
}

func TestIncludes(t *testing.T) {
	if !Includes([]string{"a", "b", "c"}, "b") {
		t.Fail()
	}
	if Includes([]string{}, "b") {
		t.Fail()
	}
}

func TestSkipEmptyLineAndComment(t *testing.T) {
	const text = `

# comment
a

# way
b

`
	var i = 0
	SkipEmptyLineAndComment(text, &i)
	if _a := text[i]; _a != 'a' {
		t.Fatalf("%s != 'a'", string(_a))
	}
	ReadLine(text, &i)
	SkipEmptyLineAndComment(text, &i)
	if _b := text[i]; _b != 'b' {
		t.Fatalf("%s != 'b'", string(_b))
	}
	ReadLine(text, &i)
	SkipEmptyLineAndComment(text, &i)
	if i != len(text) {
		t.Fail()
	}
}
func TestParseTasks(t *testing.T) {
	const kFileA = `
t1: t2
	way
t2: t3
	say
t3:
	tay
t4: t1 t2

`
	tasks := ParseTasks(kFileA)
	if len(tasks) != 4 {
		t.Fatalf("task size is not 4: %d", len(tasks))
	}
	if t1, ok := tasks["t1"]; ok {
		if t1.name != "t1" {
			t.Fail()
		}
		if len(t1.dependencies) != 1 {
			t.Fail()
		}
	}
}
func TestBuildCommands(t *testing.T) {
	const kFileA = `
t1: t2
	way
t2: t3
	say
t3:
	tay
t4: t1 t2

`
	tasks := ParseTasks(kFileA)
	result := BuildCommands(&tasks, "t3", "t1")
	expected := []string{"t3", "t2", "t1"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("%s != %s", result[i], expected[i])
		}
	}
	result = [] string{}
	result = BuildCommands(&tasks, "t4", "t3", "t2", "t1")
	expected = []string{"t3", "t2", "t1", "t4"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("%s != %s", result[i], expected[i])
		}
	}

	result = BuildCommands(&tasks, "t2", "t3")
	expected = []string{"t3", "t2"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("%s != %s", result[i], expected[i])
		}
	}
}
