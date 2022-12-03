// #[global_allocator]
// static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

extern crate diesel;

mod models;
mod schema;

use anyhow::Error;
use diesel::prelude::*;
use diesel::r2d2::{ConnectionManager, Pool, PoolError, PooledConnection};
use dotenvy::dotenv;
use models::*;
use once_cell::sync::OnceCell;
use salvo::prelude::*;
use schema::categories::dsl::*;
use tokio::runtime;

use std::env;
use std::str::FromStr;
use std::sync::Arc;
use std::thread::available_parallelism;
use uuid::Uuid;

mod server;

static POOL: OnceCell<PgPool> = OnceCell::new();
pub type PgPool = Pool<ConnectionManager<PgConnection>>;

fn connect() -> Result<PooledConnection<ConnectionManager<PgConnection>>, PoolError> {
    unsafe { POOL.get_unchecked().get() }
}

fn build_pool(database_url: &str, size: u32) -> Result<PgPool, PoolError> {
    let manager = ConnectionManager::<PgConnection>::new(database_url);
    diesel::r2d2::Pool::builder()
        .max_size(size)
        .min_idle(Some(size))
        .test_on_check_out(false)
        .idle_timeout(None)
        .max_lifetime(None)
        .build(manager)
}

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
    let mut connection = connect().unwrap();
    let results = categories
        .limit(5)
        .load::<Category>(&mut connection)
        .expect("Error loading categories");
    res.render(Json(results));
    Ok(())
}

#[handler]
async fn get_categories_spawn_blocking(res: &mut Response) {
    let results = tokio::task::spawn_blocking(|| {
        let mut connection = connect().unwrap();
        categories
            .limit(5)
            .load::<Category>(&mut connection)
            .expect("Error loading categories")
    })
    .await
    .unwrap();
    res.render(Json(results));
}

#[handler]
async fn create_category(req: &mut Request, res: &mut Response) {
    let mut connection = connect().unwrap();
    let category = req.parse_json::<NewCategory>().await.unwrap();
    let result = diesel::insert_into(categories)
        .values(&category)
        .get_result::<Category>(&mut connection)
        .expect("Error creating new category");
    res.render(Json(result));
}

#[handler]
async fn update_category(req: &mut Request, res: &mut Response) {
    let mut connection = connect().unwrap();
    let category_id = req.params().get("id").cloned().unwrap_or_default();
    let category_id = Uuid::from_str(&category_id).unwrap();
    let category = req.parse_json::<UpdateCategory>().await.unwrap();
    let result = diesel::update(categories.find(category_id))
        .set(&category)
        .get_result::<Category>(&mut connection)
        .expect("Error creating new category");
    res.render(Json(result));
}

#[handler]
async fn get_category(req: &Request, res: &mut Response) {
    let requested_id = req.params().get("id").cloned().unwrap_or_default();
    let requested_id = Uuid::from_str(&requested_id).unwrap();
    let mut connection = connect().unwrap();
    let results = categories
        .filter(id.eq(requested_id))
        .limit(1)
        .load::<Category>(&mut connection)
        .expect("Error loading specific category");
    res.render(Json(results));
}

async fn serve() {
    let router = Router::new().push(
        Router::with_path("categories")
            .get(get_categories_synch)
            .post(create_category)
            .push(
                Router::with_path("<id>")
                    .get(get_category)
                    .patch(update_category),
            ),
    );

    // Server::new(TcpListener::bind("0.0.0.0:7878"))
    //     .serve(Service::new(router))
    //     .await;

    server::builder().serve(Service::new(router)).await.unwrap();
}

fn main() {
    console_subscriber::init();
    dotenv().ok();

    let size = available_parallelism().map(|n| n.get()).unwrap_or(16);

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    POOL.set(
        build_pool(&database_url, size as u32)
            .unwrap_or_else(|_| panic!("Error connecting to {}", &database_url)),
    )
    .ok();

    println!("Cores: {size}");
    // for _ in 1..size {
    //     let rt = runtime::Runtime::new().unwrap();
    //     rt.block_on(serve());
    // }
    // println!("Started http server: 127.0.0.1:7878");

    let rt = runtime::Builder::new_multi_thread()
        .worker_threads(10)
        .max_blocking_threads(4096)
        .enable_all()
        .build()
        .unwrap();

    rt.block_on(serve());
}
