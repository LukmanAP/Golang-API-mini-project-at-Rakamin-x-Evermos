-- 0018_full_seed.down.sql
-- Remove dummy seed data (best-effort)

DELETE FROM detail_trx WHERE id_trx IN (SELECT id FROM trx WHERE kode_invoice='INV-0001');
DELETE FROM trx WHERE kode_invoice='INV-0001';

DELETE FROM log_produk WHERE slug IN ('kaos-polos-premium','headset-bluetooth');
DELETE FROM foto_produk WHERE id_produk IN (
  SELECT id FROM produk WHERE slug IN ('kaos-polos-premium','headset-bluetooth')
);
DELETE FROM produk WHERE slug IN ('kaos-polos-premium','headset-bluetooth');

DELETE FROM category WHERE nama_category IN ('Fashion','Elektronik','Kesehatan');

DELETE FROM toko WHERE nama_toko IN ('Toko Admin','Toko User');
DELETE FROM alamat WHERE `judul alamat` IN ('Rumah','Kantor');

DELETE FROM users WHERE email IN ('admin@example.com','user@example.com');