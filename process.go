package main

import (
	"io"
	"os"
	"os/exec"
)

type Process struct {
    app string
    args []string
    pipe bool
    stdout string
    stdin string
    stderr string
    stdout_pipe *os.File
    stdin_pipe *os.File
    stderr_pipe *os.File
    run_stdout_pipe *io.ReadCloser
    proc *exec.Cmd
}

func ProcessNew(app string) *Process {
    return &Process{
        app: app,
        pipe: false,
        stdout: "",
        stderr: "",
        stdin: "",
    }
}

func (p *Process) prepare() {
    if len(p.stdout) > 0 && p.pipe {
        panic("Cannot pipe stdout to file and also to another process")
    }

    if len(p.stdout) > 0 {
        if _, err := os.Stat(p.stdout); err == nil {
            file, _ := os.Open(p.stdout)
            p.stdout_pipe = file
        } else {
            file, err := os.Create(p.stdout)
            if err != nil {
                panic(err)
            }

            p.stdout_pipe = file
        }
    }
}

func (p *Process) start(previous *Process) {
    p.proc = exec.Command(p.app, p.args...)

    if previous != nil {
        if previous.pipe {
            p.proc.Stdin = *previous.run_stdout_pipe
        }
    }

    if p.stdout_pipe != nil {
        p.proc.Stdout = p.stdout_pipe
    }

    if p.pipe {
        pipe, err := p.proc.StdoutPipe()
        if err != nil {
            panic(err)
        }

        p.run_stdout_pipe = &pipe
    }

    if p.proc.Stdout == nil {
        p.proc.Stdout = os.Stdout
    }
    if p.proc.Stdin == nil {
        p.proc.Stdin = os.Stdin
    }
    if p.proc.Stderr == nil {
        p.proc.Stderr = os.Stderr
    }

    p.proc.Run()
}
