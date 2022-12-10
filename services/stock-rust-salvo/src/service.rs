use crate::{
    model::{Category, NewCategory, UpdateCategory},
    repository,
};
use anyhow::Result;
use uuid::Uuid;

pub struct Service {
    repository: repository::Repository,
}

impl Service {
    pub fn new(repository: repository::Repository) -> Self {
        Service {
            repository: repository,
        }
    }
    pub fn get_categories(&self) -> Result<Vec<Category>> {
        self.repository.get_categories()
    }
    pub fn create_category(&self, category: NewCategory) -> Result<Category> {
        self.repository.create_category(&category)
    }
    pub fn update_category(&self, id: Uuid, category: UpdateCategory) -> Result<Category> {
        self.repository.update_category(&id, &category)
    }
    pub fn get_category(&self, id: Uuid) -> Result<Category> {
        self.repository.get_category(&id)
    }
    pub fn delete_category(&self, id: Uuid) -> Result<()> {
        self.repository.delete_category(&id)
    }
}
