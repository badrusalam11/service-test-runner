-- +migrate Up
CREATE TABLE IF NOT EXISTS tbl_projects (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO tbl_projects (name, url) VALUES
  ('web1', 'http://localhost:5000'),
  ('mobile1', 'http://localhost:5001');
