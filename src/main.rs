use std::env;
use std::path::PathBuf;
use std::process::Command;
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
    let mut args = utils::split_string(prompt);
    let command = args.remove(0);

    let mut child = Command::new(command)
        .args(args)
        .spawn()
        .unwrap();

    let _ = child.wait();
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
