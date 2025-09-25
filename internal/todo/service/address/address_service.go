package address

import (
	"context"
	"strings"
	"sync"
	"time"

	addrmodel "project-evermos/internal/todo/model/address"
	repo "project-evermos/internal/todo/repository/address"
)

// Re-export repository error sentinels
var (
	ErrNotFound = repo.ErrNotFound
	ErrUpstream = repo.ErrUpstream
	ErrTimeout  = repo.ErrTimeout
)

// Repo abstracts EMSIFA repository for easier testing.
type Repo interface {
	ListProvinces(ctx context.Context) ([]addrmodel.Province, error)
	GetProvince(ctx context.Context, id string) (*addrmodel.Province, error)
	ListRegencies(ctx context.Context, provID string) ([]addrmodel.Regency, error)
	GetRegency(ctx context.Context, id string) (*addrmodel.Regency, error)
}

type Service struct {
	repo Repo
	ttl  time.Duration
	mc   *memCache
}

func NewService(r Repo, ttl time.Duration) *Service {
	var c *memCache
	if ttl > 0 {
		c = newMemCache()
	}
	return &Service{repo: r, ttl: ttl, mc: c}
}

// --------- Public API used by handlers ---------

// ListProvinces supports optional search (case-insensitive substring), limit (1..100, 0 means all), and page (1-based).
func (s *Service) ListProvinces(ctx context.Context, search string, limit, page int) ([]addrmodel.Province, error) {
	// Cache the full provinces list, then filter/paginate per request
	const keyAll = "provinces_all"
	var items []addrmodel.Province
	if s.mc != nil {
		if v, ok := s.mc.get(keyAll); ok {
			items = v.([]addrmodel.Province)
		} else {
			list, err := s.repo.ListProvinces(ctx)
			if err != nil { return nil, err }
			items = list
			s.mc.set(keyAll, items, s.ttl)
		}
	} else {
		list, err := s.repo.ListProvinces(ctx)
		if err != nil { return nil, err }
		items = list
	}

	// search
	t := strings.TrimSpace(search)
	if t != "" {
		low := strings.ToLower(t)
		filtered := make([]addrmodel.Province, 0, len(items))
		for _, p := range items {
			if strings.Contains(strings.ToLower(p.Name), low) {
				filtered = append(filtered, p)
			}
		}
		items = filtered
	}
	// limit/page
	if limit <= 0 { return items, nil }
	if limit > 100 { limit = 100 }
	if page <= 0 { page = 1 }
	offset := (page - 1) * limit
	if offset >= len(items) { return []addrmodel.Province{}, nil }
	end := offset + limit
	if end > len(items) { end = len(items) }
	return items[offset:end], nil
}

func (s *Service) ListCities(ctx context.Context, provID string) ([]addrmodel.Regency, error) {
	key := "cities_" + provID
	var items []addrmodel.Regency
	if s.mc != nil {
		if v, ok := s.mc.get(key); ok {
			items = v.([]addrmodel.Regency)
		} else {
			list, err := s.repo.ListRegencies(ctx, provID)
			if err != nil { return nil, err }
			items = list
			s.mc.set(key, items, s.ttl)
		}
	} else {
		list, err := s.repo.ListRegencies(ctx, provID)
		if err != nil { return nil, err }
		items = list
	}
	return items, nil
}

func (s *Service) DetailProvince(ctx context.Context, id string) (*addrmodel.Province, error) {
	key := "province_" + id
	if s.mc != nil {
		if v, ok := s.mc.get(key); ok {
			p := v.(addrmodel.Province)
			return &p, nil
		}
	}
	p, err := s.repo.GetProvince(ctx, id)
	if err != nil { return nil, err }
	if s.mc != nil { s.mc.set(key, *p, s.ttl) }
	return p, nil
}

func (s *Service) DetailCity(ctx context.Context, id string) (*addrmodel.Regency, error) {
	key := "city_" + id
	if s.mc != nil {
		if v, ok := s.mc.get(key); ok {
			r := v.(addrmodel.Regency)
			return &r, nil
		}
	}
	r, err := s.repo.GetRegency(ctx, id)
	if err != nil { return nil, err }
	if s.mc != nil { s.mc.set(key, *r, s.ttl) }
	return r, nil
}

// --------- tiny in-memory cache ---------

type cacheEntry struct {
	val      interface{}
	expireAt time.Time
}

type memCache struct {
	mu sync.RWMutex
	m  map[string]cacheEntry
}

func newMemCache() *memCache { return &memCache{m: make(map[string]cacheEntry)} }

func (c *memCache) get(key string) (interface{}, bool) {
	c.mu.RLock()
	e, ok := c.m[key]
	c.mu.RUnlock()
	if !ok { return nil, false }
	if !e.expireAt.IsZero() && time.Now().After(e.expireAt) {
		c.mu.Lock()
		delete(c.m, key)
		c.mu.Unlock()
		return nil, false
	}
	return e.val, true
}

func (c *memCache) set(key string, v interface{}, ttl time.Duration) {
	var exp time.Time
	if ttl > 0 { exp = time.Now().Add(ttl) }
	c.mu.Lock()
	c.m[key] = cacheEntry{val: v, expireAt: exp}
	c.mu.Unlock()
}