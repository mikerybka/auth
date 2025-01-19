package auth

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/mikerybka/twilio"
)

type Proxy struct {
	DB           *DB
	TwilioClient *twilio.Client
	BackendURL   *url.URL
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/auth") {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/auth")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		s := &Server{
			DB:           p.DB,
			TwilioClient: p.TwilioClient,
		}
		s.ServeHTTP(w, r)
		return
	}

	httputil.NewSingleHostReverseProxy(p.BackendURL).ServeHTTP(w, r)
}
