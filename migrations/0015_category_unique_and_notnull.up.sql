-- Create categories table following requirements
CREATE TABLE IF NOT EXISTS category (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nama_category VARCHAR(255) NOT NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  UNIQUE KEY uq_category_nama (nama_category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;