use std::process::Child;

#[derive(Debug)]
pub struct Process<'a> {
    pub app: String,
    pub args: Vec<String>,
    pub pipe: bool,
    pub stdout: &'a str,
    pub stderr: &'a str,
    pub stdin: &'a str,
    pub child: Option<Child>
}

impl<'a> Process<'a> {
    pub fn init(app: String) -> Process<'a> {
        let proc = Process {
            app,
            args: Vec::new(),
            pipe: false,
            stdout: "",
            stderr: "",
            stdin: "",
            child: None
        };

        proc
    }
}
