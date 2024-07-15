package http

import (
	"context"
	"errors"
	"io"
	"net/http"
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
	r.Use(TraceRequest())
	r.Route("/hook", func(r chi.Router) {
		r.Use(LoggingMiddleware(a.log))
		r.Post("/{name}", a.ForwardHook)
	})
	r.Get("/health", a.Health)

	return r
}

func (a *Adapter) ForwardHook(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK"}`))
}
