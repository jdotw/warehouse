use crate::model::{Category, NewCategory, UpdateCategory};
use crate::repository::sea::query::Query;
use crate::service::Service;
use crate::transport::{Engine, SERVICE};
use actix_web::cookie::time::Duration;
use actix_web::{delete, get, patch, post, web, App, HttpResponse, HttpServer, Responder, Result};
use once_cell::sync::OnceCell;
use sea_orm::{ConnectOptions, Database, DbConn};
use std::env;
use uuid::Uuid;

pub struct ActixEngine {
    host: String,
    port: u16,
}

fn service() -> &'static Service {
    SERVICE.get().unwrap()
}

pub static SEADB: OnceCell<DbConn> = OnceCell::new();

#[actix_web::main]
async fn serve(host: &str, port: u16) {
    let db_url = env::var("DATABASE_URL").expect("DATABASE_URL is not set in .env file");
    let mut opt = ConnectOptions::new(db_url);
    opt.max_connections(100);
    opt.min_connections(100);
    let conn = Database::connect(opt).await.unwrap();
    let _ = SEADB.set(conn);

    let _server = HttpServer::new(|| {
        App::new()
            .service(get_categories)
            .service(create_category)
            .service(get_category)
            .service(update_category)
            .service(delete_category)
    })
    .bind((host, port))
    .unwrap()
    .run()
    .await;
}

impl Engine for ActixEngine {
    fn new(host: String, port: u16) -> Self {
        ActixEngine {
            host: host.clone(),
            port: port,
        }
    }

    fn serve_and_await(&self) {
        serve(self.host.as_str(), self.port);
    }
}

fn dbconn() -> &'static DbConn {
    unsafe { SEADB.get_unchecked() }
}

#[get("/categories")]
async fn get_categories() -> Result<impl Responder> {
    let mut result = Vec::with_capacity(2);
    result.push(Category {
        id: Uuid::new_v4(),
        name: "Canned 1".to_string(),
    });
    result.push(Category {
        id: Uuid::new_v4(),
        name: "Canned 2".to_string(),
    });
    Ok(HttpResponse::Ok().json(result))
}

// #[get("/categories")]
// async fn get_categories() -> Result<impl Responder> {
//     // let result = service().get_categories().unwrap();
//     // let result = Query::find_categories_in_page(dbconn(), 0, 64)
//     //     .await
//     //     .unwrap();
//     let result = Query::find_categories_in_page(dbconn(), 0, 64)
//         .await
//         .unwrap();
//     Ok(HttpResponse::Ok().json(result.0.as_slice()))
// }

#[post("/categories")]
async fn create_category(item: web::Json<NewCategory>) -> Result<impl Responder> {
    let result = service().create_category(item.0).unwrap();
    Ok(HttpResponse::Created().json(result))
}

#[patch("/categories/{categoryid}")]
async fn update_category(
    path: web::Path<Uuid>,
    item: web::Json<UpdateCategory>,
) -> Result<impl Responder> {
    let category_id = path.into_inner();
    let result = service().update_category(category_id, item.0).unwrap();
    Ok(HttpResponse::Ok().json(result))
}

#[get("/categories/{category_id}")]
async fn get_category(path: web::Path<Uuid>) -> Result<impl Responder> {
    let category_id = path.into_inner();
    let result = service().get_category(Uuid::from(category_id)).unwrap();
    Ok(HttpResponse::Ok().json(result))
}

#[delete("/categories/{category_id}")]
async fn delete_category(path: web::Path<Uuid>) -> Result<impl Responder> {
    let category_id = path.into_inner();
    service().delete_category(category_id).unwrap();
    Ok(HttpResponse::NoContent())
}
