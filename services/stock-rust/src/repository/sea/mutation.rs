use crate::repository::sea::entity::{category, category::Entity as Category};
use sea_orm::*;
use uuid::Uuid;

pub struct Mutation;

impl Mutation {
    pub async fn create_category(
        db: &DbConn,
        form_data: category::Model,
    ) -> Result<category::ActiveModel, DbErr> {
        category::ActiveModel {
            name: Set(form_data.name.to_owned()),
            ..Default::default()
        }
        .save(db)
        .await
    }

    pub async fn update_category_by_id(
        db: &DbConn,
        id: Uuid,
        form_data: category::Model,
    ) -> Result<category::Model, DbErr> {
        let category: category::ActiveModel = Category::find_by_id(id)
            .one(db)
            .await?
            .ok_or(DbErr::Custom("Cannot find category.".to_owned()))
            .map(Into::into)?;

        category::ActiveModel {
            id: category.id,
            name: Set(form_data.name.to_owned()),
        }
        .update(db)
        .await
    }

    pub async fn delete_category(db: &DbConn, id: Uuid) -> Result<DeleteResult, DbErr> {
        let category: category::ActiveModel = Category::find_by_id(id)
            .one(db)
            .await?
            .ok_or(DbErr::Custom("Cannot find category.".to_owned()))
            .map(Into::into)?;

        category.delete(db).await
    }

    pub async fn delete_all_categorys(db: &DbConn) -> Result<DeleteResult, DbErr> {
        Category::delete_many().exec(db).await
    }
}
