use actix::System;
use log::warn;
use log::Log;
use multi_log;
use systemd::journal;

fn main() {
    let a: Box<dyn Log> = Box::new(
        pretty_env_logger::formatted_timed_builder()
            .filter_level(log::LevelFilter::Debug)
            .build(),
    );
    let b: Box<dyn Log> = Box::new(journal::JournalLog);

    let loggers: Vec<Box<dyn Log>> = vec![a, b];

    multi_log::MultiLogger::init(loggers, log::Level::Debug).unwrap();

    let system = System::new();
    system.block_on(async {  });

    system.run().unwrap();
}
