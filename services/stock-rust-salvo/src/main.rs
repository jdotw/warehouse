// #[global_allocator]
// static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

mod model;
mod repository;
mod service;
mod transport;

use dotenvy::dotenv;
use repository::diesel::DieselRepository;
use repository::Repository;
use service::DefaultService;
use service::Service;
use std::env;
use transport::salvo::SalvoTransport;
use transport::Transport;

#[tokio::main]
async fn main() {
    console_subscriber::init();
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let repository = DieselRepository::new(database_url);
    let service = DefaultService::new(repository);
    let transport = SalvoTransport::new(service, "0.0.0.0".to_string(), 7878);
    transport.serve().await
}
