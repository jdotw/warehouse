use crate::model::{Category, NewCategory, UpdateCategory};
use crate::repository::diesel::DieselEngine;

use anyhow::Error;
use uuid::Uuid;

pub mod diesel;

pub trait Engine {
    fn new(connection_string: String) -> Self
    where
        Self: Sized;
    fn get_categories(&self) -> Result<Vec<Category>, Error>;
    fn get_category(&self, id: &Uuid) -> Result<Category, Error>;
    fn create_category(&self, category: &NewCategory) -> Result<Category, Error>;
    fn update_category(
        &self,
        category_id: &Uuid,
        category: &UpdateCategory,
    ) -> Result<Category, Error>;
    fn delete_category(&self, id: &Uuid) -> Result<(), Error>;
}
pub struct Repository {
    engine: Box<dyn Engine + Sync + Send>,
}

impl Repository {
    pub fn get_categories(&self) -> Result<Vec<Category>, Error> {
        self.engine.get_categories()
    }
    pub fn get_category(&self, id: &Uuid) -> Result<Category, Error> {
        self.engine.get_category(id)
    }
    pub fn create_category(&self, category: &NewCategory) -> Result<Category, Error> {
        self.engine.create_category(category)
    }
    pub fn update_category(
        &self,
        category_id: &Uuid,
        category: &UpdateCategory,
    ) -> Result<Category, Error> {
        self.engine.update_category(category_id, category)
    }
    pub fn delete_category(&self, id: &Uuid) -> Result<(), Error> {
        self.engine.delete_category(id)
    }
}

pub struct RepositoryBuilder {
    connection_string: Option<String>,
}

impl RepositoryBuilder {
    pub fn new() -> Self {
        return RepositoryBuilder {
            connection_string: None,
        };
    }
    pub fn connection_string(mut self, connection_string: String) -> Self {
        self.connection_string = Some(connection_string);
        self
    }
    pub fn build(self) -> Repository {
        let conn_string = self.connection_string.clone();
        let engine = DieselEngine::new(conn_string.unwrap());

        Repository {
            engine: Box::new(engine),
        }
    }
}
