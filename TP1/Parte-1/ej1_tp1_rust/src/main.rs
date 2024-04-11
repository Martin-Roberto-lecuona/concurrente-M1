use nix::sys::wait::{wait};
use nix::unistd::{fork, ForkResult};
use std::collections::HashMap;
use std::io::{Error, ErrorKind};
use std::process::id;
use std::{os::unix::process::parent_id, process::Command};

fn main()
{
    let process_tree: HashMap<char, Vec<char>> = HashMap::from([
        ('A', vec!['B', 'C', 'D']),
        ('B', vec!['E', 'F']),
        ('D', vec!['G']),
        ('F', vec!['H', 'I']),
        ('I', vec!['J']),
    ]);

    let seconds = 10;

    if let Err(e) = create_fork(&process_tree, &'A', &seconds)
    {
        eprintln!("Error: {}", e);
        std::process::exit(1);
    }
}

fn create_fork(
    process_tree: &HashMap<char, Vec<char>>,
    initial_node: &char,
    seconds: &u32,
) -> Result<(), Error>
{
    println!("Node {} id {} pid {}", initial_node, id(), parent_id());
    sleep_seconds(&1);

    match process_tree.get(&initial_node)
     {
        Some(children) =>
        {
            for &next_child in children
            {
                match unsafe
                {
                  fork()
                }
                {
                    Ok(ForkResult::Parent { .. }) => {}
                    Ok(ForkResult::Child) =>
                    {
                        if let Err(e) = create_fork(process_tree, &next_child, seconds)
                        {
                            eprintln!("Error in child node: {}", e);
                            std::process::exit(1)
                        }
                        std::process::exit(0)
                    }
                    Err(e) => return Err(Error::new(ErrorKind::Other, format!("Error in fork: {}", e))),
                }
            }
            for _ in 0..children.len()
            {
                let _ = wait();
            }
        }
        None => sleep_seconds(seconds),
    }

    println!("Ending node: {}", initial_node);
    Ok(())
}

fn sleep_seconds(seconds: &u32)
{
    Command::new("sleep")
        .arg(&seconds.to_string())
        .output()
        .expect("failure to execute process");
}