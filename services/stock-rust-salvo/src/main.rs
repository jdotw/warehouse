// #[global_allocator]
// static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

extern crate diesel;

mod model;
mod repository;
mod server;

use anyhow::Error;
use dotenvy::dotenv;
use model::*;
use repository::diesel::DieselRepository;
use repository::Repository;
use salvo::prelude::*;
use std::env;
use std::thread::available_parallelism;
use uuid::Uuid;

#[handler]
async fn get_categories_canned(res: &mut Response) {
    let results: [Category; 2] = [
        Category {
            id: Uuid::new_v4(),
            name: String::new(),
        },
        Category {
            id: Uuid::new_v4(),
            name: String::new(),
        },
    ];
    res.render(Json(results));
}

#[handler]
async fn get_categories_synch(res: &mut Response) -> Result<(), Error> {
    // res.render(Json(results));
    Ok(())
}

#[handler]
async fn create_category(req: &mut Request, res: &mut Response) -> Result<(), Error> {
    // res.render(Json(result));
    Ok(())
}

#[handler]
async fn update_category(req: &mut Request, res: &mut Response) -> Result<(), Error> {
    // res.render(Json(result));
    Ok(())
}

#[handler]
async fn get_category(req: &Request, res: &mut Response) -> Result<(), Error> {
    // res.render(Json(results));
    Ok(())
}

#[handler]
async fn delete_category(req: &Request, res: &mut Response) -> Result<(), Error> {
    // res.render(Json(results));
    Ok(())
}

async fn serve() {
    let router = Router::new().push(
        Router::with_path("categories")
            .get(get_categories_synch)
            .post(create_category)
            .push(
                Router::with_path("<id>")
                    .get(get_category)
                    .patch(update_category)
                    .delete(delete_category),
            ),
    );

    // Server::new(TcpListener::bind("0.0.0.0:7878"))
    //     .serve(Service::new(router))
    //     .await;

    server::builder().serve(Service::new(router)).await.unwrap();
}

#[tokio::main]
async fn main() {
    console_subscriber::init();
    dotenv().ok();

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let repository = DieselRepository::new(database_url);
    repository.build_pool();

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

    serve().await
}
