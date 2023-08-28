package main

import (
	"os/user"
	"path"

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

    for i := 0; i < splitLen; i++ {
        item := split[i]

        if !reParse {
            processes = append(processes, ProcessNew(item))
            reParse = true
            continue
        }

        proc := processes[len(processes)-1]

        if item != ">" && item != "<" && item != "|" && item != "2>" {
            proc.args = append(proc.args, item)
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
    greet()

    user, _ := user.Current()
    entries, _ := passwd.Parse()
    home_dir := path.Join(entries[user.Username].Home, ".ookami_history")

    l, err := readline.NewEx(&readline.Config{
        Prompt: "> ",
        HistoryFile: home_dir,
        EOFPrompt: "exit",

        HistorySearchFold: true,
    })

    if err != nil {
        panic(err)
    }

    defer l.Close()
    l.CaptureExitSignal()

    for {
        line, err := l.Readline()
        if err != nil {
            break
        }
        if line == "exit" {
            break
        }
        scanner := ScannerNew(line)
        handleProcess(scanner)
    }
}
