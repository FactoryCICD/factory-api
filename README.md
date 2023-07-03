# FactoryCI

Rust Applications
MySQL Database

A workflow can be made up of multiple pipelines
A pipeline can be made up of multiple stages
A stage can perform multiple executions

## API

Webhook Handler

1. Handle webhooks from GitHub
2. POST /hooks/github
   1. Parse webhook payload
   2. Read project information from database
   3. Call the orchestrator to queue a new workflow

Project Manager

1. GET /project?id=123
2. POST /project
3. PUT /project?id=123
4. DELETE /project?id=123

Projects Table

| Field | Type | Description |
| ----- | ---- | ----------- |
| id | uuid | Project ID |
| Owner | string |
| Project | string |
| url | string | GitHub URL |
| webhook_secret | string | GitHub Webhook Secret used to verify x-hub-signature-256 header |

Owner Table

| Field | Type | Description |
| ----- | ---- | ----------- |
| id | uuid | Owner ID |
| name | string | Owner Name |

Orchestrator

1. Queue new workflows (function called by webhook handler)
   1. Use the project information to provision new build agent
   2. Use Consul for service discovery
   3. Use Kubernetes for container (build agent) orchestration

## Build Agent

The build agent is an image that contains the following capabilities:

1. Checkout the source code from GitHub
2. Parse the .factory directory
   1. All *.hcl files are parsed and converted into a single configuration
3. Determine which pipelines are to be executed (checking out the filter stuff)
4. Execute pipelines according to stages
5. Execute the stages

## Flows

### Setting up a new project

Read of the repo of the owner

1. User installs the GitHub app (manually)
1. User creates a new project
   1. If they are not an existing owner, create a new owner
      1. This will authenticate them with GitHub
   2. Sets up the webhook automatically

### Running a workflow

Pre-requisites: The project must exist

1. Once a webhook event is recieved from GitHub, the webhook handler will queue a new workflow
2. The orchestrator will be reading from the queue and provision a new pod to run the workflow
   1. The pod will have an installation of the factory_agent cli tool
3. The agent will read the project information from the database (including auth mechanism)
4. The agent will authenticate with GitHub
5. The agent will read the .factory directory from the repository (prior to checking out everything else)
   1. Syntax check the .hcl files
   2. Conditionally check whether which pipelines are to be executed
6. Execute pipelines according to stages
   1. Checkout code
