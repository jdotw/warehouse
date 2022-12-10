// #[global_allocator]
// static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

extern crate diesel;

mod model;
mod repository;
mod transport;

use anyhow::Error;
use dotenvy::dotenv;
use model::*;
use repository::diesel::DieselRepository;
use repository::Repository;
use salvo::prelude::*;
use std::env;
use std::thread::available_parallelism;
use transport::salvo::SalvoTransport;
use transport::Transport;
use uuid::Uuid;

#[tokio::main]
async fn main() {
    console_subscriber::init();
    dotenv().ok();

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let _repository = DieselRepository::new(database_url);

    let size = available_parallelism().map(|n| n.get()).unwrap_or(16);

    println!("Cores: {size}");
    // for _ in 1..size {
    //     let rt = runtime::Runtime::new().unwrap();
    //     rt.block_on(serve());
    // }
    // println!("Started http server: 127.0.0.1:7878");

    // let rt = runtime::Builder::new_multi_thread()
    //     .worker_threads(10)
    //     .max_blocking_threads(4096)
    //     .enable_all()
    //     .build()
    //     .unwrap();
    // rt.block_on(serve());

    let transport = SalvoTransport::new("0.0.0.0".to_string(), 7878);
    transport.serve().await
}
