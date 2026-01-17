package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/drummerzzz/goxios/src/cache"
	"go.uber.org/zap"
)

// Config contém as configurações para o fluxo OAuth2 client_credentials.
type Config struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	Scopes       []string

	// ExtraParams injeta parâmetros adicionais no form (ex: audience, resource, etc).
	ExtraParams map[string]string

	// Cache habilita cache externo (ex: Redis) para compartilhar token entre instâncias.
	Cache cache.TokenCache

	// RefreshBefore controla o "leeway" para renovar antes de expirar.
	// Default: 30s.
	RefreshBefore time.Duration

	// Now existe pra testes; se nil usa time.Now.
	Now func() time.Time
}

// TokenSource gerencia a obtenção e renovação de tokens OAuth2.
type TokenSource struct {
	httpClient *http.Client
	logger     *zap.Logger
	cfg        Config

	mu            sync.Mutex
	token         string
	expiresAt     time.Time
	refreshBefore time.Duration
	now           func() time.Time
}

// NewTokenSource cria um novo TokenSource.
func NewTokenSource(httpClient *http.Client, logger *zap.Logger, cfg Config) *TokenSource {
	if cfg.RefreshBefore == 0 {
		cfg.RefreshBefore = 30 * time.Second
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if logger == nil {
		logger = zap.NewNop()
	}
	return &TokenSource{
		httpClient:    httpClient,
		logger:        logger,
		cfg:           cfg,
		refreshBefore: cfg.RefreshBefore,
		now:           cfg.Now,
	}
}

// Apply aplica o token de acesso no header Authorization da request.
func (s *TokenSource) Apply(req *http.Request) error {
	if s == nil || req == nil {
		return nil
	}
	tok, err := s.Token(req.Context())
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	return nil
}

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type oauthCachedToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Token retorna um token válido, buscando do cache ou do endpoint se necessário.
func (s *TokenSource) Token(ctx context.Context) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.token != "" && !s.shouldRefreshLocked() {
		return s.token, nil
	}

	if s.cfg.Cache != nil {
		key := s.cacheKey()
		if cached, err := s.cfg.Cache.Get(ctx, key); err == nil && cached != "" {
			var ct oauthCachedToken
			if json.Unmarshal([]byte(cached), &ct) == nil && ct.AccessToken != "" {
				s.token = ct.AccessToken
				if ct.ExpiresIn > 0 {
					s.expiresAt = s.now().Add(time.Duration(ct.ExpiresIn) * time.Second)
				} else {
					s.expiresAt = s.now().Add(s.refreshBefore)
				}
				return s.token, nil
			}
		}
	}

	tok, expiresIn, err := s.fetchToken(ctx)
	if err != nil {
		s.logger.Debug("goxios oauth: failed to fetch token", zap.Error(err))
		return "", err
	}

	s.token = tok
	var ttl time.Duration
	if expiresIn > 0 {
		ttl = time.Duration(expiresIn) * time.Second
		s.expiresAt = s.now().Add(ttl)
	} else {
		ttl = s.refreshBefore
		s.expiresAt = s.now().Add(ttl)
	}

	if s.cfg.Cache != nil {
		key := s.cacheKey()
		ct := oauthCachedToken{
			AccessToken: tok,
			ExpiresIn:   expiresIn,
		}
		if b, err := json.Marshal(ct); err == nil {
			_ = s.cfg.Cache.Set(ctx, key, string(b), ttl)
		}
	}

	return s.token, nil
}

func (s *TokenSource) shouldRefreshLocked() bool {
	if s.expiresAt.IsZero() {
		return true
	}
	return s.now().Add(s.refreshBefore).After(s.expiresAt)
}

func (s *TokenSource) cacheKey() string {
	sum := sha256.Sum256([]byte(s.cfg.ClientID + ":" + s.cfg.ClientSecret))
	return "goxios:oauth:" + hex.EncodeToString(sum[:])
}

func (s *TokenSource) fetchToken(ctx context.Context) (string, int64, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	if len(s.cfg.Scopes) > 0 {
		form.Set("scope", strings.Join(s.cfg.Scopes, " "))
	}
	for k, v := range s.cfg.ExtraParams {
		form.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.cfg.ClientID, s.cfg.ClientSecret)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return "", 0, errors.New("oauth: token endpoint returned error: " + resp.Status + " body=" + string(b))
	}

	var tr oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", 0, err
	}
	if tr.AccessToken == "" {
		return "", 0, errors.New("oauth: empty access_token")
	}
	return tr.AccessToken, tr.ExpiresIn, nil
}
