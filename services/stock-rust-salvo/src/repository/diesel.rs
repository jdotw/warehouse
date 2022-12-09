extern crate diesel;

pub mod schema;

use crate::model::{Category, NewCategory, UpdateCategory};
use crate::repository::Repository;
use schema::categories::dsl;

use anyhow::Error;
use diesel::prelude::*;
use diesel::r2d2::{ConnectionManager, Pool, PoolError, PooledConnection};
use once_cell::sync::OnceCell;
use uuid::Uuid;

pub struct DieselRepository {
    pub connection_string: String,
}

pub type PgPool = Pool<ConnectionManager<PgConnection>>;

static POOL: OnceCell<PgPool> = OnceCell::new();

fn connect() -> Result<PooledConnection<ConnectionManager<PgConnection>>, PoolError> {
    unsafe { POOL.get_unchecked().get() }
}

impl Repository for DieselRepository {
    fn new(connection_string: String) -> DieselRepository {
        DieselRepository {
            connection_string: connection_string,
        }
    }
    fn build_pool(&self) -> Result<(), Error> {
        let manager = ConnectionManager::<PgConnection>::new(&self.connection_string);
        let result = diesel::r2d2::Pool::builder()
            .max_size(100)
            .min_idle(Some(100))
            .test_on_check_out(false)
            .idle_timeout(None)
            .max_lifetime(None)
            .build(manager);
        Ok(())
    }

    fn get_categories(&self) -> Result<Vec<Category>, Error> {
        let mut connection = connect().unwrap();
        let results = dsl::categories
            .limit(5)
            .load::<Category>(&mut connection)
            .expect("Error loading categories");
        Ok(results)
    }
    fn get_category(&self, id: &Uuid) -> Result<Vec<Category>, Error> {
        let mut connection = connect().unwrap();
        let results = dsl::categories
            .filter(dsl::id.eq(id))
            .limit(1)
            .load::<Category>(&mut connection)
            .expect("Error loading specific category");
        return Ok(results);
    }
    fn create_category(&self, category: &NewCategory) -> Result<Category, Error> {
        let mut connection = connect().unwrap();
        let result = diesel::insert_into(dsl::categories)
            .values(category)
            .get_result::<Category>(&mut connection)
            .expect("Error creating new category");
        Ok(result)
    }
    fn update_category(
        &self,
        category_id: &Uuid,
        category: &UpdateCategory,
    ) -> Result<Category, Error> {
        let mut connection = connect().unwrap();
        let result = diesel::update(dsl::categories.find(category_id))
            .set(category)
            .get_result::<Category>(&mut connection)
            .expect("Error creating new category");
        Ok(result)
    }
    fn delete_category(&self, id: &Uuid) -> Result<(), Error> {
        Ok(())
    }
}
