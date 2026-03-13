mod bootstrap;
mod config;
mod protocol;
mod state;
mod transport;

#[tokio::main]
async fn main() {
    env_logger::init();

    log::info!("RedForgeC2 agent starting");
    if let Err(e) = bootstrap::run().await {
        log::error!("agent bootstrap failed: {}", e);
        std::process::exit(1);
    }
}
