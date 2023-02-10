BEGIN;

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
Set default_with_oids = FALSE;

-- EXTENSIONS --

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- TABLES --

CREATE TABLE users
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone CHARACTER VARYING (255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    first_name CHARACTER VARYING (255),
    last_name CHARACTER VARYING (255),
    status_id UUID NOT NULL,
    image_id UUID,
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE roles
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE statuses
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE users_roles
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE organisations
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    user_id UUID NOT NULL,
    address TEXT,
    phone CHARACTER VARYING (255)
);

CREATE TABLE categories
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name CHARACTER VARYING (255) NOT NULL,
    parent_id TEXT,
    level SMALLINT NOT NULL DEFAULT 0,
    organisation_id UUID NOT NULL
);

CREATE TABLE items
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    price DOUBLE PRECISION DEFAULT 0,
    category_id UUID NOT NULL,
    organisation_id UUID NOT NULL
);

CREATE TABLE categories_items
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE NOT NULL,
    items_id UUID REFERENCES items(id) ON DELETE CASCADE NOT NULL
);

-- DATA --

INSERT INTO statuses (name) VALUES ('active'), ('inactive'), ('blocked');

INSERT INTO roles (name) VALUES ('admin'), ('manager'), ('user');

COMMIT;