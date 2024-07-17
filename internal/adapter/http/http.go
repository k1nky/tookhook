package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	DefaultReadTimeout  = 10 * time.Second
	DefaultWriteTimeout = 10 * time.Second
	DefaultCloseTimeout = 5 * time.Second
)

type Adapter struct {
	hooker hookService
	log    logger
}

func New(log logger, hooker hookService) *Adapter {
	a := &Adapter{
		log:    log,
		hooker: hooker,
	}

	return a
}

func (a *Adapter) ListenAndServe(ctx context.Context, addr string) {
	srv := &http.Server{
		Handler:      a.buildRouter(),
		Addr:         addr,
		WriteTimeout: DefaultWriteTimeout,
		ReadTimeout:  DefaultReadTimeout,
	}
	go func() {
		a.log.Infof("listen %s", addr)
		if err := srv.ListenAndServe(); err != nil {
			a.log.Debugf("http server was closed")
			if !errors.Is(err, http.ErrServerClosed) {
				a.log.Errorf("unexpected server closing: %v", err)
			}
		}
	}()
	go func() {
		<-ctx.Done()
		a.log.Debugf("closing http server")
		c, cancel := context.WithTimeout(context.Background(), DefaultCloseTimeout)
		defer cancel()
		srv.Shutdown(c)
	}()
}

func (a *Adapter) buildRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/hook", func(r chi.Router) {
		r.Use(TraceRequest(), LoggingMiddleware(a.log))
		r.Post("/{name}", a.ForwardHook)
	})
	r.Route("/-", func(r chi.Router) {
		r.Use(TraceRequest(), LoggingMiddleware(a.log))
		r.Get("/reload", a.Reload)
	})
	r.Get("/health", a.Health)

	return r
}

func unescapeBody(data []byte) []byte {
	s, err := url.QueryUnescape(string(data))
	if err != nil {
		return data
	}
	return []byte(s)
}

func (a *Adapter) ForwardHook(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.Contains(r.Header.Get("content-type"), "application/x-www-form-urlencoded") {
		data = unescapeBody(data)
	}
	if err = a.hooker.Forward(r.Context(), name, data); err != nil {
		requestId := r.Context().Value(KeyRequestId)
		a.log.Errorf("request id=%d failed %v", requestId, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Adapter) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	pluginstatus := a.hooker.Status(r.Context())
	response := make(map[string]interface{})
	response["status"] = "OK"
	response["hooker"] = pluginstatus
	body, err := json.Marshal(response)
	if err != nil {
		a.log.Errorf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func (a *Adapter) Reload(w http.ResponseWriter, r *http.Request) {
	err := a.hooker.Reload(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
