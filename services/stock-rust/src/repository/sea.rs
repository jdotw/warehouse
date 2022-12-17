pub mod entity;
pub mod mutation;
pub mod query;

use crate::model::{Category, NewCategory, UpdateCategory};
use crate::repository::sea::query::Query;
use crate::repository::Engine;
use sea_orm::{ConnectOptions, Database, DatabaseConnection, DbConn};
use tokio::runtime::{Handle, Runtime};

use anyhow::{anyhow, Error};
use once_cell::sync::OnceCell;
use std::env;
use uuid::Uuid;

pub static SEA_DB: OnceCell<DatabaseConnection> = OnceCell::new();

fn connect(connection_string: &String) -> &DatabaseConnection {
    SEA_DB.get_or_init(|| {
        let mut opt = ConnectOptions::new(connection_string.clone());
        opt.max_connections(100);

        // Execute the future, blocking the current thread until completion
        let conn = Handle::current().block_on(Database::connect(opt)).unwrap();
        println!("Hey conn is {:?}", conn);
        return conn;
    })
}

pub struct SeaEngine {
    connection_string: String,
}

impl Engine for SeaEngine {
    fn new(connection_string: String) -> SeaEngine {
        SeaEngine {
            connection_string: connection_string,
        }
    }

    fn get_categories(&self) -> Result<Vec<Category>, Error> {
        let conn = connect(&self.connection_string);
        let result = Handle::current()
            .block_on(Query::find_categories_in_page(&conn, 0, 64))
            .unwrap();
        let mut categories = Vec::new();
        for entity in result.0.iter() {
            let category = entity.to_category();
            categories.push(category);
        }
        Ok(categories)
    }
    fn get_category(&self, id: &Uuid) -> Result<Category, Error> {
        Err(anyhow!("Not implemented"))
    }
    fn create_category(&self, category: &NewCategory) -> Result<Category, Error> {
        Err(anyhow!("Not implemented"))
    }
    fn update_category(
        &self,
        category_id: &Uuid,
        category: &UpdateCategory,
    ) -> Result<Category, Error> {
        Err(anyhow!("Not implemented"))
    }
    fn delete_category(&self, id: &Uuid) -> Result<(), Error> {
        Err(anyhow!("Not implemented"))
    }
}
