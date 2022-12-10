use futures::future::BoxFuture;

pub mod salvo;

pub trait Transport {
    fn new(host: String, port: u16) -> Self;
    fn serve(&self) -> BoxFuture<'static, ()>;
}
