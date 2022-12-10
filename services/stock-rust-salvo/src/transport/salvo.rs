use crate::model::Category;
use crate::service::Service;
use anyhow::Error;
use anyhow::Result;
use futures::future::BoxFuture;
use salvo::hyper;
use salvo::hyper::server::conn::AddrIncoming;
use salvo::prelude::*;
use std::io;
use std::marker::Send;
use std::net::{Ipv4Addr, SocketAddr};
use std::str::FromStr;
use tokio::net::{TcpListener, TcpSocket};
use uuid::Uuid;

use crate::transport::Transport;

pub struct SalvoTransport<'a> {
    host: String,
    port: u16,
    service: Box<dyn Service + Sync + Send + 'a>,
    handler: Option<GetCategoriesHandler<'a>>,
}

impl SalvoTransport<'_> {}

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

struct GetCategoriesHandler<'a> {
    transport: &'a SalvoTransport<'a>,
}

#[async_trait]
impl Handler for GetCategoriesHandler<'static> {
    async fn handle(
        &self,
        req: &mut Request,
        depot: &mut Depot,
        res: &mut Response,
        ctrl: &mut FlowCtrl,
    ) {
        let results = self.transport.service.get_categories().unwrap();
        res.render(Json(results));
        ctrl.call_next(req, depot, res).await;
    }
}

#[handler]
async fn get_categories(res: &mut Response) -> Result<(), Error> {
    // let results = service.get_categories().unwrap();
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

impl<'a> Transport<'a> for SalvoTransport<'static> {
    fn new(service: impl Service + Send + Sync + 'static, host: String, port: u16) -> Self {
        let transport = SalvoTransport {
            host: host,
            port: port,
            service: Box::new(service),
            handler: None,
        };
        transport.handler = Some(GetCategoriesHandler {
            transport: &transport,
        });
        transport
    }
    fn serve(&self) -> BoxFuture<'static, ()> {
        let router = Router::new().push(
            Router::with_path("categories")
                .get(self.handler.unwrap())
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
