use crate::model::{Category, NewCategory, UpdateCategory};
use anyhow::Error;
use anyhow::Result;
use futures::future::{ready, BoxFuture};
use salvo::hyper;
use salvo::hyper::server::conn::AddrIncoming;
use salvo::prelude::*;
use std::future::Future;
use std::io;
use std::net::{Ipv4Addr, SocketAddr};
use std::str::FromStr;
use tokio::net::{TcpListener, TcpSocket};
use uuid::Uuid;

use crate::transport::Transport;

pub struct SalvoTransport {
    host: String,
    port: u16,
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

impl Transport for SalvoTransport {
    fn new(host: String, port: u16) -> Self {
        SalvoTransport {
            host: host,
            port: port,
        }
    }
    fn serve(&self) -> BoxFuture<'static, ()> {
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
        let addr = SocketAddr::from((Ipv4Addr::from_str(&self.host).unwrap(), self.port));
        let listener = reuse_listener(addr).expect("couldn't bind to addr");
        let incoming = AddrIncoming::from_listener(listener).unwrap();
        let server = hyper::Server::builder(incoming)
            .http1_only(true)
            .tcp_nodelay(true);
        Box::pin(async {
            let _res = server.serve(Service::new(router)).await;
        })
    }
}
