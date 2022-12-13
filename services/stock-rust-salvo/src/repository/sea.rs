pub mod mutation;
pub mod query;

// use crate::model::{Category, NewCategory, UpdateCategory};
// use crate::repository::sea::query::{Query}
// use crate::repository::Engine;
// use sea_orm::DbConn;

// use anyhow::Error;
// use uuid::Uuid;

// pub struct SeaEngine {
//   db: DbConn
//     // pool: Pool<ConnectionManager<PgConnection>>,
// }

// impl SeaEngine {
//     // fn connect(&self) -> Result<PooledConnection<ConnectionManager<PgConnection>>, PoolError> {
//     //     self.pool.get()
//     // }
// }

// impl Engine for SeaEngine {
//     fn new(connection_string: String) -> SeaEngine {
//         SeaEngine {}
//     }

//     fn get_categories(&self) -> Result<Vec<Category>, Error> {
//         let result = Query::find_categories_in_page(&self.db, 0, 64)
//         Ok(())
//     }
//     fn get_category(&self, id: &Uuid) -> Result<Category, Error> {}
//     fn create_category(&self, category: &NewCategory) -> Result<Category, Error> {}
//     fn update_category(
//         &self,
//         category_id: &Uuid,
//         category: &UpdateCategory,
//     ) -> Result<Category, Error> {
//     }
//     fn delete_category(&self, id: &Uuid) -> Result<(), Error> {}
// }
