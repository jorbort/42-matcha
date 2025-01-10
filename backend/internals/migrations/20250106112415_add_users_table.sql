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
	profile_picture_one varchar(255) ,
	profile_picture_two varchar(255),
	profile_picture_three varchar(255),
	profile_picture_four varchar(255),
	profile_picture_five varchar(255)
);

CREATE table users (
	id bigserial PRIMARY KEY,
	username  varchar(255) UNIQUE not null,
	first_name varchar(255) not null,
	last_name varchar(255) not null,
	profile_info int UNIQUE REFERENCES profile_info(id),
	email varchar(255) UNIQUE not null,
	validated boolean not null default false,
	completed boolean not null default false,
	password bytea not null,
	fame_index float not null
);
3

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
