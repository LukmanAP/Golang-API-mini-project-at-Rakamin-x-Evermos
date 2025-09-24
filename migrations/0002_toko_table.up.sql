-- 0002_toko_table.up.sql
CREATE TABLE IF NOT EXISTS toko (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_user INT,
  nama_toko VARCHAR(255),
  url_foto VARCHAR(255),
  updated_at DATE,
  created_at DATE,
  CONSTRAINT fk_toko_user 
    FOREIGN KEY (id_user) REFERENCES users(id)
    ON UPDATE CASCADE ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;