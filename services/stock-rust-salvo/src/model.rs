extern crate diesel;

use crate::repository::diesel::schema::categories;
use diesel::prelude::*;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Queryable, Serialize, Deserialize, Clone, Debug)]
pub struct Category {
    pub id: Uuid,
    pub name: String,
}

#[derive(Deserialize, Insertable, Debug)]
#[diesel(table_name = categories)]
pub struct NewCategory {
    pub name: String,
}

#[derive(Deserialize, AsChangeset, Debug)]
#[diesel(table_name = categories)]
pub struct UpdateCategory {
    pub name: String,
}
