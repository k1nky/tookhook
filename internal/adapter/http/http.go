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
	"github.com/k1nky/tookhook/internal/entity/hooks"
)

const (
	DefaultReadTimeout  = 10 * time.Second
	DefaultWriteTimeout = 10 * time.Second
	DefaultCloseTimeout = 5 * time.Second
)

type Adapter struct {
	hs  hookService
	ms  monitorService
	log logger
	rs  rulesService
}

func New(log logger, hooker hookService, monitor monitorService, rs rulesService) *Adapter {
	a := &Adapter{
		log: log,
		hs:  hooker,
		ms:  monitor,
		rs:  rs,
	}

	return a
}

func readFormToJSON(r *http.Request) (data []byte, err error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	if data, err = json.Marshal(r.Form); err != nil {
		return nil, err
	}
	return
}

func newIncomeRequest(r *http.Request) (*hooks.HookRequest, error) {
	var (
		err error
	)
	ir := &hooks.HookRequest{
		Meta: hooks.HookRequestMeta{
			Name: chi.URLParam(r, "name"),
		},
		Content: hooks.HookRequestBody{
			Type: r.Header.Get("content-type"),
		},
	}
	if strings.Contains(ir.Content.Type, "application/x-www-form-urlencoded") {
		ir.Content.Body, err = readFormToJSON(r)
	} else {
		ir.Content.Body, err = io.ReadAll(r.Body)
	}
	requestId := r.Context().Value(KeyRequestId)
	if requestId == nil {
		ir.Meta.ID = 0
	} else {
		ir.Meta.ID = requestId.(uint64)
	}
	return ir, err
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
	ir, err := newIncomeRequest(r)
	if err != nil {
		a.log.Errorf("request id=%d has bad data: %v", ir.Meta.ID, err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if err = a.hs.Forward(r.Context(), ir); err != nil {
		a.log.Errorf("request id=%d failed %v", ir.Meta.ID, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Adapter) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	status := a.ms.Status(r.Context())
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
