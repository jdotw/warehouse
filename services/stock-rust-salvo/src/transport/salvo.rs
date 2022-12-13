extern crate salvo;

use crate::model::{Category, NewCategory, UpdateCategory};
use crate::repository::sea::query::Query;
use crate::service::Service;
use crate::transport::{Engine, SERVICE};
use anyhow::{Error, Result};
use once_cell::sync::OnceCell;
use salvo::hyper;
use salvo::hyper::server::conn::AddrIncoming;
use salvo::prelude::*;
use sea_orm::{ConnectOptions, Database, DbConn};
use std::env;
use std::io;
use std::net::{Ipv4Addr, SocketAddr};
use std::str::FromStr;
use tokio::net::{TcpListener, TcpSocket};
use uuid::Uuid;

pub struct SalvoEngine {
    host: String,
    port: u16,
}

pub static SEADB: OnceCell<DbConn> = OnceCell::new();

fn service() -> &'static Service {
    SERVICE.get().unwrap()
}

#[tokio::main]
async fn serve(host: &str, port: u16) {
    let db_url = env::var("DATABASE_URL").expect("DATABASE_URL is not set in .env file");
    let mut opt = ConnectOptions::new(db_url);
    opt.max_connections(100);
    let conn = Database::connect(opt).await.unwrap();
    let _setres = SEADB.set(conn);
    let router = Router::new().push(
        Router::with_path("categories")
            .get(get_categories_canned)
            .post(create_category)
            .push(
                Router::with_path("<id>")
                    .get(get_category)
                    .patch(update_category)
                    .delete(delete_category),
            ),
    );
    let addr = SocketAddr::from((Ipv4Addr::from_str(host).unwrap(), port));
    let listener = reuse_listener(addr).expect("couldn't bind to addr");
    let incoming = AddrIncoming::from_listener(listener).unwrap();
    let server = hyper::Server::builder(incoming)
        .http1_only(true)
        .tcp_nodelay(true);
    let _res = server.serve(salvo::Service::new(router)).await;
}

impl Engine for SalvoEngine {
    fn new(host: String, port: u16) -> Self {
        SalvoEngine {
            host: host,
            port: port,
        }
    }
    fn serve_and_await(&self) {
        serve(self.host.as_str(), self.port);
    }
}

fn reuse_listener(addr: SocketAddr) -> io::Result<TcpListener> {
    let socket = match addr {
        SocketAddr::V4(_) => TcpSocket::new_v4()?,
        SocketAddr::V6(_) => TcpSocket::new_v6()?,
    };

    #[cfg(unix)]
    {
        if let Err(e) = socket.set_reuseport(true) {
            eprintln!("error setting SO_REUSEPORT: {}", e);
        }
    }

    socket.set_reuseaddr(true)?;
    socket.bind(addr)?;
    socket.listen(1024)
}

// Handlers

fn dbconn() -> &'static DbConn {
    unsafe { SEADB.get_unchecked() }
}

#[handler]
async fn get_categories(res: &mut Response) -> Result<(), Error> {
    // let result = service().get_categories().unwrap();
    let result = Query::find_categories_in_page(dbconn(), 1, 64)
        .await
        .unwrap();
    res.render(Json(result.0));
    Ok(())
}

#[handler]
async fn get_categories_canned(res: &mut Response) -> Result<(), Error> {
    let mut result = Vec::with_capacity(2);
    result.push(Category {
        id: Uuid::new_v4(),
        name: "Canned 1".to_string(),
    });
    result.push(Category {
        id: Uuid::new_v4(),
        name: "Canned 2".to_string(),
    });
    res.render(Json(result));
    Ok(())
}

#[handler]
async fn create_category(req: &mut Request, res: &mut Response) -> Result<(), Error> {
    let category = req
        .parse_json::<NewCategory>()
        .await
        .expect("failed to parse category");
    let result = service().create_category(category).unwrap();
    res.render(Json(result));
    Ok(())
}

#[handler]
async fn update_category(req: &mut Request, res: &mut Response) -> Result<(), Error> {
    let id = Uuid::from_str(req.params().get("id").unwrap()).unwrap();
    let category = req
        .parse_json::<UpdateCategory>()
        .await
        .expect("failed to parse category");
    let result = service().update_category(id, category).unwrap();
    res.render(Json(result));
    Ok(())
}

#[handler]
async fn get_category(req: &Request, res: &mut Response) -> Result<(), Error> {
    let id = Uuid::from_str(req.params().get("id").unwrap()).unwrap();
    let result = service().get_category(id).unwrap();
    res.render(Json(result));
    Ok(())
}

#[handler]
async fn delete_category(req: &Request, res: &mut Response) -> Result<(), Error> {
    let id = Uuid::from_str(req.params().get("id").unwrap()).unwrap();
    let result = service().delete_category(id).unwrap();
    res.render(Json(result));
    Ok(())
}
