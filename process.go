package main

import (
	"io"
	"os"
	"os/exec"
	"syscall"
)

type Process struct {
    app string
    args []string
    pipe bool
    stdout string
    stdin string
    stderr string
    pipe_pipe_reader *io.PipeReader
    pipe_pipe_writer *io.PipeWriter
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

    if p.pipe {
        p.pipe_pipe_reader, p.pipe_pipe_writer = io.Pipe()
    }

    if len(p.stdout) > 0 {
        file, _ := os.OpenFile(p.stdout, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
        p.stdout_pipe = file
    }

    if len(p.stdin) > 0 {
        file, _ := os.OpenFile(p.stdin, os.O_RDONLY, 0644)
        p.stdin_pipe = file
    }

    if len(p.stderr) > 0 {
        file, _ := os.OpenFile(p.stderr, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
        p.stderr_pipe = file
    }
}

func (p *Process) start(previous *Process) {
    p.proc = exec.Command(p.app, p.args...)

    if previous != nil {
        if previous.pipe {
            p.proc.Stdin = previous.pipe_pipe_reader
        }
    }

    if p.stdout_pipe != nil {
        p.proc.Stdout = p.stdout_pipe
    }

    if p.stderr_pipe != nil {
        p.proc.Stderr = p.stderr_pipe
    }

    if p.stdin_pipe != nil {
        stdin, err := p.proc.StdinPipe()
        if err != nil {
            panic(err)
        }

        io.Copy(stdin, p.stdin_pipe)
        stdin.Close()
    }

    if p.pipe {
        p.proc.Stdout = p.pipe_pipe_writer
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

    p.proc.Start()
    if previous != nil {
        if previous.pipe {
            previous.end()
            previous.pipe_pipe_writer.Close()
        }
    }
}

func (p *Process) end() {
    err := p.proc.Wait()
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            exitStatus := exitErr.Sys().(syscall.WaitStatus).ExitStatus()
            if exitStatus == -1 {
            } else {
                panic(err)
            }
        }
    }

    if p.stdout_pipe != nil {
        p.stdout_pipe.Close()
    }

    if p.stderr_pipe != nil {
        p.stderr_pipe.Close()
    }

    if p.stdin_pipe != nil {
        p.stdin_pipe.Close()
    } 
}
