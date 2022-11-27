// #[global_allocator]
// static ALLOC: snmalloc_rs::SnMalloc = snmalloc_rs::SnMalloc;

use self::models::*;
use self::schema::categories::dsl as schema;
use diesel::prelude::*;
use diesel::r2d2::{ConnectionManager, Pool, PoolError, PooledConnection};
use once_cell::sync::OnceCell;
use salvo::prelude::*;
use std::str::FromStr;
use std::sync::Arc;
use stock_rust_salvo::*;
use uuid::Uuid;

static POOL: OnceCell<PgPool> = OnceCell::new();

fn connect() -> Result<PooledConnection<ConnectionManager<PgConnection>>, PoolError> {
    unsafe { POOL.get_unchecked().get() }
}

#[handler]
async fn get_categories(res: &mut Response) {
    let results = tokio::task::spawn_blocking(|| {
        let mut connection = connect().unwrap();
        schema::categories
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
    let result = diesel::insert_into(schema::categories)
        .values(&category)
        .get_result::<Category>(&mut connection)
        .expect("Error creating new category");
    res.render(Json(result));
}

#[handler]
async fn update_category(req: &mut Request, res: &mut Response) {
    let mut connection = connect().unwrap();
    let id = req.params().get("id").cloned().unwrap_or_default();
    let id = Uuid::from_str(&id).unwrap();
    let category = req.parse_json::<UpdateCategory>().await.unwrap();
    let result = diesel::update(schema::categories.find(id))
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
    let results = schema::categories
        .filter(schema::id.eq(requested_id))
        .limit(1)
        .load::<Category>(&mut connection)
        .expect("Error loading specific category");
    res.render(Json(results));
}

#[tokio::main()]
async fn main() {
    let pool = create_pool();
    POOL.set(pool).unwrap();

    let router = Arc::new(
        Router::new().push(
            Router::with_path("categories")
                .get(get_categories)
                .post(create_category)
                .push(
                    Router::with_path("<id>")
                        .get(get_category)
                        .patch(update_category),
                ),
        ),
    );

    Server::new(TcpListener::bind("0.0.0.0:7878"))
        .serve(Service::new(router))
        .await;
}
