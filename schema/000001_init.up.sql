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
    is_main BOOLEAN DEFAULT FALSE,
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

CREATE TABLE brands
(
    id SERIAL PRIMARY KEY,
    name CHARACTER VARYING (500) NOT NULL
);

CREATE TABLE items
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name_tm TEXT NOT NULL,
    name_ru TEXT NOT NULL,
    name_tr TEXT NOT NULL,
    name_en TEXT NOT NULL,
    description_tm TEXT,
    description_ru TEXT,
    description_tr TEXT,
    description_en TEXT,
    internal_id TEXT,
    price DOUBLE PRECISION DEFAULT 0,
    rating DOUBLE PRECISION DEFAULT 0,
    comments_qty INTEGER DEFAULT 0,
    category_id UUID NOT NULL,
    organization_id UUID REFERENCES organizations(id) NOT NULL,
    brand_id INTEGER REFERENCES brands(id) NOT NULL DEFAULT 1,
    is_deleted BOOLEAN DEFAULT FALSE,
    user_created UUID REFERENCES users(id),
    user_updated UUID REFERENCES users(id),
    user_deleted UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    updated_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'gmt'),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE specification
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES items(id) NOT NULL,
    organization_id UUID REFERENCES organizations(id) NOT NULL,
    name_tm TEXT NOT NULL,
    name_ru TEXT NOT NULL,
    name_tr TEXT NOT NULL,
    name_en TEXT NOT NULL,
    description_tm TEXT,
    description_ru TEXT,
    description_tr TEXT,
    description_en TEXT,
    value TEXT
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
    rating DOUBLE PRECISION DEFAULT 0,
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

INSERT INTO brands (name) VALUES ('other');

-- FUNCTIONS --

CREATE OR REPLACE FUNCTION item_rating_comments_inc() RETURNS TRIGGER AS $$
DECLARE
    qty INTEGER;
    rat DOUBLE PRECISION;
BEGIN
    SELECT comments_qty INTO qty FROM items WHERE id = NEW.item_id;
    SELECT AVG(rating) INTO rat FROM comments WHERE item_id = NEW.item_id AND is_deleted = false;
    UPDATE items SET comments_qty = qty + 1, rating = rat WHERE id = NEW.item_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION item_rating_comments_dec() RETURNS TRIGGER AS $$
DECLARE
    qty INTEGER;
    rat DOUBLE PRECISION;
BEGIN
    SELECT comments_qty INTO qty FROM items WHERE id = NEW.item_id;
    if OLD.is_deleted = false THEN
        SELECT AVG(rating) INTO rat FROM comments WHERE item_id = NEW.item_id AND is_deleted = false;

        if rat IS NULL THEN
            UPDATE items SET comments_qty = qty - 1, rating = 0 WHERE id = NEW.item_id;
        ELSE
            UPDATE items SET comments_qty = qty - 1, rating = rat WHERE id = NEW.item_id;
        END IF;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION item_rating_update() RETURNS TRIGGER AS $$
DECLARE
    rat DOUBLE PRECISION;
BEGIN
    SELECT AVG(rating) INTO rat FROM comments WHERE item_id = NEW.item_id AND is_deleted = false;

    if rat IS NULL THEN
        UPDATE items SET rating = 0 WHERE id = NEW.item_id;
    ELSE
        UPDATE items SET rating = rat WHERE id = NEW.item_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- TRIGGERS --

CREATE OR REPLACE TRIGGER item_rating_comments_inc_t AFTER INSERT ON comments FOR EACH ROW EXECUTE PROCEDURE item_rating_comments_inc();

CREATE OR REPLACE TRIGGER item_rating_comments_dec_t AFTER UPDATE OF is_deleted ON comments FOR EACH ROW EXECUTE PROCEDURE item_rating_comments_dec();

CREATE OR REPLACE TRIGGER item_rating_update_t AFTER UPDATE OF rating ON comments FOR EACH ROW EXECUTE PROCEDURE item_rating_update();

COMMIT;