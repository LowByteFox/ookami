package main

import (
	"bufio"
	"strings"
	"unicode"
)

type Scanner struct {
    reader *bufio.Reader
}

func ScannerNew(text string) Scanner {
    return Scanner {
        reader: bufio.NewReader(strings.NewReader(text)),
    }
}

func (s *Scanner) next() string {
    var out strings.Builder
    var escaped bool
    var quoteChar rune

    for {
        c, _, err := s.reader.ReadRune()
        if err != nil {
            break
        }

        if escaped {
            out.WriteRune(c)
            escaped = false
        } else if c == '\\' {
            escaped = true
        } else if quoteChar != 0 {
            if c == quoteChar {
                quoteChar = 0
            } else {
                out.WriteRune(c)
            }
        } else if c == '\'' || c == '"' {
            quoteChar = c
        } else if unicode.IsSpace(c) {
            if out.Len() > 0 {
                break
            }
            continue
        } else {
            out.WriteRune(c)
        }
    }

    if out.Len() > 0 {
        return out.String()
    }
    return ""
}

