use crate::schema::categories;
use diesel::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Queryable, Serialize, Deserialize)]
pub struct Category {
    pub id: uuid::Uuid,
    pub name: String,
}

#[derive(Deserialize, Insertable)]
#[diesel(table_name = categories)]
pub struct NewCategory {
    pub name: String,
}

#[derive(Deserialize, AsChangeset)]
#[diesel(table_name = categories)]
pub struct UpdateCategory {
    pub name: String,
}
