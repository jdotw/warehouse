use crate::service::Service;
use crate::transport::salvo::SalvoEngine;
use futures::future::BoxFuture;
use once_cell::sync::OnceCell;

pub mod salvo;

pub struct Transport {
    engine: Box<dyn Engine>,
}

pub static SERVICE: OnceCell<Service> = OnceCell::new();

pub trait Engine {
    fn new(host: String, port: u16) -> Self
    where
        Self: Sized;
    fn serve(&self) -> BoxFuture<'static, ()>;
}

impl Transport {
    pub fn serve(&self) -> BoxFuture<'static, ()> {
        self.engine.serve()
    }
}

pub struct TransportBuilder {
    host: String,
    port: u16,
}

impl TransportBuilder {
    pub fn new() -> Self {
        return TransportBuilder {
            host: String::from("0.0.0.0"),
            port: 7878,
        };
    }
    pub fn service(self, service: Service) -> Self {
        let _res = SERVICE.set(service);
        self
    }
    pub fn host(mut self, host: &str) -> Self {
        self.host = String::from(host);
        self
    }
    pub fn port(mut self, port: u16) -> Self {
        self.port = port;
        self
    }
    pub fn build(self) -> Transport {
        let engine = SalvoEngine::new(self.host, self.port);
        Transport {
            engine: Box::new(engine),
        }
    }
}
