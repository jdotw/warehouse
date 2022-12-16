use crate::entity::{category, category::Entity as Category};
use sea_orm::*;
use uuid::Uuid;

pub struct Query;

impl Query {
    pub async fn find_category_by_id(
        db: &DbConn,
        id: Uuid,
    ) -> Result<Option<category::Model>, DbErr> {
        Category::find_by_id(id).one(db).await
    }

    /// If ok, returns (category models, num pages).
    pub async fn find_categories_in_page(
        db: &DbConn,
        page: usize,
        categories_per_page: usize,
    ) -> Result<(Vec<category::Model>, usize), DbErr> {
        let res = Category::find().all(db).await;
        Ok((res.unwrap(), 1))
        // // Setup paginator
        // let paginator = Category::find()
        //     .order_by_asc(category::Column::Id)
        //     .paginate(db, categories_per_page);
        // let num_pages = paginator.num_pages().await?;

        // // Fetch paginated categories
        // paginator.fetch_page(page - 1).await.map(|p| (p, num_pages))
    }
}
