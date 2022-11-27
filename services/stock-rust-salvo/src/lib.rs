pub mod models;
pub mod schema;

extern crate diesel;
use diesel::{
    pg::PgConnection,
    r2d2::{ConnectionManager, Pool},
};
use dotenvy::dotenv;
use std::env;

pub type PgPool = Pool<ConnectionManager<PgConnection>>;

pub fn create_pool() -> PgPool {
    // dotenv().ok();
    // let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    // PgConnection::establish(&database_url)
    //     .unwrap_or_else(|_| panic!("Error connecting to {}", database_url))

    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    PgPool::builder()
        .max_size(100)
        .build(ConnectionManager::new(database_url))
        .expect("failed to create db connection_pool")
}
