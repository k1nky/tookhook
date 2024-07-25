package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
	hs      hookService
	monitor monitorService
	log     logger
	rs      rulesService
}

func New(log logger, hooker hookService, monitor monitorService, rs rulesService) *Adapter {
	a := &Adapter{
		log:     log,
		hs:      hooker,
		monitor: monitor,
		rs:      rs,
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

func (a *Adapter) ForwardHook(w http.ResponseWriter, r *http.Request) {
	var (
		data []byte
		err  error
	)
	requestId := r.Context().Value(KeyRequestId)
	name := chi.URLParam(r, "name")

	badRequest := func(error) {
		a.log.Errorf("request id=%d has bad data: %v", requestId, err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if strings.Contains(r.Header.Get("content-type"), "application/x-www-form") {
		if err := r.ParseForm(); err != nil {
			badRequest(err)
			return
		}
		if data, err = json.Marshal(r.Form); err != nil {
			badRequest(err)
			return
		}
	} else {
		data, err = io.ReadAll(r.Body)
		if err != nil {
			badRequest(err)
			return
		}
	}
	if err = a.hs.Forward(r.Context(), name, data); err != nil {
		a.log.Errorf("request id=%d failed %v", requestId, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Adapter) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	status := a.monitor.Status(r.Context())
	body, err := json.Marshal(status)
	if err != nil {
		a.log.Errorf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func (a *Adapter) Reload(w http.ResponseWriter, r *http.Request) {
	err := a.rs.Load(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
