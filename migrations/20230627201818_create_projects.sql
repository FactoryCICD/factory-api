-- Add migration script here

-- Create Tables

CREATE TABLE projects
(
    id UUID default gen_random_uuid() primary key,
    url varchar(255) not null,
    webhook_secret varchar(255) not null
);

-- Initialize Data

INSERT INTO projects (url, webhook_secret) VALUES ('https://github.com/Cwagne17/MantaNet.git', 'secret_string');