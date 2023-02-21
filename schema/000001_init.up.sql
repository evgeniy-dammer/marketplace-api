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

CREATE TABLE users_statuses
(
    id SMALLSERIAL PRIMARY KEY,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE users
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone CHARACTER VARYING (255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    first_name CHARACTER VARYING (255),
    last_name CHARACTER VARYING (255),
    status_id SMALLINT REFERENCES users_statuses(id) DEFAULT 1,
    is_deleted BOOLEAN DEFAULT FALSE,
    user_created UUID REFERENCES users(id),
    user_updated UUID REFERENCES users(id),
    user_deleted UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE roles
(
    id SMALLSERIAL PRIMARY KEY NOT NULL,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE users_roles
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    role_id SMALLINT REFERENCES roles(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE organizations
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL,
    address TEXT,
    phone CHARACTER VARYING (255),
    is_deleted BOOLEAN DEFAULT FALSE,
    user_created UUID REFERENCES users(id),
    user_updated UUID REFERENCES users(id),
    user_deleted UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE images_types
(
    id SMALLSERIAL PRIMARY KEY,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE images
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    object_id UUID NOT NULL,
    type SMALLINT REFERENCES images_types(id) ON DELETE CASCADE NOT NULL,
    origin UUID,
    middle UUID,
    small UUID,
    organization_id UUID REFERENCES organizations(id) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    user_created UUID REFERENCES users(id),
    user_updated UUID REFERENCES users(id),
    user_deleted UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE categories
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name_tm CHARACTER VARYING (255) NOT NULL,
    name_ru CHARACTER VARYING (255) NOT NULL,
    name_tr CHARACTER VARYING (255) NOT NULL,
    name_en CHARACTER VARYING (255) NOT NULL,
    parent_id TEXT,
    level SMALLINT NOT NULL DEFAULT 0,
    organization_id UUID REFERENCES organizations(id) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    user_created UUID REFERENCES users(id),
    user_updated UUID REFERENCES users(id),
    user_deleted UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE items
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name_tm TEXT NOT NULL,
    name_ru TEXT NOT NULL,
    name_tr TEXT NOT NULL,
    name_en TEXT NOT NULL,
    image_id UUID REFERENCES images(id),
    price DOUBLE PRECISION DEFAULT 0,
    category_id UUID NOT NULL,
    organization_id UUID REFERENCES organizations(id) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    user_created UUID REFERENCES users(id),
    user_updated UUID REFERENCES users(id),
    user_deleted UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE categories_items
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE NOT NULL,
    items_id UUID REFERENCES items(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE tables
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    organization_id UUID REFERENCES organizations(id) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    user_created UUID REFERENCES users(id),
    user_updated UUID REFERENCES users(id),
    user_deleted UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE orders_statuses
(
    id SMALLSERIAL PRIMARY KEY,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE orders
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    organization_id UUID REFERENCES organizations(id) NOT NULL,
    table_id UUID REFERENCES tables(id),
    status_id SMALLINT REFERENCES orders_statuses(id) DEFAULT 1,
    totalsum DOUBLE PRECISION DEFAULT 0,
    is_deleted BOOLEAN DEFAULT FALSE,
    user_created UUID REFERENCES users(id),
    user_updated UUID REFERENCES users(id),
    user_deleted UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    deleted_at TIMESTAMPTZ

    CONSTRAINT valid_totalsum CHECK ( totalsum >= 0 )
);

CREATE TABLE orders_items
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) NOT NULL,
    item_id UUID REFERENCES items(id) NOT NULL,
    quantity DOUBLE PRECISION NOT NULL,
    unitprise DOUBLE PRECISION NOT NULL,
    totalprice DOUBLE PRECISION NOT NULL

    CONSTRAINT valid_quantity CHECK ( quantity >= 0 )
    CONSTRAINT valid_unitprise CHECK ( unitprise >= 0 )
    CONSTRAINT valid_totalprice CHECK ( totalprice >= 0 )
);

CREATE TABLE comments_statuses
(
    id SMALLSERIAL PRIMARY KEY,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE comments
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES items(id) NOT NULL,
    organization_id UUID REFERENCES organizations(id) NOT NULL,
    content TEXT NOT NULL,
    status_id SMALLINT REFERENCES comments_statuses(id) DEFAULT 1,
    is_deleted BOOLEAN DEFAULT FALSE,
    user_created UUID REFERENCES users(id),
    user_updated UUID REFERENCES users(id),
    user_deleted UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    deleted_at TIMESTAMPTZ
);

-- DATA --

INSERT INTO users_statuses (name) VALUES ('active'), ('inactive'), ('blocked');

INSERT INTO roles (name) VALUES ('admin'), ('analyst'), ('vendor'), ('operator'), ('customer');

INSERT INTO orders_statuses (name) VALUES ('new'), ('approved'), ('canceled'), ('in process'), ('on the way'), ('shipped'), ('returned'), ('payed');

INSERT INTO images_types (name) VALUES ('user'), ('item');

INSERT INTO comments_statuses (name) VALUES ('new'), ('approved'), ('canceled');

COMMIT;