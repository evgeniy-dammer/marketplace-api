-- +goose Up
-- +goose StatementBegin

-- EXTENSIONS --

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- TABLES --

CREATE TABLE IF NOT EXISTS users_statuses
(
    id SMALLSERIAL PRIMARY KEY,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users
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

CREATE TABLE IF NOT EXISTS roles
(
    id SMALLSERIAL PRIMARY KEY NOT NULL,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users_roles
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    role_id SMALLINT REFERENCES roles(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS organizations
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

CREATE TABLE IF NOT EXISTS images_types
(
    id SMALLSERIAL PRIMARY KEY,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS images
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

CREATE TABLE IF NOT EXISTS categories
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

CREATE TABLE IF NOT EXISTS brands
(
    id SERIAL PRIMARY KEY,
    name CHARACTER VARYING (500) NOT NULL
);

CREATE TABLE IF NOT EXISTS items
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

CREATE TABLE IF NOT EXISTS specification
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

CREATE TABLE IF NOT EXISTS categories_items
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE NOT NULL,
    items_id UUID REFERENCES items(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS tables
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

CREATE TABLE IF NOT EXISTS orders_statuses
(
    id SMALLSERIAL PRIMARY KEY,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS orders
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

CREATE TABLE IF NOT EXISTS orders_items
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

CREATE TABLE IF NOT EXISTS comments_statuses
(
    id SMALLSERIAL PRIMARY KEY,
    name CHARACTER VARYING (50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS comments
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

CREATE TABLE IF NOT EXISTS users_favorites
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    item_id UUID REFERENCES items(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS casbin_rule
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ptype TEXT DEFAULT 'p',
    v0 TEXT DEFAULT '',
    v1 TEXT DEFAULT '',
    v2 TEXT DEFAULT '',
    v3 TEXT DEFAULT '',
    v4 TEXT DEFAULT '',
    v5 TEXT DEFAULT ''
);

-- DATA --

INSERT INTO users_statuses (name) VALUES ('active'), ('inactive'), ('blocked');

INSERT INTO roles (name) VALUES ('admin'), ('analyst'), ('vendor'), ('operator'), ('customer');

INSERT INTO orders_statuses (name) VALUES ('new'), ('approved'), ('canceled'), ('in process'), ('on the way'), ('shipped'), ('returned'), ('payed');

INSERT INTO images_types (name) VALUES ('user'), ('item');

INSERT INTO comments_statuses (name) VALUES ('new'), ('approved'), ('canceled');

INSERT INTO brands (name) VALUES ('other');

INSERT INTO casbin_rule (v0, v1, v2, v3)
VALUES
   ('customer', 'users', 'get', 'deny'),
   ('customer', 'user', 'get', 'allow'),
   ('customer', 'user', 'post', 'deny'),
   ('customer', 'user', 'patch', 'deny'),
   ('customer', 'user', 'delete', 'deny'),
   ('customer', 'roles', 'get', 'deny'),
   ('customer', 'organizations', 'get', 'allow'),
   ('customer', 'organization', 'get', 'allow'),
   ('customer', 'organization', 'post', 'deny'),
   ('customer', 'organization', 'patch', 'deny'),
   ('customer', 'organization', 'delete', 'deny'),
   ('customer', 'categories', 'get', 'allow'),
   ('customer', 'category', 'get', 'allow'),
   ('customer', 'category', 'post', 'deny'),
   ('customer', 'category', 'patch', 'deny'),
   ('customer', 'category', 'delete', 'deny'),
   ('customer', 'items', 'get', 'allow'),
   ('customer', 'item', 'get', 'allow'),
   ('customer', 'item', 'post', 'deny'),
   ('customer', 'item', 'patch', 'deny'),
   ('customer', 'item', 'delete', 'deny'),
   ('customer', 'tables', 'get', 'allow'),
   ('customer', 'table', 'get', 'allow'),
   ('customer', 'table', 'post', 'deny'),
   ('customer', 'table', 'patch', 'deny'),
   ('customer', 'table', 'delete', 'deny'),
   ('customer', 'orders', 'get', 'allow'),
   ('customer', 'order', 'get', 'allow'),
   ('customer', 'order', 'post', 'deny'),
   ('customer', 'order', 'patch', 'deny'),
   ('customer', 'order', 'delete', 'deny'),
   ('customer', 'images', 'get', 'deny'),
   ('customer', 'image', 'get', 'deny'),
   ('customer', 'image', 'post', 'deny'),
   ('customer', 'image', 'patch', 'deny'),
   ('customer', 'image', 'delete', 'deny'),
   ('customer', 'comments', 'get', 'allow'),
   ('customer', 'comment', 'get', 'allow'),
   ('customer', 'comment', 'post', 'deny'),
   ('customer', 'comment', 'patch', 'deny'),
   ('customer', 'comment', 'delete', 'deny'),
   ('customer', 'specifications', 'get', 'deny'),
   ('customer', 'specification', 'get', 'deny'),
   ('customer', 'specification', 'post', 'deny'),
   ('customer', 'specification', 'patch', 'deny'),
   ('customer', 'specification', 'delete', 'deny'),
   ('customer', 'favorite', 'post', 'allow'),
   ('customer', 'favorite', 'delete', 'allow'),
   ('customer', 'rules', 'get', 'deny'),
   ('customer', 'rule', 'get', 'deny'),
   ('customer', 'rule', 'post', 'deny'),
   ('customer', 'rule', 'patch', 'deny'),
   ('customer', 'rule', 'delete', 'deny'),
   ('customer', 'metrics', 'get', 'deny'),

   ('operator', 'users', 'get', 'deny'),
   ('operator', 'user', 'get', 'deny'),
   ('operator', 'user', 'post', 'deny'),
   ('operator', 'user', 'patch', 'deny'),
   ('operator', 'user', 'delete', 'deny'),
   ('operator', 'roles', 'get', 'deny'),
   ('operator', 'organizations', 'get', 'deny'),
   ('operator', 'organization', 'get', 'deny'),
   ('operator', 'organization', 'post', 'deny'),
   ('operator', 'organization', 'patch', 'deny'),
   ('operator', 'organization', 'delete', 'deny'),
   ('operator', 'categories', 'get', 'allow'),
   ('operator', 'category', 'get', 'allow'),
   ('operator', 'category', 'post', 'deny'),
   ('operator', 'category', 'patch', 'deny'),
   ('operator', 'category', 'delete', 'deny'),
   ('operator', 'items', 'get', 'allow'),
   ('operator', 'item', 'get', 'allow'),
   ('operator', 'item', 'post', 'allow'),
   ('operator', 'item', 'patch', 'allow'),
   ('operator', 'item', 'delete', 'allow'),
   ('operator', 'tables', 'get', 'allow'),
   ('operator', 'table', 'get', 'allow'),
   ('operator', 'table', 'post', 'allow'),
   ('operator', 'table', 'patch', 'allow'),
   ('operator', 'table', 'delete', 'allow'),
   ('operator', 'orders', 'get', 'allow'),
   ('operator', 'order', 'get', 'allow'),
   ('operator', 'order', 'post', 'allow'),
   ('operator', 'order', 'patch', 'allow'),
   ('operator', 'order', 'delete', 'deny'),
   ('operator', 'images', 'get', 'allow'),
   ('operator', 'image', 'get', 'allow'),
   ('operator', 'image', 'post', 'allow'),
   ('operator', 'image', 'patch', 'allow'),
   ('operator', 'image', 'delete', 'allow'),
   ('operator', 'comments', 'get', 'allow'),
   ('operator', 'comment', 'get', 'allow'),
   ('operator', 'comment', 'post', 'deny'),
   ('operator', 'comment', 'patch', 'deny'),
   ('operator', 'comment', 'delete', 'deny'),
   ('operator', 'specifications', 'get', 'allow'),
   ('operator', 'specification', 'get', 'allow'),
   ('operator', 'specification', 'post', 'allow'),
   ('operator', 'specification', 'patch', 'allow'),
   ('operator', 'specification', 'delete', 'allow'),
   ('operator', 'favorite', 'post', 'deny'),
   ('operator', 'favorite', 'delete', 'deny'),
   ('operator', 'rules', 'get', 'deny'),
   ('operator', 'rule', 'get', 'deny'),
   ('operator', 'rule', 'post', 'deny'),
   ('operator', 'rule', 'patch', 'deny'),
   ('operator', 'rule', 'delete', 'deny'),
   ('operator', 'metrics', 'get', 'deny'),

   ('vendor', 'users', 'get', 'deny'),
   ('vendor', 'user', 'get', 'deny'),
   ('vendor', 'user', 'post', 'deny'),
   ('vendor', 'user', 'patch', 'deny'),
   ('vendor', 'user', 'delete', 'deny'),
   ('vendor', 'roles', 'get', 'deny'),
   ('vendor', 'organizations', 'get', 'deny'),
   ('vendor', 'organization', 'get', 'allow'),
   ('vendor', 'organization', 'post', 'deny'),
   ('vendor', 'organization', 'patch', 'allow'),
   ('vendor', 'organization', 'delete', 'deny'),
   ('vendor', 'categories', 'get', 'allow'),
   ('vendor', 'category', 'get', 'allow'),
   ('vendor', 'category', 'post', 'deny'),
   ('vendor', 'category', 'patch', 'deny'),
   ('vendor', 'category', 'delete', 'deny'),
   ('vendor', 'items', 'get', 'allow'),
   ('vendor', 'item', 'get', 'allow'),
   ('vendor', 'item', 'post', 'allow'),
   ('vendor', 'item', 'patch', 'allow'),
   ('vendor', 'item', 'delete', 'allow'),
   ('vendor', 'tables', 'get', 'allow'),
   ('vendor', 'table', 'get', 'allow'),
   ('vendor', 'table', 'post', 'allow'),
   ('vendor', 'table', 'patch', 'allow'),
   ('vendor', 'table', 'delete', 'allow'),
   ('vendor', 'orders', 'get', 'allow'),
   ('vendor', 'order', 'get', 'allow'),
   ('vendor', 'order', 'post', 'allow'),
   ('vendor', 'order', 'patch', 'allow'),
   ('vendor', 'order', 'delete', 'allow'),
   ('vendor', 'images', 'get', 'allow'),
   ('vendor', 'image', 'get', 'allow'),
   ('vendor', 'image', 'post', 'allow'),
   ('vendor', 'image', 'patch', 'allow'),
   ('vendor', 'image', 'delete', 'allow'),
   ('vendor', 'comments', 'get', 'allow'),
   ('vendor', 'comment', 'get', 'allow'),
   ('vendor', 'comment', 'post', 'deny'),
   ('vendor', 'comment', 'patch', 'deny'),
   ('vendor', 'comment', 'delete', 'deny'),
   ('vendor', 'specifications', 'get', 'allow'),
   ('vendor', 'specification', 'get', 'allow'),
   ('vendor', 'specification', 'post', 'allow'),
   ('vendor', 'specification', 'patch', 'allow'),
   ('vendor', 'specification', 'delete', 'allow'),
   ('vendor', 'favorite', 'post', 'deny'),
   ('vendor', 'favorite', 'delete', 'deny'),
   ('vendor', 'rules', 'get', 'deny'),
   ('vendor', 'rule', 'get', 'deny'),
   ('vendor', 'rule', 'post', 'deny'),
   ('vendor', 'rule', 'patch', 'deny'),
   ('vendor', 'rule', 'delete', 'deny'),
   ('vendor', 'metrics', 'get', 'deny'),

   ('analyst', 'users', 'get', 'allow'),
   ('analyst', 'user', 'get', 'allow'),
   ('analyst', 'user', 'post', 'allow'),
   ('analyst', 'user', 'patch', 'deny'),
   ('analyst', 'user', 'delete', 'deny'),
   ('analyst', 'roles', 'get', 'allow'),
   ('analyst', 'organizations', 'get', 'allow'),
   ('analyst', 'organization', 'get', 'allow'),
   ('analyst', 'organization', 'post', 'allow'),
   ('analyst', 'organization', 'patch', 'deny'),
   ('analyst', 'organization', 'delete', 'deny'),
   ('analyst', 'categories', 'get', 'allow'),
   ('analyst', 'category', 'get', 'allow'),
   ('analyst', 'category', 'post', 'allow'),
   ('analyst', 'category', 'patch', 'allow'),
   ('analyst', 'category', 'delete', 'deny'),
   ('analyst', 'items', 'get', 'allow'),
   ('analyst', 'item', 'get', 'allow'),
   ('analyst', 'item', 'post', 'deny'),
   ('analyst', 'item', 'patch', 'deny'),
   ('analyst', 'item', 'delete', 'deny'),
   ('analyst', 'tables', 'get', 'allow'),
   ('analyst', 'table', 'get', 'allow'),
   ('analyst', 'table', 'post', 'deny'),
   ('analyst', 'table', 'patch', 'deny'),
   ('analyst', 'table', 'delete', 'deny'),
   ('analyst', 'orders', 'get', 'allow'),
   ('analyst', 'order', 'get', 'allow'),
   ('analyst', 'order', 'post', 'allow'),
   ('analyst', 'order', 'patch', 'allow'),
   ('analyst', 'order', 'delete', 'deny'),
   ('analyst', 'images', 'get', 'allow'),
   ('analyst', 'image', 'get', 'allow'),
   ('analyst', 'image', 'post', 'deny'),
   ('analyst', 'image', 'patch', 'deny'),
   ('analyst', 'image', 'delete', 'deny'),
   ('analyst', 'comments', 'get', 'allow'),
   ('analyst', 'comment', 'get', 'allow'),
   ('analyst', 'comment', 'post', 'allow'),
   ('analyst', 'comment', 'patch', 'allow'),
   ('analyst', 'comment', 'delete', 'allow'),
   ('analyst', 'specifications', 'get', 'allow'),
   ('analyst', 'specification', 'get', 'allow'),
   ('analyst', 'specification', 'post', 'deny'),
   ('analyst', 'specification', 'patch', 'deny'),
   ('analyst', 'specification', 'delete', 'deny'),
   ('analyst', 'favorite', 'post', 'deny'),
   ('analyst', 'favorite', 'delete', 'deny'),
   ('analyst', 'rules', 'get', 'deny'),
   ('analyst', 'rule', 'get', 'deny'),
   ('analyst', 'rule', 'post', 'deny'),
   ('analyst', 'rule', 'patch', 'deny'),
   ('analyst', 'rule', 'delete', 'deny'),
   ('analyst', 'metrics', 'get', 'deny'),

   ('admin', 'users', 'get', 'allow'),
   ('admin', 'user', 'get', 'allow'),
   ('admin', 'user', 'post', 'allow'),
   ('admin', 'user', 'patch', 'allow'),
   ('admin', 'user', 'delete', 'allow'),
   ('admin', 'roles', 'get', 'allow'),
   ('admin', 'organizations', 'get', 'allow'),
   ('admin', 'organization', 'get', 'allow'),
   ('admin', 'organization', 'post', 'allow'),
   ('admin', 'organization', 'patch', 'allow'),
   ('admin', 'organization', 'delete', 'allow'),
   ('admin', 'categories', 'get', 'allow'),
   ('admin', 'category', 'get', 'allow'),
   ('admin', 'category', 'post', 'allow'),
   ('admin', 'category', 'patch', 'allow'),
   ('admin', 'category', 'delete', 'allow'),
   ('admin', 'items', 'get', 'allow'),
   ('admin', 'item', 'get', 'allow'),
   ('admin', 'item', 'post', 'allow'),
   ('admin', 'item', 'patch', 'allow'),
   ('admin', 'item', 'delete', 'allow'),
   ('admin', 'tables', 'get', 'allow'),
   ('admin', 'table', 'get', 'allow'),
   ('admin', 'table', 'post', 'allow'),
   ('admin', 'table', 'patch', 'allow'),
   ('admin', 'table', 'delete', 'allow'),
   ('admin', 'orders', 'get', 'allow'),
   ('admin', 'order', 'get', 'allow'),
   ('admin', 'order', 'post', 'allow'),
   ('admin', 'order', 'patch', 'allow'),
   ('admin', 'order', 'delete', 'allow'),
   ('admin', 'images', 'get', 'allow'),
   ('admin', 'image', 'get', 'allow'),
   ('admin', 'image', 'post', 'allow'),
   ('admin', 'image', 'patch', 'allow'),
   ('admin', 'image', 'delete', 'allow'),
   ('admin', 'comments', 'get', 'allow'),
   ('admin', 'comment', 'get', 'allow'),
   ('admin', 'comment', 'post', 'allow'),
   ('admin', 'comment', 'patch', 'allow'),
   ('admin', 'comment', 'delete', 'allow'),
   ('admin', 'specifications', 'get', 'allow'),
   ('admin', 'specification', 'get', 'allow'),
   ('admin', 'specification', 'post', 'allow'),
   ('admin', 'specification', 'patch', 'allow'),
   ('admin', 'specification', 'delete', 'allow'),
   ('admin', 'favorite', 'post', 'allow'),
   ('admin', 'favorite', 'delete', 'allow'),
   ('admin', 'rules', 'get', 'allow'),
   ('admin', 'rule', 'get', 'allow'),
   ('admin', 'rule', 'post', 'allow'),
   ('admin', 'rule', 'patch', 'allow'),
   ('admin', 'rule', 'delete', 'allow'),
   ('admin', 'metrics', 'get', 'allow');

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

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- TABLES --

DROP TABLE IF EXISTS categories_items;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users_favorites;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS comments_statuses;
DROP TABLE IF EXISTS orders_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS orders_statuses;
DROP TABLE IF EXISTS specification;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS tables;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS images_types;
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS brands;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS users_statuses;
DROP TABLE IF EXISTS casbin_rule;

-- +goose StatementEnd