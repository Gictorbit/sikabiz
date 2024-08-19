BEGIN;
CREATE TABLE IF NOT EXISTS users
(
    id           uuid         NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name         VARCHAR(500) NOT NULL,
    email        VARCHAR(500) NOT NULL UNIQUE,
    phone_number VARCHAR(250) NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses
(
    id       BIGSERIAL PRIMARY KEY,
    street   VARCHAR(255) NOT NULL,
    city     VARCHAR(255) NOT NULL,
    state    VARCHAR(255) NOT NULL,
    zip_code VARCHAR(255) NOT NULL,
    country  VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_address
(
    user_id    uuid   NOT NULL,
    address_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, address_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (address_id) REFERENCES addresses (id) ON DELETE CASCADE
);

COMMIT;