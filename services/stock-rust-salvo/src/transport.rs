use crate::service::Service;
use futures::future::BoxFuture;

pub mod salvo;

pub trait Transport<'a> {
    fn new(service: impl Service + Send + Sync + 'static, host: String, port: u16) -> Self;
    fn serve(&self) -> BoxFuture<'static, ()>;
}
