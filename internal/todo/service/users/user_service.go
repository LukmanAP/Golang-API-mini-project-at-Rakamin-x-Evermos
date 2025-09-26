package users

import (
    "errors"
    "net/mail"
    "strings"
    "time"

    repo "project-evermos/internal/todo/repository/users"
    model "project-evermos/internal/todo/model/users"

    "golang.org/x/crypto/bcrypt"
)

var (
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrNotFound     = errors.New("not_found")
    ErrDuplicate    = errors.New("duplicate")
    ErrBadRequest   = errors.New("bad_request")
)

type Service struct{ repo *repo.Repository }

func NewService(r *repo.Repository) *Service { return &Service{repo: r} }

// -------- Profile --------
func (s *Service) GetProfile(userID uint) (*model.User, error) {
    u, err := s.repo.FindByID(userID)
    if err != nil { return nil, ErrNotFound }
    return u, nil
}

type UpdateProfileInput struct {
    Nama          string
    KataSandi     string
    NoTelp        string
    TanggalLahir  string // dd/MM/yyyy
    Pekerjaan     string
    Email         string
    IDProvinsi    string
    IDKota        string
}

func parseDateDMY(s string) (*time.Time, error) {
    s = strings.TrimSpace(s)
    if s == "" { return nil, nil }
    t, err := time.Parse("02/01/2006", s)
    if err != nil { return nil, err }
    return &t, nil
}

func (s *Service) UpdateProfile(userID uint, in UpdateProfileInput) error {
    // validations
    var errs []string
    if strings.TrimSpace(in.Nama) == "" { errs = append(errs, "nama wajib diisi") }
    if strings.TrimSpace(in.Pekerjaan) == "" { errs = append(errs, "pekerjaan wajib diisi") }
    if _, err := mail.ParseAddress(strings.TrimSpace(in.Email)); err != nil { errs = append(errs, "Email tidak valid") }
    // phone regex: digits 10-15
    cleanedPhone := strings.TrimSpace(in.NoTelp)
    if cleanedPhone == "" || len(cleanedPhone) < 10 || len(cleanedPhone) > 15 || strings.ContainsAny(cleanedPhone, " ^+-.()abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
        errs = append(errs, "no_telp harus 10-15 digit")
    }

    if strings.TrimSpace(in.IDProvinsi) == "" || strings.TrimSpace(in.IDKota) == "" { errs = append(errs, "id_provinsi dan id_kota wajib diisi") }
    // Placeholder: simple check kota belongs to prov
    if !strings.HasPrefix(strings.TrimSpace(in.IDKota), strings.TrimSpace(in.IDProvinsi)) && strings.TrimSpace(in.IDKota) != strings.TrimSpace(in.IDProvinsi) {
        errs = append(errs, "id_kota tidak sesuai dengan id_provinsi")
    }

    var tgl *time.Time
    if strings.TrimSpace(in.TanggalLahir) != "" {
        dt, err := parseDateDMY(in.TanggalLahir)
        if err != nil { errs = append(errs, "tanggal_Lahir format harus dd/MM/yyyy") } else { tgl = dt }
    }

    if len(in.KataSandi) > 0 && len(in.KataSandi) < 6 { errs = append(errs, "kata_sandi minimal 6 karakter") }

    if len(errs) > 0 { return errors.New(strings.Join(errs, ", ")) }

    // uniqueness checks
    if u2, err := s.repo.FindByEmail(in.Email); err == nil && u2.ID != userID { return ErrDuplicate }
    if u3, err := s.repo.FindByPhone(in.NoTelp); err == nil && u3.ID != userID { return ErrDuplicate }

    // build update model
    u := &model.User{
        ID:           userID,
        Nama:         strings.TrimSpace(in.Nama),
        NoTelp:       cleanedPhone,
        Pekerjaan:    strings.TrimSpace(in.Pekerjaan),
        Email:        strings.TrimSpace(in.Email),
        IDProvinsi:   strings.TrimSpace(in.IDProvinsi),
        IDKota:       strings.TrimSpace(in.IDKota),
        TanggalLahir: tgl,
    }

    if strings.TrimSpace(in.KataSandi) != "" {
        h, err := bcrypt.GenerateFromPassword([]byte(in.KataSandi), bcrypt.DefaultCost)
        if err != nil { return err }
        u.KataSandi = string(h)
    }

    if err := s.repo.UpdateSelf(u); err != nil { return err }
    return nil
}

// -------- Alamat --------
func (s *Service) ListAlamat(userID uint, title string) ([]model.Alamat, error) {
    return s.repo.ListAlamatByUser(userID, title)
}

func (s *Service) GetAlamat(userID, id uint) (*model.Alamat, error) {
    // Fetch alamat by ID first
    a, err := s.repo.GetAlamatByID(id)
    if err != nil { return nil, ErrNotFound }
    // Ownership enforcement
    if a.IDUser != userID { return nil, ErrForbidden }
    return a, nil
}

type CreateAlamatInput struct {
    JudulAlamat  string
    NamaPenerima string
    NoTelp       string
    DetailAlamat string
}

func (s *Service) CreateAlamat(userID uint, in CreateAlamatInput) (uint, error) {
    var errs []string
    if strings.TrimSpace(in.JudulAlamat) == "" { errs = append(errs, "judul_alamat wajib diisi") }
    if strings.TrimSpace(in.NamaPenerima) == "" { errs = append(errs, "nama_penerima wajib diisi") }
    if strings.TrimSpace(in.NoTelp) == "" { errs = append(errs, "no_telp wajib diisi") }
    if strings.TrimSpace(in.DetailAlamat) == "" { errs = append(errs, "detail_alamat wajib diisi") }
    if len(errs) > 0 { return 0, errors.New(strings.Join(errs, ", ")) }

    a := &model.Alamat{
        IDUser:       userID,
        JudulAlamat:  strings.TrimSpace(in.JudulAlamat),
        NamaPenerima: strings.TrimSpace(in.NamaPenerima),
        NoTelp:       strings.TrimSpace(in.NoTelp),
        DetailAlamat: strings.TrimSpace(in.DetailAlamat),
    }
    if err := s.repo.CreateAlamat(a); err != nil { return 0, err }
    return a.ID, nil
}

type UpdateAlamatInput struct {
    JudulAlamat  string
    NamaPenerima string
    NoTelp       string
    DetailAlamat string
}

func (s *Service) UpdateAlamat(userID, id uint, in UpdateAlamatInput) error {
    // ownership check: fetch by ID regardless of owner
    a0, err := s.repo.GetAlamatByID(id)
    if err != nil { return ErrNotFound }
    if a0.IDUser != userID { return ErrForbidden }

    var errs []string
    if strings.TrimSpace(in.JudulAlamat) == "" { errs = append(errs, "judul_alamat wajib diisi") }
    if strings.TrimSpace(in.NamaPenerima) == "" { errs = append(errs, "nama_penerima wajib diisi") }
    if strings.TrimSpace(in.NoTelp) == "" { errs = append(errs, "no_telp wajib diisi") }
    if strings.TrimSpace(in.DetailAlamat) == "" { errs = append(errs, "detail_alamat wajib diisi") }
    if len(errs) > 0 { return errors.New(strings.Join(errs, ", ")) }

    a := &model.Alamat{ID: id, JudulAlamat: strings.TrimSpace(in.JudulAlamat), NamaPenerima: strings.TrimSpace(in.NamaPenerima), NoTelp: strings.TrimSpace(in.NoTelp), DetailAlamat: strings.TrimSpace(in.DetailAlamat)}
    return s.repo.UpdateAlamatForUser(userID, a)
}

func (s *Service) DeleteAlamat(userID, id uint) error {
    // ownership check: fetch by ID regardless of owner
    a0, err := s.repo.GetAlamatByID(id)
    if err != nil { return ErrNotFound }
    if a0.IDUser != userID { return ErrForbidden }
    return s.repo.DeleteAlamatForUser(userID, id)
}