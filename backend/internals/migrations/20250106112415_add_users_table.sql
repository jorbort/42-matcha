-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE table profile_info(
	id bigserial PRIMARY KEY,
	gender varchar(50) ,
	sexual_orientation varchar(50) ,
	age int ,
	bio varchar(255) ,
	interests text[],
	location geometry(Point, 4326),
);

CREATE table users (
	id bigserial PRIMARY KEY,
	username  varchar(255) UNIQUE not null,
	first_name varchar(255) not null,
	last_name varchar(255) not null,
	profile_info int UNIQUE REFERENCES profile_info(id) on delete cascade,
	email varchar(255) UNIQUE not null,
	validated boolean not null default false,
	completed boolean not null default false,
	password varchar(255) not null,
	fame_index float not null,
	validation_code bytea not null
);

CREATE table user_images(
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_number INT NOT NULL CHECK (image_number BETWEEN 1 AND 5),
    image_url varchar(255) NOT NULL,
    UNIQUE(profile_id, image_number)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
