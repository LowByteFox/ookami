package main

import (
	"os"
	"os/signal"
	"os/user"
	"path"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
	"github.com/willdonnelly/passwd"
)

func handleProcess(scan Scanner) {
    var split []string
    var processes []*Process
    var reParse bool = false

    for {
        tmp := scan.next()
        if len(tmp) == 0 {
            break
        }

        split = append(split, tmp)
    }

    splitLen := len(split)

    if splitLen == 2 {
        if split[0] == "cd" {
            os.Chdir(split[1])
            return
        }
    }

    for i := 0; i < splitLen; i++ {
        item := split[i]

        if !reParse {
            processes = append(processes, ProcessNew(item))
            reParse = true
            continue
        }

        proc := processes[len(processes)-1]

        if item != ">" && item != "<" && item != "|" && item != "2>" {
            finalItem := ""
            maybeEnv := strings.Split(item, "$")
            for _, i := range maybeEnv {
                if val, exists := os.LookupEnv(i); exists {
                    finalItem += val
                } else {
                    finalItem += i
                }
            }
            proc.args = append(proc.args, finalItem)
        } else {
            hasNext := i + 1 < splitLen
            if item == ">" && hasNext {
                proc.stdout = split[i + 1]
                i++
            } else if item == "<" && hasNext {
                proc.stdin = split[i + 1]
                i++
            } else if item == "2>" && hasNext {
                proc.stderr = split[i + 1]
                i++
            } else if item == "|" && hasNext {
                reParse = false
                proc.pipe = true
            }
        }
    }

    splitLen = len(processes)
    for i := 0; i < splitLen; i++ {
        proc := processes[i]
        proc.prepare()
    }

    var previous *Process = nil

    for i := 0; i < splitLen; i++ {
        proc := processes[i]
        proc.start(previous)
        previous = proc
    }

    previous.end()
}

func main() {
    os.Setenv("SHELL", "ookami")

    args := os.Args[1:]
    if len(args) > 0 {
        startScript(args[0])
        os.Exit(0)
    }

    greet()

    user, _ := user.Current()
    entries, _ := passwd.Parse()
    home_dir := path.Join(entries[user.Username].Home, ".ookami_history")

    l, err := readline.NewEx(&readline.Config{
        Prompt: "> ",
        HistoryFile: home_dir,

        HistorySearchFold: true,
    })

    if err != nil {
        panic(err)
    }

    defer l.Close()

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-interrupt
    }()

    for {
        line, err := l.Readline()
        if err == readline.ErrInterrupt {
            continue
        }
        if err != nil {
            break
        }
        if line == "banner" {
            banner()
            continue
        }
        if line == "exit" {
            break
        }
        scanner := ScannerNew(line)
        handleProcess(scanner)
    }
}
