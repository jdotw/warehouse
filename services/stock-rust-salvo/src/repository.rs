use crate::model::{Category, NewCategory, UpdateCategory};

use anyhow::Error;
use uuid::Uuid;

pub mod diesel;

pub trait Repository {
    fn new(connection_string: String) -> Self;
    fn get_categories(&self) -> Result<Vec<Category>, Error>;
    fn get_category(&self, id: &Uuid) -> Result<Vec<Category>, Error>;
    fn create_category(&self, category: &NewCategory) -> Result<Category, Error>;
    fn update_category(
        &self,
        category_id: &Uuid,
        category: &UpdateCategory,
    ) -> Result<Category, Error>;
    fn delete_category(&self, id: &Uuid) -> Result<(), Error>;
}