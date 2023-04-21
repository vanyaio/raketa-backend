CREATE TABLE users (
    id INTEGER PRIMARY KEY 
);

CREATE TYPE task_status AS ENUM (
    'open',
    'closed',
    'declined'
);

CREATE TABLE tasks (
  url VARCHAR(255) PRIMARY KEY,
  assigned_id INTEGER REFERENCES users(id),
  status task_status
);