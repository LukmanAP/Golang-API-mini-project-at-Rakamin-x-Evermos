package transaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	trxmodel "project-evermos/internal/todo/model/transaction"
	trxrepo "project-evermos/internal/todo/repository/transaction"

	"gorm.io/gorm"
)

type Service struct {
	repo *trxrepo.Repository
	db   *gorm.DB
}

func NewService(repo *trxrepo.Repository) *Service {
	return &Service{repo: repo, db: repo.DB}
}

// Response structures matching requested format
type TrxListResponse struct {
	Data  []TrxItem `json:"data"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
}

type TrxItem struct {
	ID          uint            `json:"id"`
	HargaTotal  int             `json:"harga_total"`
	KodeInvoice string          `json:"kode_invoice"`
	MethodBayar string          `json:"method_bayar"`
	AlamatKirim AlamatKirimResp `json:"alamat_kirim"`
	DetailTrx   []DetailTrxResp `json:"detail_trx"`
}

type AlamatKirimResp struct {
	ID           uint   `json:"id"`
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	DetailAlamat string `json:"detail_alamat"`
}

type DetailTrxResp struct {
	Product    ProductResp `json:"product"`
	Toko       TokoResp    `json:"toko"`
	Kuantitas  int         `json:"kuantitas"`
	HargaTotal int         `json:"harga_total"`
}

type ProductResp struct {
	ID            uint         `json:"id"`
	NamaProduk    string       `json:"nama_produk"`
	Slug          string       `json:"slug"`
	HargaReseller int          `json:"harga_reseler"`
	HargaKonsumen int          `json:"harga_konsumen"`
	Deskripsi     string       `json:"deskripsi"`
	Toko          TokoResp     `json:"toko"`
	Category      CategoryResp `json:"category"`
	Photos        []PhotoResp  `json:"photos"`
}

type TokoResp struct {
	ID       uint   `json:"id"`
	NamaToko string `json:"nama_toko"`
	URLFoto  string `json:"url_foto"`
}

type CategoryResp struct {
	ID           uint   `json:"id"`
	NamaCategory string `json:"nama_category"`
}

type PhotoResp struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	URL       string `json:"url"`
}

type CreateRequest struct {
	MethodBayar string          `json:"method_bayar"`
	AlamatKirim uint            `json:"alamat_kirim"`
	DetailTrx   []CreateItemReq `json:"detail_trx"`
}

type CreateItemReq struct {
	ProductID uint `json:"product_id"`
	Kuantitas int  `json:"kuantitas"`
}

// List returns user's transactions with pagination
func (s *Service) List(userID uint, limit, page int) (*TrxListResponse, error) {
	trxs, _, err := s.repo.ListTrxByUser(userID, limit, page)
	if err != nil {
		return nil, err
	}

	items := make([]TrxItem, 0, len(trxs))
	for _, t := range trxs {
		item, err := s.buildTrxItem(&t)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return &TrxListResponse{Data: items, Page: page, Limit: limit}, nil
}

// GetByID returns transaction details if owned by user
func (s *Service) GetByID(trxID, userID uint) (*TrxItem, error) {
	owner, err := s.repo.GetOwnerUserIDOfTrx(trxID)
	if err != nil {
		return nil, err
	}
	if owner != userID {
		return nil, errors.New("forbidden")
	}

	trx, err := s.repo.GetTrxByID(trxID)
	if err != nil {
		return nil, err
	}
	if trx == nil {
		return nil, errors.New("not found")
	}

	return s.buildTrxItem(trx)
}

// Create creates a new transaction with validation and snapshot
func (s *Service) Create(userID uint, req CreateRequest) (uint, error) {
	// Validate alamat ownership
	alamat, err := s.repo.GetAlamatByID(req.AlamatKirim)
	if err != nil {
		return 0, err
	}
	if alamat == nil {
		return 0, errors.New("alamat not found")
	}
	if alamat.IDUser != userID {
		return 0, errors.New("alamat not owned by user")
	}

	// Validate products and calculate total
	var hargaTotal int
	items := make([]trxmodel.DetailTrx, 0, len(req.DetailTrx))
	logs := make([]trxmodel.LogProduk, 0, len(req.DetailTrx))

	for _, item := range req.DetailTrx {
		if item.Kuantitas <= 0 {
			return 0, errors.New("kuantitas must be > 0")
		}

		prod, err1 := s.repo.GetProductByID(item.ProductID)
		if err1 != nil {
			return 0, err
		}
		if prod == nil {
			return 0, errors.New("product not found")
		}

		hargaSatuan, _ := strconv.Atoi(prod.HargaKonsumen)
		hargaItem := hargaSatuan * item.Kuantitas
		hargaTotal += hargaItem

		items = append(items, trxmodel.DetailTrx{
			Kuantitas:  item.Kuantitas,
			HargaTotal: hargaItem,
			IDToko:     prod.IDToko,
		})

		// Create product snapshot
		photoURLs := make([]string, len(prod.Photos))
		for i, p := range prod.Photos {
			photoURLs[i] = p.URL
		}

		logs = append(logs, trxmodel.LogProduk{
			IDProduk:      prod.ID,
			NamaProduk:    prod.NamaProduk,
			Slug:          prod.Slug,
			HargaReseller: prod.HargaReseller,
			HargaKonsumen: prod.HargaKonsumen,
			Deskripsi:     prod.Deskripsi,
			IDToko:        prod.IDToko,
			IDCategory:    prod.IDCategory,
			PhotosJSON:    trxrepo.MarshalPhotos(photoURLs),
		})
	}

	// Generate invoice code
	kodeInvoice := fmt.Sprintf("INV-%d", time.Now().Unix())

	// Create transaction within DB transaction
	var trxID uint
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create main transaction
		trx := &trxmodel.Trx{
			IDUser:           userID,
			AlamatPengiriman: req.AlamatKirim,
			HargaTotal:       hargaTotal,
			KodeInvoice:      kodeInvoice,
			MethodBayar:      req.MethodBayar,
		}
		if err1 := s.repo.CreateTrx(tx, trx); err1 != nil {
			return err
		}
		trxID = trx.ID

		// Create product logs first
		for i := range logs {
			if err2 := s.repo.CreateLogProduk(tx, &logs[i]); err2 != nil {
				return err
			}
			items[i].IDLogProduk = logs[i].ID
			items[i].IDTrx = trxID
		}

		// Create detail items
		if err3 := s.repo.CreateDetailItems(tx, items); err3 != nil {
			return err
		}

		// Optional: Reduce stock
		for _, item := range req.DetailTrx {
			if err4 := s.repo.UpdateProductStock(tx, item.ProductID, item.Kuantitas); err4 != nil {
				// Log warning but don't fail transaction for stock issues
				continue
			}
		}

		return nil
	})

	return trxID, err
}

// buildTrxItem constructs response with joined data
func (s *Service) buildTrxItem(trx *trxmodel.Trx) (*TrxItem, error) {
	// Get alamat
	alamat, err := s.repo.GetAlamatByID(trx.AlamatPengiriman)
	if err != nil {
		return nil, err
	}

	alamatResp := AlamatKirimResp{
		ID:           alamat.ID,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat,
	}

	// Get detail items
	details, err := s.repo.GetDetailItems(trx.ID)
	if err != nil {
		return nil, err
	}

	detailResp := make([]DetailTrxResp, 0, len(details))
	for _, detail := range details {
		// Get product log
		log, err := s.repo.GetLogProdukByID(detail.IDLogProduk)
		if err != nil {
			return nil, err
		}

		// Get toko
		toko, err := s.repo.GetTokoByID(log.IDToko)
		if err != nil {
			return nil, err
		}

		// Get category
		cat, err := s.repo.GetCategoryByID(log.IDCategory)
		if err != nil {
			return nil, err
		}

		hargaRes, _ := strconv.Atoi(log.HargaReseller)
		hargaKons, _ := strconv.Atoi(log.HargaKonsumen)

		prodResp := ProductResp{
			ID:            log.IDProduk,
			NamaProduk:    log.NamaProduk,
			Slug:          log.Slug,
			HargaReseller: hargaRes,
			HargaKonsumen: hargaKons,
			Deskripsi:     log.Deskripsi,
			Toko:          TokoResp{ID: toko.ID, NamaToko: toko.NamaToko, URLFoto: toko.UrlFoto},
			Category:      CategoryResp{ID: cat.ID, NamaCategory: cat.NamaCategory},
			Photos:        []PhotoResp{},
		}
		// Parse photos JSON into []PhotoResp
		var urls []string
		_ = json.Unmarshal([]byte(log.PhotosJSON), &urls)
		photos := make([]PhotoResp, 0, len(urls))
		for _, u := range urls {
			photos = append(photos, PhotoResp{ID: 0, ProductID: log.IDProduk, URL: u})
		}
		prodResp.Photos = photos

		detailResp = append(detailResp, DetailTrxResp{
			Product:    prodResp,
			Toko:       TokoResp{ID: toko.ID, NamaToko: toko.NamaToko, URLFoto: toko.UrlFoto},
			Kuantitas:  detail.Kuantitas,
			HargaTotal: detail.HargaTotal,
		})
	}

	return &TrxItem{
		ID:          trx.ID,
		HargaTotal:  trx.HargaTotal,
		KodeInvoice: trx.KodeInvoice,
		MethodBayar: trx.MethodBayar,
		AlamatKirim: alamatResp,
		DetailTrx:   detailResp,
	}, nil
}
