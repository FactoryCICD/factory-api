use std::net::TcpListener;

use crate::api::{
    github_webhook, routes::create_project, routes::delete_project, routes::get_project,
    routes::health_check, update_project,
};
use actix_web::{dev::Server, web::Data, App, HttpServer};
use sqlx::{postgres::PgPoolOptions, PgPool};
use tracing::info;

use super::configuration::{DatabaseSettings, Settings};

pub struct Application {
    port: u16,
    server: Server,
}

impl Application {
    pub async fn build(configuration: Settings) -> Result<Self, std::io::Error> {
        let connection_pool = get_connection_pool(&configuration.database);

        let address = format!(
            "{}:{}",
            configuration.application.host, configuration.application.port
        );
        let listener = TcpListener::bind(address)?;
        let port = listener.local_addr().unwrap().port();
        let server = run(
            listener,
            connection_pool,
            configuration.application.base_url,
        )?;

        Ok(Self { port, server })
    }

    pub fn port(&self) -> u16 {
        self.port
    }

    pub async fn run_until_stopped(self) -> Result<(), std::io::Error> {
        info!("Server running on port: {}", self.port);
        self.server.await
    }
}

#[derive(Debug)]
pub struct ApplicationBaseUrl(pub String);

#[derive(Debug)]
pub struct ApplicationPort(pub u16);

pub fn run(
    listener: TcpListener,
    connection_pool: PgPool,
    base_url: String,
) -> Result<Server, std::io::Error> {
    let connection = Data::new(connection_pool);
    let base_url = Data::new(ApplicationBaseUrl(base_url));
    let port = Data::new(ApplicationPort(
        listener.local_addr().expect("Cannot Get Port").port(),
    ));
    info!("Creating Server at: {}", port.0);
    let server = HttpServer::new(move || {
        App::new()
            .service(health_check)
            .service(github_webhook)
            .service(get_project)
            .service(create_project)
            .service(delete_project)
            .service(update_project)
            .app_data(connection.clone())
            .app_data(base_url.clone())
            .app_data(port.clone())
    })
    .listen(listener)?
    .run();

    Ok(server)
}

pub fn get_connection_pool(configuration: &DatabaseSettings) -> PgPool {
    PgPoolOptions::new()
        .acquire_timeout(std::time::Duration::from_secs(2))
        .connect_lazy_with(configuration.with_db())
}
