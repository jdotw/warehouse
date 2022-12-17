// #[global_allocator]
// static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

mod model;
mod repository;
mod service;
mod transport;

use dotenvy::dotenv;
use std::env;

use repository::RepositoryBuilder;
use service::Service;
use transport::TransportBuilder;

fn main() {
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let repository = RepositoryBuilder::new()
        .connection_string(database_url)
        .build();
    let service = Service::new(repository);
    let transport = TransportBuilder::new()
        .service(service)
        .port(7878)
        .host("0.0.0.0")
        .build();
    transport.serve_and_await();
}
