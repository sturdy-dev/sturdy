package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"mash/pkg/jwt"
	service_jwt "mash/pkg/jwt/service"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ErrUnauthenticated = fmt.Errorf("unauthenticated")
	ErrForbidden       = fmt.Errorf("forbidden")
)

var (
	authExpirationSecondsMetric = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "sturdy_auth_jwt_expires_in_seconds",
		Buckets: []float64{
			0,
			float64(time.Hour * 24 / time.Second), // 1 day
			float64(time.Hour * 24 * 7 / time.Second), // 1 week
			float64(time.Hour * 24 * 14 / time.Second),
			float64(time.Hour * 24 * 20 / time.Second),
			float64(time.Hour * 24 * 25 / time.Second), // Most JWTs should expire in 25 to 30 days, using more detailed buckets in this range
			float64(time.Hour * 24 * 26 / time.Second),
			float64(time.Hour * 24 * 27 / time.Second),
			float64(time.Hour * 24 * 28 / time.Second),
			float64(time.Hour * 24 * 29 / time.Second),
			float64(time.Hour * 24 * 30 / time.Second),
			float64(time.Hour * 24 * 31 / time.Second),
		},
	})
)

func SubjectFromRequest(r *http.Request, jwtService *service_jwt.Service) (*Subject, error) {
	jwt, _, err := jwtFromRequest(r, jwtService)
	if err != nil {
		return nil, err
	}
	return subjectFromToken(jwt), nil
}

func jwtFromRequest(r *http.Request, jwtService *service_jwt.Service) (*jwt.Token, bool, error) {
	token, fromHeader := tokenFromHeaders(r.Header)
	var fromCookies bool
	if !fromHeader {
		token, fromCookies = tokenFromCookies(r.Cookies())
	}

	jwtToken, err := jwtService.Verify(r.Context(), token, jwt.TokenTypeAuth, jwt.TokenTypeCI)
	if errors.Is(err, service_jwt.ErrInvalidToken) || errors.Is(err, service_jwt.ErrTokenExpired) {
		return nil, false, ErrUnauthenticated
	} else if err != nil {
		return nil, false, fmt.Errorf("failed to verify token: %w", err)
	}
	expiresIn := time.Until(jwtToken.ExpiresAt)

	authExpirationSecondsMetric.Observe(float64(expiresIn))

	shouldRefresh := jwtToken.Type == jwt.TokenTypeAuth && fromCookies && expiresIn < refreshThreshold
	return jwtToken, shouldRefresh, nil
}

var (
	oneDay           = time.Hour * 24
	oneMonth         = 30 * oneDay
	refreshThreshold = 28 * oneDay
)

const (
	authHeaderName = "Authorization"
	tokenPrefix    = "bearer "
)

func tokenFromHeaders(headers http.Header) (string, bool) {
	authHeader := headers.Get(authHeaderName)
	if len(authHeader) < len(tokenPrefix) {
		return "", false
	}
	prefix := authHeader[:len(tokenPrefix)]
	if !strings.EqualFold(prefix, tokenPrefix) {
		return "", false
	}
	return strings.TrimPrefix(authHeader, tokenPrefix), true
}

const (
	authCookieName = "auth"
)

func tokenFromCookies(cookies []*http.Cookie) (string, bool) {
	for _, cookie := range cookies {
		if cookie.Name == authCookieName {
			return cookie.Value, true
		}
	}
	return "", false
}

func setAuthCookie(w http.ResponseWriter, secure bool, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    url.QueryEscape(value),
		MaxAge:   int(oneMonth), // Expire in 30 days
		Path:     "/",
		Domain:   "",
		SameSite: http.SameSiteStrictMode, // Don't send in external frames etc
		Secure:   secure,                  // HTTPS and localhost only
		HttpOnly: true,                    // Don't expose to JS
	})
}

func RemoveAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		MaxAge:   -1, // Expire!
		Path:     "/",
		Domain:   "",
		SameSite: http.SameSiteStrictMode, // Don't send in external frames etc
		HttpOnly: true,                    // Don't expose to JS
	})
}
