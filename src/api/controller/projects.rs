use sqlx::PgPool;
use uuid::Uuid;

use crate::api::{
    database,
    error::ApiError,
    projects::{NewProject, Project},
};

pub async fn create_project(project: NewProject, conn: &PgPool) -> Result<(), ApiError> {
    let id = Uuid::new_v4();
    let project = Project::new(id.to_string(), project.url, project.webhook_secret);
    // TODO: Validate The URL
    database::projects::insert_project(conn, project).await?;
    Ok(())
}
