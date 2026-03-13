use log::info;
use std::time::Duration;

fn main() {
    env_logger::init();

    info!("RedForgeC2 agent starting");

    // TODO: Implement bootstrap logic (config, transport, tasks)
    // This stub will be replaced by the full agent runtime.

    loop {
        info!("agent heartbeat (placeholder)");
        std::thread::sleep(Duration::from_secs(30));
    }

    // Note: In a real agent, we would gracefully shutdown on signals.
}
