-- migrate:up
CREATE TABLE "users" (
    "id" uuid DEFAULT uuid_generate_v4(),
    "username" varchar NOT NULL UNIQUE,
    "password" varchar NOT NULL,
    "email" varchar NOT NULL DEFAULT '' UNIQUE,
    "phone_num" varchar NOT NULL DEFAULT '' UNIQUE,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    PRIMARY KEY ("id")
);

CREATE TABLE "products" (
    "id" uuid DEFAULT uuid_generate_v4(),
    "name" varchar NOT NULL,
    "price" integer NOT NULL,
    "description" varchar NOT NULL DEFAULT '',
    "image" varchar NOT NULL DEFAULT '',
    "category" varchar NOT NULL DEFAULT '',
    "stock" integer NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    PRIMARY KEY ("id")
);

CREATE TABLE "user_cart" (
    "id" uuid DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "product_id" uuid NOT NULL,
    "quantity" integer NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "orders" (
    "id" uuid DEFAULT uuid_generate_v4(),
    "status" varchar NOT NULL DEFAULT 'pending',
    "user_id" uuid NOT NULL,
    "product_id" uuid NOT NULL,
    "quantity" integer NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    PRIMARY KEY ("id")
);

ALTER TABLE "user_cart" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_cart" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "orders" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

CREATE INDEX ON "users" ("username");
CREATE INDEX ON "products" ("name");
CREATE INDEX ON "user_cart" ("user_id");
CREATE INDEX ON "user_cart" ("product_id");

-- migrate:down
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS user_cart;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;