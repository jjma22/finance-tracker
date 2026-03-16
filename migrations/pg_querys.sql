-- Create initial database
CREATE DATABASE expenses
WITH  OWNER = postgres

-- Create expenses table

CREATE TABLE expenses  (
    id INT,
    name VARCHAR(255),
    price VARCHAR(255),
    sku VARCHAR(255),
    dateadded timestamp,
    lastupdate timestamp
)

-- Update id to become primary key

ALTER TABLE expenses ADD primary key (id INT)