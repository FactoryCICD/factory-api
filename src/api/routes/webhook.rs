use actix_web::{post, web::Json, HttpRequest, HttpResponse};
use serde::{Deserialize, Serialize};
use serde_json::Value;

use crate::api::error::ApiError;

pub enum GithubEvent {
    /// This event occurs when there is activity relating to which repositories a GitHub App installation can access.
    /// All GitHub Apps receive this event by default. You cannot manually subscribe to this event.
    InstallationRepositories,
    /// This event occurs when there is activity relating to a GitHub App installation.
    /// All GitHub Apps receive this event by default. You cannot manually subscribe to this event
    Installation,
    /// This event occurs when a commit or tag is pushed.
    Push,
    /// This event occurs when there is activity on a pull request.
    PullRequest,
}

impl From<&str> for GithubEvent {
    fn from(event: &str) -> Self {
        match event {
            "installation_repositories" => GithubEvent::InstallationRepositories,
            "installation" => GithubEvent::Installation,
            "push" => GithubEvent::Push,
            "pull_request" => GithubEvent::PullRequest,
            _ => panic!("Invalid event"),
        }
    }
}

#[derive(Deserialize, Serialize)]
pub struct WebhookEvent {
    sender: Sender,
    #[serde(default)]
    repository: Repo,
    /// Present when the event type is push
    after: Option<String>,
    /// Present when the event type is pull_request
    pull_request: Option<PullRequest>,
    /// Present when the event type is installation_repositories
    repositories_added: Option<Vec<RepositoryAccess>>,
    /// Present when the event type is installation_repositories
    repositories_removed: Option<Vec<RepositoryAccess>>,
}
#[derive(Deserialize, Serialize)]

struct Sender {
    id: u64,
    login: String,
}
#[derive(Deserialize, Serialize, Default)]

struct Repo {
    id: u64,
    name: String,
    url: String,
}

#[derive(Deserialize, Serialize)]
struct RepositoryAccess {
    /// ID of the repository
    id: u64,
    /// The name of the repository
    name: String,
    /// Whether the repository is private
    private: bool,
}

#[derive(Deserialize, Serialize)]
struct PullRequest {
    created_at: String,
    id: u64,
    title: String,
    state: PullRequestState,
}

#[derive(Deserialize, Serialize)]
enum PullRequestState {
    #[serde(rename = "open")]
    Open,
    #[serde(rename = "closed")]
    Closed,
}

#[post("/github/webhook")]
pub async fn github_webhook(
    payload: Json<Value>,
    req: HttpRequest,
) -> Result<HttpResponse, ApiError> {
    let event = req
        .headers()
        .get("x-github-event")
        .unwrap()
        .to_str()
        .unwrap();
    println!("[EVENT]: {}", event);
    let payload = serde_json::from_value::<WebhookEvent>(payload.0).unwrap();
    println!("{}", serde_json::to_string_pretty(&payload).unwrap());
    let event: GithubEvent = req
        .headers()
        .get("x-github-event")
        .unwrap()
        .to_str()
        .unwrap()
        .into();

    match event {
        GithubEvent::InstallationRepositories => {
            println!("InstallationRepositories");
        }
        GithubEvent::Installation => {
            println!("Installation");
        }
        GithubEvent::Push => {
            println!("Push");
        }
        GithubEvent::PullRequest => {
            println!("PullRequest");
        }
    }

    Ok(HttpResponse::Ok().finish())
}
