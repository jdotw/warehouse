use crate::model::{NewCategory, UpdateCategory};
use crate::service::Service;
use crate::transport::{Engine, SERVICE};
use actix_web::{delete, get, patch, post, web, App, HttpResponse, HttpServer, Responder, Result};
use uuid::Uuid;

pub struct ActixEngine {
    host: String,
    port: u16,
}

fn service() -> &'static Service {
    SERVICE.get().unwrap()
}

#[actix_web::main]
async fn serve(host: &str, port: u16) {
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

#[get("/categories")]
async fn get_categories() -> Result<impl Responder> {
    let result = service().get_categories().unwrap();
    Ok(HttpResponse::Ok().json(result))
}

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
