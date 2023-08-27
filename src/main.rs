use std::env;
use std::path::PathBuf;
use std::process::Child;
use std::process::Command;
use std::process::Stdio;
use std::process::exit;

mod utils;

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
    let mut splited: Vec<Vec<String>> = Vec::new();
    let mut temp = Vec::new();

    for item in parsed {
        if item == "|" {
            splited.push(temp.clone());
            temp.clear();
        } else {
            temp.push(item);
        }
    }
    splited.push(temp.clone());

    let mut iter = splited.iter_mut().peekable();
    let mut previous = None;

    while let Some(args) = iter.next() {
        let command = args.remove(0);

        let stdin = previous.map_or(Stdio::inherit(), |output: Child| {
            Stdio::from(output.stdout.unwrap())
        });

        let stdout = if iter.peek().is_some() {
            Stdio::piped()
        } else {
            Stdio::inherit()
        };

        let output = Command::new(command)
            .args(args)
            .stdin(stdin)
            .stdout(stdout)
            .spawn();

        match output {
            Ok(o) => {
                previous = Some(o);
            }
            Err(e) => {
                previous = None;
                eprintln!("{}", e);
            }
        }
    }

    if let Some(mut last) = previous {
        last.wait().unwrap();
    }
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

    loop {
        let line = readline(&mut rl, "> ");
        if line != "exit" {
            let _ = rl.add_history_entry(line.clone());
            process_prompt(&line);
        } else { break; }
    }

    rl.save_history(home_dir).unwrap();
}
