-- Categories (idempotent via ON DUPLICATE KEY)
INSERT INTO category (nama_category, created_at, updated_at) VALUES
('Fashion', NOW(), NOW()),
('Elektronik', NOW(), NOW()),
('Kesehatan', NOW(), NOW()),
('Rumah Tangga', NOW(), NOW()),
('Olahraga', NOW(), NOW()),
('Makanan & Minuman', NOW(), NOW()),
('Kosmetik', NOW(), NOW()),
('Bayi & Anak', NOW(), NOW()),
('Perlengkapan Ibadah', NOW(), NOW()),
('Aksesoris', NOW(), NOW()),
('Sepatu', NOW(), NOW()),
('Tas', NOW(), NOW()),
('Peralatan Dapur', NOW(), NOW()),
('Perlengkapan Sekolah', NOW(), NOW()),
('Dekorasi', NOW(), NOW()),
('Elektronik Rumah', NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at=VALUES(updated_at);


