pub mod schema;

use crate::model::{Category, NewCategory, UpdateCategory};
use crate::repository::Engine;
use schema::categories::dsl;

use anyhow::Error;
use diesel::prelude::*;
use diesel::r2d2::{ConnectionManager, Pool, PoolError, PooledConnection};
use uuid::Uuid;

pub struct DieselEngine {
    pool: Pool<ConnectionManager<PgConnection>>,
}

impl DieselEngine {
    fn connect(&self) -> Result<PooledConnection<ConnectionManager<PgConnection>>, PoolError> {
        self.pool.get()
    }
}

impl Engine for DieselEngine {
    fn new(connection_string: String) -> DieselEngine {
        let manager = ConnectionManager::<PgConnection>::new(connection_string);
        DieselEngine {
            pool: diesel::r2d2::Pool::builder()
                .max_size(100)
                .min_idle(Some(100))
                .test_on_check_out(false)
                .idle_timeout(None)
                .max_lifetime(None)
                .build(manager)
                .unwrap(),
        }
    }

    fn get_categories(&self) -> Result<Vec<Category>, Error> {
        let mut connection = self.connect().unwrap();
        let results = dsl::categories
            .load::<Category>(&mut connection)
            .expect("Error loading categories");
        Ok(results)
    }
    fn get_category(&self, id: &Uuid) -> Result<Category, Error> {
        let mut connection = self.connect().unwrap();
        let results = dsl::categories
            .filter(dsl::id.eq(id))
            .load::<Category>(&mut connection)
            .expect("Error loading specific category");
        Ok(results[0].clone())
    }
    fn create_category(&self, category: &NewCategory) -> Result<Category, Error> {
        let mut connection = self.connect().unwrap();
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
        let mut connection = self.connect().unwrap();
        let result = diesel::update(dsl::categories.find(category_id))
            .set(category)
            .get_result::<Category>(&mut connection)
            .expect("Error creating new category");
        Ok(result)
    }
    fn delete_category(&self, id: &Uuid) -> Result<(), Error> {
        let mut connection = self.connect().unwrap();
        diesel::delete(dsl::categories.find(id))
            .execute(&mut connection)
            .expect("Error deleting category");
        Ok(())
    }
}
