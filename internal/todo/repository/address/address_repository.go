package address

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	addrmodel "project-evermos/internal/todo/model/address"
)

var (
	ErrNotFound = errors.New("not_found")
	ErrUpstream = errors.New("upstream_error")
	ErrTimeout  = errors.New("timeout")
)

type Repository struct {
	baseURL string
	client  *http.Client
	retry   int
}

func NewRepository(baseURL string, timeoutMS int, retry int) *Repository {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://www.emsifa.com/api-wilayah-indonesia/api"
	}
	if timeoutMS <= 0 {
		timeoutMS = 5000
	}
	if retry < 0 {
		retry = 0
	}
	return &Repository{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{Timeout: time.Duration(timeoutMS) * time.Millisecond},
		retry:  retry,
	}
}

func (r *Repository) getJSON(ctx context.Context, path string, out interface{}) error {
	url := r.baseURL + path
	var lastErr error
	attempts := r.retry + 1
	for i := 0; i < attempts; i++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err := r.client.Do(req)
		if err != nil {
			lastErr = classifyHTTPError(err)
			if errors.Is(lastErr, ErrTimeout) && i < attempts-1 { continue }
			return lastErr
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			b, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
			if i < attempts-1 { continue }
			return ErrUpstream
		}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(out); err != nil {
			return ErrUpstream
		}
		return nil
	}
	if lastErr == nil { lastErr = ErrUpstream }
	return lastErr
}

func classifyHTTPError(err error) error {
	if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
		return ErrTimeout
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return ErrUpstream
}

// Public methods
func (r *Repository) ListProvinces(ctx context.Context) ([]addrmodel.Province, error) {
	var data []addrmodel.Province
	if err := r.getJSON(ctx, "/provinces.json", &data); err != nil { return nil, err }
	return data, nil
}

func (r *Repository) GetProvince(ctx context.Context, id string) (*addrmodel.Province, error) {
	var data addrmodel.Province
	if err := r.getJSON(ctx, "/province/"+id+".json", &data); err != nil { return nil, err }
	return &data, nil
}

func (r *Repository) ListRegencies(ctx context.Context, provID string) ([]addrmodel.Regency, error) {
	var data []addrmodel.Regency
	if err := r.getJSON(ctx, "/regencies/"+provID+".json", &data); err != nil { return nil, err }
	return data, nil
}

func (r *Repository) GetRegency(ctx context.Context, id string) (*addrmodel.Regency, error) {
	var data addrmodel.Regency
	if err := r.getJSON(ctx, "/regency/"+id+".json", &data); err != nil { return nil, err }
	return &data, nil
}