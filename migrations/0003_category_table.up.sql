-- 0003_category_table.up.sql
CREATE TABLE IF NOT EXISTS category (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nama_category VARCHAR(255),
  created_at DATE,
  updated_at DATE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;