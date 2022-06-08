-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2022-03-10T12:34:03.155Z

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL DEFAULT '',
  "phone" varchar NOT NULL DEFAULT '',
  "password_change_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
