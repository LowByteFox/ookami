use std::env;
use std::path::PathBuf;
use std::process::Child;
use std::process::Command;
use std::process::Stdio;
use std::process::exit;

mod utils;
mod process;
mod startup;

use process::Process;

use rustyline::Editor;
use rustyline::history;

fn readline(rl: &mut Editor<(), history::FileHistory>, prompt: &str) -> String {
    let line = rl.readline(prompt);
    return line.unwrap_or(String::from("exit"));
}

fn get_home_directory() -> Option<String> {
    env::var("HOME").ok()
}

fn process_prompt(prompt: &String) {
    let parsed = utils::split_string(prompt);
    let mut parsed_iter = parsed.iter().peekable();
    let mut splited: Vec<Process> = Vec::new();
    let mut re_parse = false;

    while let Some(item) = parsed_iter.next() {
        if !re_parse {
            splited.push(Process::init(item.to_owned()));
            re_parse = true;
            continue;
        }
        let length = splited.len();
        let proc = splited.get_mut(length - 1).unwrap();
        if item != ">" && item != "<" && item != "|" && item != "2>" {
            proc.args.push(item.to_owned());
        } else {
            if item == ">" {
                proc.stdout = if parsed_iter.peek().is_some() {
                    parsed_iter.next().unwrap()
                } else {
                    ""
                }
            } else if item == "<" {
                proc.stdin = if parsed_iter.peek().is_some() {
                    parsed_iter.next().unwrap()
                } else {
                    ""
                }
            } else if item == "2>" {
                proc.stderr = if parsed_iter.peek().is_some() {
                    parsed_iter.next().unwrap()
                } else {
                    ""
                }
            } else if item == "|" {
                re_parse = false;
                proc.pipe = true;
            }
        }
    }

    println!("{:?}", splited);
}

fn main() {
    let mut rl = Editor::<(), history::FileHistory>::new().unwrap();

    let maybe_home = get_home_directory();
    let mut home_dir = match maybe_home {
        Some(ref path) => path.as_str(),
        None => "??"
    };

    if home_dir == "??" {
        eprintln!("I was unable to determine your home directory");
        exit(1); 
    }
    let mut home_path = PathBuf::new();
    home_path.push(home_dir);
    home_path.push(".ookami_history");

    home_dir = home_path.to_str().unwrap();

    if rl.load_history(home_dir).is_err() {
        println!("No previous history.");
    }

    startup::draw();

    loop {
        let line = readline(&mut rl, "> ");
        if line != "exit" {
            let _ = rl.add_history_entry(line.clone());
            process_prompt(&line);
        } else { break; }
    }

    rl.save_history(home_dir).unwrap();
}
