package http

// --- Swagger Models for accurate schemas and example values ---
// Auth
// AuthLoginRequest represents login payload
// swagger:model
type AuthLoginRequest struct {
    NoTelp    string `json:"no_telp" example:"08123456789"`
    KataSandi string `json:"kata_sandi" example:"password123"`
}

// AuthRegisterRequest represents registration payload
// swagger:model
type AuthRegisterRequest struct {
    Nama         string `json:"nama" example:"John Doe"`
    KataSandi    string `json:"kata_sandi" example:"password123"`
    NoTelp       string `json:"no_telp" example:"08123456789"`
    TanggalLahir string `json:"tanggal_Lahir" example:"01/01/1990"`
    Pekerjaan    string `json:"pekerjaan" example:"Reseller"`
    Email        string `json:"email" example:"john@example.com"`
    IDProvinsi   string `json:"id_provinsi" example:"11"`
    IDKota       string `json:"id_kota" example:"1101"`
}

// ProvinceRef basic province reference
// swagger:model
type ProvinceRef struct {
    ID   string `json:"id" example:"11"`
    Name string `json:"name" example:"Aceh"`
}

// CityRef basic city reference
// swagger:model
type CityRef struct {
    ID         string `json:"id" example:"1101"`
    ProvinceID string `json:"province_id" example:"11"`
    Name       string `json:"name" example:"Kota Banda Aceh"`
}

// AuthLoginData returned in login success
// swagger:model
type AuthLoginData struct {
    Nama          string     `json:"nama" example:"John Doe"`
    NoTelp        string     `json:"no_telp" example:"08123456789"`
    TanggalLahir  string     `json:"tanggal_Lahir" example:"01/01/1990"`
    Tentang       string     `json:"tentang" example:"Saya reseller Evermos"`
    Pekerjaan     string     `json:"pekerjaan" example:"Reseller"`
    Email         string     `json:"email" example:"john@example.com"`
    IDProvinsi    ProvinceRef `json:"id_provinsi"`
    IDKota        CityRef     `json:"id_kota"`
    Token         string     `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// AuthLoginResponse envelope
// swagger:model
type AuthLoginResponse struct {
    Status  bool           `json:"status" example:"true"`
    Message string         `json:"message" example:"Succeed to POST data"`
    Errors  []string       `json:"errors" example:""`
    Data    AuthLoginData  `json:"data"`
}

// AuthRegisterResponse envelope
// swagger:model
type AuthRegisterResponse struct {
    Status  bool     `json:"status" example:"true"`
    Message string   `json:"message" example:"Succeed to POST data"`
    Errors  []string `json:"errors" example:""`
    Data    string   `json:"data" example:"Register Succeed"`
}

// ErrorResponse generic error envelope
// swagger:model
type ErrorResponse struct {
    Status  bool     `json:"status" example:"false"`
    Message string   `json:"message" example:"Failed to POST data"`
    Errors  []string `json:"errors" example:"Invalid JSON"`
    Data    interface{} `json:"data"`
}

// Users
// UserProfileData returned by GET /user
// swagger:model
type UserProfileData struct {
    ID            uint        `json:"id" example:"1"`
    Nama          string      `json:"nama" example:"John Doe"`
    NoTelp        string      `json:"no_telp" example:"08123456789"`
    TanggalLahir  string      `json:"tanggal_Lahir" example:"01/01/1990"`
    Pekerjaan     string      `json:"pekerjaan" example:"Reseller"`
    Email         string      `json:"email" example:"john@example.com"`
    IDProvinsi    ProvinceRef `json:"id_provinsi"`
    IDKota        CityRef     `json:"id_kota"`
}

// UserProfileResponse envelope
// swagger:model
type UserProfileResponse struct {
    Status  bool            `json:"status" example:"true"`
    Message string          `json:"message" example:"Succeed to GET data"`
    Errors  []string        `json:"errors" example:""`
    Data    UserProfileData `json:"data"`
}

// Product models
// swagger:model
type ProductPhoto struct {
    ID        uint   `json:"id" example:"1"`
    ProductID uint   `json:"product_id" example:"10"`
    URL       string `json:"url" example:"https://files.local/uploads/products/1758869029454234400-IMG_2867_11zon.jpg"`
}

// swagger:model
type ProductStore struct {
    ID       uint   `json:"id" example:"5"`
    NamaToko string `json:"nama_toko" example:"Toko Budi"`
    URLFoto  string `json:"url_foto" example:"https://files.local/uploads/stores/toko-1758868233503052000.jpg"`
}

// swagger:model
type ProductCategory struct {
    ID           uint   `json:"id" example:"2"`
    NamaCategory string `json:"nama_category" example:"Fashion"`
}

// swagger:model
type Product struct {
    ID             uint            `json:"id" example:"10"`
    NamaProduk     string          `json:"nama_produk" example:"Kemeja Pria Lengan Panjang"`
    Slug           string          `json:"slug" example:"kemeja-pria-lengan-panjang"`
    HargaReseller  int             `json:"harga_reseller" example:"90000"`
    HargaKonsumen  int             `json:"harga_konsumen" example:"120000"`
    Stok           int             `json:"stok" example:"50"`
    Deskripsi      string          `json:"deskripsi" example:"Bahan katun, nyaman dipakai"`
    Toko           ProductStore    `json:"toko"`
    Category       ProductCategory `json:"category"`
    Photos         []ProductPhoto  `json:"photos"`
}

// swagger:model
type ProductListResponse struct {
    Status  bool      `json:"status" example:"true"`
    Message string    `json:"message" example:"Succeed to GET data"`
    Errors  []string  `json:"errors" example:""`
    Data    []Product `json:"data"`
}

// swagger:model
type ProductDetailResponse struct {
    Status  bool     `json:"status" example:"true"`
    Message string   `json:"message" example:"Succeed to GET data"`
    Errors  []string `json:"errors" example:""`
    Data    Product  `json:"data"`
}

// swagger:model
type APIResponseID struct {
    Status  bool   `json:"status" example:"true"`
    Message string `json:"message" example:"Succeed to POST data"`
    Errors  []string `json:"errors" example:""`
    Data    uint   `json:"data" example:"10"`
}

// swagger:model
type APIResponseString struct {
    Status  bool   `json:"status" example:"true"`
    Message string `json:"message" example:"Succeed to PUT data"`
    Errors  []string `json:"errors" example:""`
    Data    string `json:"data" example:""`
}

// Swagger documentation stubs
// The following dummy functions only serve for Swagger annotations.

// @Summary Health check
// @Description Health check endpoint to verify server status
// @Tags Health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func SwaggerHealthCheck() {}

// @Summary Login user
// @Description Authenticate user dengan nomor telepon dan kata sandi
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body AuthLoginRequest true "Login credentials"
// @Success 200 {object} AuthLoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/login [post]
func SwaggerAuthLogin() {}

// @Summary Register user
// @Description Register akun pengguna baru
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body AuthRegisterRequest true "Registration data"
// @Success 200 {object} AuthRegisterResponse "Registration successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Router /auth/register [post]
func SwaggerAuthRegister() {}

// @Summary Get my store
// @Description Get current user's store information
// @Tags Toko
// @Security BearerAuth
// @Produce json
// @Success 200 {object} APIResponseString "Store information"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /toko/my [get]
func SwaggerTokoGetMy() {}

// @Summary Update store
// @Description Update store information dengan optional photo upload
// @Tags Toko
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id_toko path integer true "Store ID" example(5)
// @Param nama_toko formData string false "Store name" example(Toko Budi)
// @Param photo formData file false "Store photo (jpg, jpeg, png)"
// @Success 200 {object} APIResponseString "Update successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not found"
// @Router /toko/{id_toko} [put]
func SwaggerTokoUpdate() {}

// @Summary Get store by ID
// @Description Get public store information by ID
// @Tags Toko
// @Produce json
// @Param id_toko path integer true "Store ID" example(5)
// @Success 200 {object} APIResponseString "Store information"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 404 {object} ErrorResponse "Not found"
// @Router /toko/{id_toko} [get]
func SwaggerTokoGetByID() {}

// @Summary List stores
// @Description Get list of all stores dengan pagination dan pencarian
// @Tags Toko
// @Produce json
// @Param limit query integer false "Results per page" default(10) example(10)
// @Param page query integer false "Page number" default(1) example(1)
// @Param nama query string false "Filter by store name" example(Budi)
// @Success 200 {object} APIResponseString "List of stores"
// @Router /toko [get]
func SwaggerTokoList() {}

// @Summary Get user profile
// @Description Get current user's profile information
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserProfileResponse "User profile"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /user [get]
func SwaggerUserGetProfile() {}

// @Summary Update user profile
// @Description Update current user's profile information
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body AuthRegisterRequest true "Profile update data"
// @Success 200 {object} APIResponseString "Update successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /user [put]
func SwaggerUserUpdateProfile() {}

// @Summary List user addresses
// @Description Get list of user's delivery addresses
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param judul_alamat query string false "Filter by address title" example(Home)
// @Success 200 {object} APIResponseString "List of addresses"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /user/alamat [get]
func SwaggerUserListAlamat() {}

// @Summary Create address
// @Description Create new delivery address
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body object true "Address data" SchemaExample({"judul_alamat":"Home","nama_penerima":"John Doe","no_telp":"08123456789","detail_alamat":"Jl. Sudirman No. 1"})
// @Success 200 {object} APIResponseID "Address created"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /user/alamat [post]
func SwaggerUserCreateAlamat() {}

// @Summary Get address by ID
// @Description Get specific address by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path integer true "Address ID" example(1)
// @Success 200 {object} APIResponseString "Address details"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not found"
// @Router /user/alamat/{id} [get]
func SwaggerUserGetAlamat() {}

// @Summary Update address
// @Description Update delivery address by ID
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path integer true "Address ID" example(1)
// @Param body body object true "Address update data" SchemaExample({"judul_alamat":"Office","nama_penerima":"John Doe","no_telp":"08123456789","detail_alamat":"Jl. Thamrin No. 2"})
// @Success 200 {object} APIResponseString "Address updated"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not found"
// @Router /user/alamat/{id} [put]
func SwaggerUserUpdateAlamat() {}

// @Summary Delete address
// @Description Delete delivery address by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path integer true "Address ID" example(1)
// @Success 200 {object} APIResponseString "Address deleted"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not found"
// @Router /user/alamat/{id} [delete]
func SwaggerUserDeleteAlamat() {}

// @Summary List products
// @Description Get list of products dengan filtering dan pagination
// @Tags Product
// @Produce json
// @Param limit query integer false "Results per page" default(10) example(10)
// @Param page query integer false "Page number" default(1) example(1)
// @Param nama_produk query string false "Filter by product name" example(Kemeja)
// @Param category_id query integer false "Filter by category ID" example(2)
// @Param toko_id query integer false "Filter by store ID" example(5)
// @Param min_harga query integer false "Minimum price filter" example(50000)
// @Param max_harga query integer false "Maximum price filter" example(150000)
// @Success 200 {object} ProductListResponse "List of products"
// @Router /product [get]
func SwaggerProductList() {}

// @Summary Get product by ID
// @Description Get specific product details by ID
// @Tags Product
// @Produce json
// @Param id path integer true "Product ID" example(10)
// @Success 200 {object} ProductDetailResponse "Product details"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Router /product/{id} [get]
func SwaggerProductGetByID() {}

// @Summary Create product
// @Description Create new product dengan upload foto
// @Tags Product
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param nama_produk formData string true "Product name" example(Kemeja Pria Lengan Panjang)
// @Param category_id formData integer true "Category ID" example(2)
// @Param harga_reseller formData integer true "Reseller price" example(90000)
// @Param harga_konsumen formData integer true "Consumer price" example(120000)
// @Param stok formData integer true "Stock quantity" example(50)
// @Param deskripsi formData string false "Product description" example(Bahan katun, nyaman dipakai)
// @Param photos formData file false "Product photos (multiple files supported)"
// @Success 200 {object} APIResponseID "Product created"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /product [post]
func SwaggerProductCreate() {}

// @Summary Update product
// @Description Update product information dengan optional upload foto
// @Tags Product
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path integer true "Product ID" example(10)
// @Param nama_produk formData string false "Product name" example(Kemeja Pria Premium)
// @Param category_id formData integer false "Category ID" example(3)
// @Param harga_reseller formData integer false "Reseller price" example(95000)
// @Param harga_konsumen formData integer false "Consumer price" example(125000)
// @Param stok formData integer false "Stock quantity" example(60)
// @Param deskripsi formData string false "Product description" example(Bahan katun premium)
// @Param photos formData file false "Product photos (multiple files supported)"
// @Success 200 {object} APIResponseString "Product updated"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not found"
// @Router /product/{id} [put]
func SwaggerProductUpdate() {}

// @Summary Delete product
// @Description Delete product by ID
// @Tags Product
// @Security BearerAuth
// @Produce json
// @Param id path integer true "Product ID" example(10)
// @Success 200 {object} APIResponseString "Product deleted"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not found"
// @Router /product/{id} [delete]
func SwaggerProductDelete() {}

// --- Address Swagger models ---
// swagger:model
type AddressProvincesResponse struct {
    Status  bool          `json:"status" example:"true"`
    Message string        `json:"message" example:"Succeed to get data"`
    Errors  []string      `json:"errors" example:""`
    Data    []ProvinceRef `json:"data"`
}

// swagger:model
type AddressCitiesResponse struct {
    Status  bool     `json:"status" example:"true"`
    Message string   `json:"message" example:"Succeed to get data"`
    Errors  []string `json:"errors" example:""`
    Data    []CityRef `json:"data"`
}

// swagger:model
type AddressProvinceDetailResponse struct {
    Status  bool       `json:"status" example:"true"`
    Message string     `json:"message" example:"Succeed to get data"`
    Errors  []string   `json:"errors" example:""`
    Data    ProvinceRef `json:"data"`
}

// swagger:model
type AddressCityDetailResponse struct {
    Status  bool    `json:"status" example:"true"`
    Message string  `json:"message" example:"Succeed to get data"`
    Errors  []string `json:"errors" example:""`
    Data    CityRef `json:"data"`
}

// --- Category Swagger models ---
// swagger:model
type CategoryItem struct {
    ID           uint   `json:"id" example:"1"`
    NamaCategory string `json:"nama_category" example:"Fashion"`
}

// swagger:model
type CategoryListResponse struct {
    Status  bool           `json:"status" example:"true"`
    Message string         `json:"message" example:"Succeed to GET data"`
    Errors  []string       `json:"errors" example:""`
    Data    []CategoryItem `json:"data"`
}

// swagger:model
type CategoryDetailResponse struct {
    Status  bool         `json:"status" example:"true"`
    Message string       `json:"message" example:"Succeed to GET data"`
    Errors  []string     `json:"errors" example:""`
    Data    CategoryItem `json:"data"`
}

// swagger:model
type CategoryCreateRequest struct {
    NamaCategory string `json:"nama_category" example:"Fashion"`
}

// swagger:model
type CategoryUpdateRequest struct {
    NamaCategory string `json:"nama_category" example:"Electronics - Updated"`
}

// --- Transaction Swagger models ---
// swagger:model
type TrxPhoto struct {
    ID        uint   `json:"id" example:"1"`
    ProductID uint   `json:"product_id" example:"10"`
    URL       string `json:"url" example:"https://files.local/uploads/products/p1.jpg"`
}

// swagger:model
type TrxToko struct {
    ID       uint   `json:"id" example:"5"`
    NamaToko string `json:"nama_toko" example:"Toko Budi"`
    URLFoto  string `json:"url_foto" example:"https://files.local/uploads/stores/toko-1.jpg"`
}

// swagger:model
type TrxCategory struct {
    ID           uint   `json:"id" example:"2"`
    NamaCategory string `json:"nama_category" example:"Fashion"`
}

// swagger:model
type TrxProduct struct {
    ID            uint        `json:"id" example:"10"`
    NamaProduk    string      `json:"nama_produk" example:"Kemeja Pria Lengan Panjang"`
    Slug          string      `json:"slug" example:"kemeja-pria-lengan-panjang"`
    HargaReseller int         `json:"harga_reseler" example:"90000"`
    HargaKonsumen int         `json:"harga_konsumen" example:"120000"`
    Deskripsi     string      `json:"deskripsi" example:"Bahan katun, nyaman dipakai"`
    Toko          TrxToko     `json:"toko"`
    Category      TrxCategory `json:"category"`
    Photos        []TrxPhoto  `json:"photos"`
}

// swagger:model
type TrxAlamatKirim struct {
    ID           uint   `json:"id" example:"3"`
    JudulAlamat  string `json:"judul_alamat" example:"Rumah"`
    NamaPenerima string `json:"nama_penerima" example:"John Doe"`
    NoTelp       string `json:"no_telp" example:"08123456789"`
    DetailAlamat string `json:"detail_alamat" example:"Jl. Sudirman No. 1, Bandung"`
}

// swagger:model
type TrxDetailItem struct {
    Product    TrxProduct `json:"product"`
    Toko       TrxToko    `json:"toko"`
    Kuantitas  int        `json:"kuantitas" example:"2"`
    HargaTotal int        `json:"harga_total" example:"240000"`
}

// swagger:model
type TrxItem struct {
    ID          uint           `json:"id" example:"1"`
    HargaTotal  int            `json:"harga_total" example:"240000"`
    KodeInvoice string         `json:"kode_invoice" example:"INV-1700000000"`
    MethodBayar string         `json:"method_bayar" example:"COD"`
    AlamatKirim TrxAlamatKirim `json:"alamat_kirim"`
    DetailTrx   []TrxDetailItem `json:"detail_trx"`
}

// swagger:model
type TransactionListData struct {
    Data  []TrxItem `json:"data"`
    Page  int       `json:"page" example:"1"`
    Limit int       `json:"limit" example:"10"`
}

// swagger:model
type TransactionListResponse struct {
    Status  bool                `json:"status" example:"true"`
    Message string              `json:"message" example:"Succeed to GET data"`
    Errors  []string            `json:"errors" example:""`
    Data    TransactionListData `json:"data"`
}

// swagger:model
type TransactionDetailResponse struct {
    Status  bool    `json:"status" example:"true"`
    Message string  `json:"message" example:"Succeed to GET data"`
    Errors  []string `json:"errors" example:""`
    Data    TrxItem `json:"data"`
}

// swagger:model
type TransactionCreateItem struct {
    ProductID uint `json:"product_id" example:"10"`
    Kuantitas int  `json:"kuantitas" example:"2"`
}

// swagger:model
type TransactionCreateRequest struct {
    MethodBayar string                  `json:"method_bayar" example:"COD"`
    AlamatKirim uint                    `json:"alamat_kirim" example:"3"`
    DetailTrx   []TransactionCreateItem `json:"detail_trx"`
}

// --- Address (Province/City) with concrete models ---
// @Summary List provinces
// @Description Get list of Indonesian provinces
// @Tags Address
// @Produce json
// @Param search query string false "Search keyword"
// @Param limit query integer false "Results limit (0-100)" maximum(100)
// @Param page query integer false "Page number (>=1)" minimum(1)
// @Success 200 {object} AddressProvincesResponse "List of provinces"
// @Failure 502 {object} ErrorResponse "Bad gateway"
// @Failure 504 {object} ErrorResponse "Gateway timeout"
// @Router /provcity/listprovincies [get]
func SwaggerAddressListProvinces() {}

// @Summary List cities by province
// @Description Get list of cities in a specific province
// @Tags Address
// @Produce json
// @Param prov_id path string true "Province ID"
// @Success 200 {object} AddressCitiesResponse "List of cities"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 502 {object} ErrorResponse "Bad gateway"
// @Router /provcity/listcities/{prov_id} [get]
func SwaggerAddressListCities() {}

// @Summary Get province detail
// @Description Get detailed information about a province
// @Tags Address
// @Produce json
// @Param prov_id path string true "Province ID"
// @Success 200 {object} AddressProvinceDetailResponse "Province details"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 502 {object} ErrorResponse "Bad gateway"
// @Router /provcity/detailprovince/{prov_id} [get]
func SwaggerAddressDetailProvince() {}

// @Summary Get city detail
// @Description Get detailed information about a city
// @Tags Address
// @Produce json
// @Param city_id path string true "City ID"
// @Success 200 {object} AddressCityDetailResponse "City details"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 502 {object} ErrorResponse "Bad gateway"
// @Router /provcity/detailcity/{city_id} [get]
func SwaggerAddressDetailCity() {}

// --- Category with concrete models ---
// @Summary List categories
// @Description Get list of all product categories
// @Tags Category
// @Produce json
// @Success 200 {object} CategoryListResponse "List of categories"
// @Router /category [get]
func SwaggerCategoryList() {}

// @Summary Get category by ID
// @Description Get specific category details by ID
// @Tags Category
// @Produce json
// @Param id path integer true "Category ID" example(1)
// @Success 200 {object} CategoryDetailResponse "Category details"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Router /category/{id} [get]
func SwaggerCategoryGetByID() {}

// @Summary Create category
// @Description Create new product category (Admin only)
// @Tags Category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CategoryCreateRequest true "Category data"
// @Success 200 {object} APIResponseID "Category created"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - Admin only"
// @Router /category [post]
func SwaggerCategoryCreate() {}

// @Summary Update category
// @Description Update category by ID (Admin only)
// @Tags Category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path integer true "Category ID" example(1)
// @Param body body CategoryUpdateRequest true "Category update data"
// @Success 200 {object} APIResponseString "Category updated"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Router /category/{id} [put]
func SwaggerCategoryUpdate() {}

// @Summary Delete category
// @Description Delete category by ID (Admin only)
// @Tags Category
// @Security BearerAuth
// @Produce json
// @Param id path integer true "Category ID" example(1)
// @Success 200 {object} APIResponseString "Category deleted"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Router /category/{id} [delete]
func SwaggerCategoryDelete() {}

// --- Transaction with concrete models ---
// @Summary List transactions
// @Description Get list of user's transactions with pagination
// @Tags Transaction
// @Security BearerAuth
// @Produce json
// @Param limit query integer false "Results per page" default(10) example(10)
// @Param page query integer false "Page number" default(1) example(1)
// @Success 200 {object} TransactionListResponse "List of transactions"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /trx [get]
func SwaggerTransactionList() {}

// @Summary Get transaction by ID
// @Description Get specific transaction details by ID
// @Tags Transaction
// @Security BearerAuth
// @Produce json
// @Param id path integer true "Transaction ID" example(1)
// @Success 200 {object} TransactionDetailResponse "Transaction details"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Transaction not found"
// @Router /trx/{id} [get]
func SwaggerTransactionGetByID() {}

// @Summary Create transaction
// @Description Create new transaction
// @Tags Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body TransactionCreateRequest true "Transaction data"
// @Success 200 {object} APIResponseID "Transaction created"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /trx [post]
func SwaggerTransactionCreate() {}