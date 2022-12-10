use crate::{
    model::{Category, NewCategory, UpdateCategory},
    repository,
};
use anyhow::Result;
use anyhow::*;
use uuid::Uuid;

pub trait Service {
    fn new(repository: impl repository::Repository + 'static) -> Self
    where
        Self: Sized;
    fn get_categories(&self) -> Result<Vec<Category>>;
    fn create_category(&self, category: NewCategory) -> Result<Category>;
    fn update_category(&self, id: Uuid, category: UpdateCategory) -> Result<Category>;
    fn get_category(&self, id: Uuid) -> Result<Category>;
    fn delete_category(&self, id: Uuid) -> Result<()>;
}

pub struct DefaultService {
    repository: Box<dyn repository::Repository + 'static>,
}

impl Service for DefaultService {
    fn new(repository: impl repository::Repository + 'static) -> Self {
        DefaultService {
            repository: Box::new(repository),
        }
    }
    fn get_categories(&self) -> Result<Vec<Category>> {
        self.repository.get_categories()
    }
    fn create_category(&self, category: NewCategory) -> Result<Category> {
        self.repository.create_category(&category)
    }
    fn update_category(&self, id: Uuid, category: UpdateCategory) -> Result<Category> {
        self.repository.update_category(&id, &category)
    }
    fn get_category(&self, id: Uuid) -> Result<Category> {
        // Ok(self.repository.get_category(&id).unwrap().first().unwrap())
        Err(anyhow!("not implemented"))
    }
    fn delete_category(&self, id: Uuid) -> Result<()> {
        Err(anyhow!("not implemented"))
    }
}

unsafe impl Sync for DefaultService {}
unsafe impl Send for DefaultService {}
