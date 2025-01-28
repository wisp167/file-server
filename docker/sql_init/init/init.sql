
CREATE TABLE files
(
	id UUID PRIMARY KEY,
	create_time TIMESTAMP NOT NULL,
	update_time TIMESTAMP NOT NULL,
	file_data BYTEA NOT NULL,
	file_name VARCHAR(255) NOT NULL
);
