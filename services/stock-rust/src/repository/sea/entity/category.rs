use crate::model::Category;
use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Clone, Debug, PartialEq, Eq, DeriveEntityModel, Deserialize, Serialize)]
#[sea_orm(table_name = "categories")]
// #[derive(Queryable, Serialize, Deserialize, Clone, Debug)]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: Uuid,
    pub name: String,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl Model {
    pub fn to_category(&self) -> Category {
        Category {
            id: self.id.clone(),
            name: self.name.clone(),
        }
    }
}

impl ActiveModelBehavior for ActiveModel {}
