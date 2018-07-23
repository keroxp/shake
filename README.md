# shake
![CircleCI](https://img.shields.io/circleci/project/github/keroxp/memlim.svg?style=flat-square)

A minimal subset of GNU Make focused on running full shell scripts.

![https://gyazo.com/719f3623a591db98e6b7e8b034de5c76.png](https://gyazo.com/719f3623a591db98e6b7e8b034de5c76.png)

## Installation

With go

``
$ go get github.com/keroxp/shake
``

With wget

```bash
$ wget https://github.com/keroxp/shake/releases/download/{{version}}/shake_{{os}}_{{arch}}.zip -O shake.zip
$ unzip shake.zip
$ mv shake /usr/local/bin/shake  
```

## Usage

```bash
$ shake task1 task2 task3
```

## Why I Made It

There are several difference between Shakefile and Makefile. Frankly speaking, Shakefile doesn't have almost all features that Makefile does.    
That is, because Shakefile is designed for describing universal toil tasks in various kind of projects, not just for C.  
For example, I'm using Makefile for building, testing and deploying Docker Image in my production project.  
Obviously, Makefile is a powerful, simple, and universal task runner that has a long history and wide usages.  
However, I found that there are, just in my case, useless features designed for building C/C++ software and critical lacks for describing modern complicated project's task.  
If you have built and deployed Docker application with alpine-linux, you certainly know usefulness about shell script. Shell script is essential component for sipping modern application in arbitrary platforms.   
But Makefile can't run full shell script, especially recognize multiline expressions and variable expansion. 
 
So, I decided to make **Makefile that can run full shell scripts.**

## Author

Yusuke Sakurai (@keroxp)

## LICENSE

MIT