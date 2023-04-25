CREATE TABLE users (
    id bigint PRIMARY KEY 
);

CREATE TYPE task_status AS ENUM (
    'open',
    'closed',
    'declined'
);

CREATE TABLE tasks (
  url VARCHAR(255) PRIMARY KEY,
  assigned_id bigint REFERENCES users(id),
  status task_status,
  price bigserial
);