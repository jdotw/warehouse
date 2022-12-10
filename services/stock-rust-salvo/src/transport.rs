extern crate salvo;

use crate::model::{NewCategory, UpdateCategory};
use crate::service::Service;
use anyhow::{Error, Result};
use futures::future::BoxFuture;
use once_cell::sync::OnceCell;
use salvo::hyper;
use salvo::hyper::server::conn::AddrIncoming;
use salvo::prelude::*;
use std::io;
use std::net::{Ipv4Addr, SocketAddr};
use std::str::FromStr;
use tokio::net::{TcpListener, TcpSocket};
use uuid::Uuid;

pub struct Transport {
    host: String,
    port: u16,
}

pub static SERVICE: OnceCell<Service> = OnceCell::new();

fn service() -> &'static Service {
    SERVICE.get().unwrap()
}

impl Transport {
    pub fn new(service: Service, host: String, port: u16) -> Self {
        let _res = SERVICE.set(service);
        Transport {
            host: host,
            port: port,
        }
    }
    pub fn serve(&self) -> BoxFuture<'static, ()> {
        let router = Router::new().push(
            Router::with_path("categories")
                .get(get_categories)
                .post(create_category)
                .push(
                    Router::with_path("<id>")
                        .get(get_category)
                        .patch(update_category)
                        .delete(delete_category),
                ),
        );
        let addr = SocketAddr::from((Ipv4Addr::from_str(&self.host).unwrap(), self.port));
        let listener = reuse_listener(addr).expect("couldn't bind to addr");
        let incoming = AddrIncoming::from_listener(listener).unwrap();
        let server = hyper::Server::builder(incoming)
            .http1_only(true)
            .tcp_nodelay(true);
        Box::pin(async {
            let _res = server.serve(salvo::Service::new(router)).await;
        })
    }
}

fn reuse_listener(addr: SocketAddr) -> io::Result<TcpListener> {
    let socket = match addr {
        SocketAddr::V4(_) => TcpSocket::new_v4()?,
        SocketAddr::V6(_) => TcpSocket::new_v6()?,
    };

    #[cfg(unix)]
    {
        println!("Using set_reuseport");
        if let Err(e) = socket.set_reuseport(true) {
            eprintln!("error setting SO_REUSEPORT: {}", e);
        }
    }

    println!("reuse_listener");

    socket.set_reuseaddr(true)?;
    socket.bind(addr)?;
    socket.listen(1024)
}

// Handlers

#[handler]
async fn get_categories(res: &mut Response) -> Result<(), Error> {
    let result = service().get_categories().unwrap();
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
