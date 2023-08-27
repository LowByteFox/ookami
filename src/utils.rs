use std::str::Chars;

struct CommandSplitter<'a> {
    chars: Chars<'a>,
}

impl<'a> CommandSplitter<'a> {
    fn new(command: &'a str) -> Self {
        Self {
            chars: command.chars(),
        }
    }
}

impl<'a> Iterator for CommandSplitter<'a> {
    type Item = String;
    fn next(&mut self) -> Option<Self::Item> {
        let mut out = String::new();
        let mut escaped = false;
        let mut quote_char = None;
        while let Some(c) = self.chars.next() {
            if escaped {
                out.push(c);
                escaped = false;
            } else if c == '\\' {
                escaped = true;
            } else if let Some(qc) = quote_char {
                if c == qc {
                    quote_char = None;
                } else {
                    out.push(c);
                }
            } else if c == '\'' || c == '"' {
                quote_char = Some(c);
            } else if c.is_whitespace() {
                if !out.is_empty() {
                    return Some(out);
                } else {
                    continue;
                }
            } else {
                out.push(c);
            }
        }

        if !out.is_empty() {
            Some(out)
        } else {
            None
        }
    }
}

pub fn split_string(string: &String) -> Vec<String> {
    let ret = CommandSplitter::new(string);
    let mut args = Vec::new();
    for item in ret {
        args.push(item);
    }
    args
}
