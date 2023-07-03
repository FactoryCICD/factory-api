use std::str::FromStr;

use crate::api::controller;
use crate::api::projects::NewProject;
use crate::api::{database, error::ApiError};
use actix_web::patch;
use actix_web::web::Json;
use actix_web::{
    delete, get, post,
    web::{self, Data},
    HttpResponse,
};
use anyhow::{anyhow, Context};
use serde::Deserialize;
use sqlx::PgPool;

#[derive(Deserialize)]
struct GetProjectsQuery {
    id: Option<String>,
}

#[get("/projects")]
async fn get_project(
    query: web::Query<GetProjectsQuery>,
    conn: Data<PgPool>,
) -> Result<HttpResponse, ApiError> {
    let projects;

    if let Some(id) = &query.id {
        // Unwrap the Option<String> to a String
        projects = database::projects::get_project_by_id(uuid::Uuid::from_str(&id).unwrap(), &conn)
            .await?;
    } else {
        projects = database::projects::get_all_projects(&conn).await?;
    }

    Ok(HttpResponse::Ok().json(projects))
}

#[post("/projects")]
async fn create_project(
    project: Json<NewProject>,
    conn: Data<PgPool>,
) -> Result<HttpResponse, ApiError> {
    controller::projects::create_project(project.0, &conn).await?;
    Ok(HttpResponse::Created().finish())
}

#[delete("/projects/{id}")]
async fn delete_project(
    path: web::Path<String>,
    conn: Data<PgPool>,
) -> Result<HttpResponse, ApiError> {
    database::projects::delete_project_by_id(
        uuid::Uuid::from_str(&path.into_inner()).unwrap(),
        &conn,
    )
    .await?;
    Ok(HttpResponse::NoContent().finish())
}

#[derive(serde::Deserialize, Debug)]
pub struct UpdateProject {
    pub url: Option<String>,
}

#[patch("/projects/{id}")]
async fn update_project(
    id: web::Path<String>,
    updated_project: Json<UpdateProject>,
    conn: Data<PgPool>,
) -> Result<HttpResponse, ApiError> {
    // Check if project exists
    let project_id = uuid::Uuid::from_str(&id.into_inner())
        .context("Failed to convert id to UUID")
        .map_err(ApiError::BadRequest)?;
    let project = database::projects::get_project_by_id(project_id, &conn).await?;
    if project.len() == 0 {
        return Err(ApiError::BadRequest(anyhow!("Project does not exist")));
    }

    database::projects::update_project(&conn, project_id, updated_project.0).await?;

    Ok(HttpResponse::Ok().finish())
}
