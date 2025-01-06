-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE table users (
	id serial PRIMARY KEY,
	username varchar(255) not null,
	first_name varchar(255) not null,
	last_name varchar(255) not null,
	profile_info int UNIQUE REFERENCES profile_info(id),
	email varchar(255) not null,
	validated boolean not null,
	completed boolean not null,
	password varchar(255) not null,
	fame_index float not null,
);

CREATE table profile_info(
	id serial PRIMARY KEY,
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
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
