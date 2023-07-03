use crate::api::{error::ApiError, projects::Project, UpdateProject};
use anyhow::Context;
use sqlx::{query_builder, Executor, PgPool};
use uuid::Uuid;

#[tracing::instrument(name = "Get Project From Database", skip(conn))]
pub async fn get_project_by_id(id: Uuid, conn: &PgPool) -> Result<Vec<Project>, ApiError> {
    let projects: Vec<Project> = sqlx::query!(
        "SELECT id, url, webhook_secret FROM projects WHERE id = $1",
        id
    )
    .fetch_all(conn)
    .await
    .context("Failed to query projects")
    .map_err(ApiError::Database)?
    .iter()
    .map(|row| {
        Project::new(
            row.id.to_string(),
            row.url.clone(),
            row.webhook_secret.clone(),
        )
    })
    .collect();
    Ok(projects)
}

#[tracing::instrument(name = "Get all projects from Database", skip(conn))]
pub async fn get_all_projects(conn: &PgPool) -> Result<Vec<Project>, ApiError> {
    let projects: Vec<Project> = sqlx::query!("SELECT id, url, webhook_secret FROM projects")
        .fetch_all(conn)
        .await
        .context("Failed to query projects")
        .map_err(ApiError::Database)?
        .iter()
        .map(|row| {
            Project::new(
                row.id.to_string(),
                row.url.clone(),
                row.webhook_secret.clone(),
            )
        })
        .collect();
    Ok(projects)
}

#[tracing::instrument(name = "Insert project into database", skip(conn))]
pub async fn insert_project(conn: &PgPool, project: Project) -> Result<(), ApiError> {
    sqlx::query!(
        "INSERT INTO projects (id, url, webhook_secret) VALUES ($1, $2, $3)",
        Uuid::parse_str(&project.id()).unwrap(),
        project.url(),
        project.webhook_secret()
    )
    .execute(conn)
    .await
    .context("Failed to insert new project")
    .map_err(ApiError::Database)?;
    Ok(())
}

#[tracing::instrument(name = "Delete project by id from database", skip(conn))]
pub async fn delete_project_by_id(id: Uuid, conn: &PgPool) -> Result<(), ApiError> {
    sqlx::query!("DELETE FROM projects WHERE id = $1", id)
        .execute(conn)
        .await
        .context("Failed to delete project")
        .map_err(ApiError::Database)?;
    Ok(())
}

#[tracing::instrument(name = "Update Project in database", skip(conn))]
pub async fn update_project(
    conn: &PgPool,
    id: Uuid,
    updates: UpdateProject,
) -> Result<(), ApiError> {
    let mut query_builder = query_builder::QueryBuilder::new("UPDATE projects SET ");
    if let Some(url) = updates.url {
        query_builder.push(" url = ");
        query_builder.push_bind(url);
    }

    query_builder.push(" WHERE id = ");
    query_builder.push_bind(id);

    conn.execute(query_builder.build())
        .await
        .context("Failed to update project")
        .map_err(ApiError::Database)?;
    Ok(())
}
