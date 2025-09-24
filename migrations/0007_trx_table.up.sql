-- 0007_trx_table.up.sql
CREATE TABLE IF NOT EXISTS trx (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_user INT,
  alamat_pengiriman INT,
  harga_total INT,
  kode_invoice VARCHAR(255),
  method_bayar VARCHAR(255),
  updated_at DATE,
  created_at DATE,
  CONSTRAINT fk_trx_user
    FOREIGN KEY (id_user) REFERENCES users(id)
    ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT fk_trx_alamat
    FOREIGN KEY (alamat_pengiriman) REFERENCES alamat(id)
    ON UPDATE CASCADE ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;