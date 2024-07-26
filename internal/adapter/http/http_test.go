package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k1nky/tookhook/internal/adapter/http/mock"
	"github.com/k1nky/tookhook/internal/entity"
	log "github.com/k1nky/tookhook/pkg/logger"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type httpAdapterTestSuite struct {
	suite.Suite
	hs *mock.MockhookService
	ms *mock.MockmonitorService
	rs *mock.MockrulesService
}

func TestHTTPAdapter(t *testing.T) {
	suite.Run(t, new(httpAdapterTestSuite))
}

func (suite *httpAdapterTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.hs = mock.NewMockhookService(ctrl)
	suite.ms = mock.NewMockmonitorService(ctrl)
	suite.rs = mock.NewMockrulesService(ctrl)
}

func (suite *httpAdapterTestSuite) TestHealth() {
	a := &Adapter{
		ms: suite.ms,
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	suite.ms.EXPECT().Status(gomock.Any()).Return(entity.ServiceStatus{
		Status: entity.StatusOk,
	})
	a.Health(w, r)
	body, _ := io.ReadAll(w.Body)
	suite.Equal([]byte(`{"status":"OK","plugins":null}`), body)
	suite.Equal(http.StatusOK, w.Result().StatusCode)
}

func (suite *httpAdapterTestSuite) TestReloadSuccess() {
	a := &Adapter{
		rs: suite.rs,
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/-/reload", nil)
	suite.rs.EXPECT().Load(gomock.Any()).Return(nil)
	a.Reload(w, r)
	suite.Equal(http.StatusOK, w.Result().StatusCode)
}

func (suite *httpAdapterTestSuite) TestReloadFailed() {
	a := &Adapter{
		rs: suite.rs,
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/-/reload", nil)
	suite.rs.EXPECT().Load(gomock.Any()).Return(errors.New("unexpected error"))
	a.Reload(w, r)
	suite.Equal(http.StatusInternalServerError, w.Result().StatusCode)
}

func (suite *httpAdapterTestSuite) TestForwardSuccess() {
	a := &Adapter{
		hs: suite.hs,
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/hook/test", nil)
	suite.hs.EXPECT().Forward(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	a.ForwardHook(w, r)
	suite.Equal(http.StatusOK, w.Result().StatusCode)
}

func (suite *httpAdapterTestSuite) TestForwardFailed() {
	a := &Adapter{
		hs:  suite.hs,
		log: &log.Blackhole{},
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/hook/test", nil)
	suite.hs.EXPECT().Forward(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
	a.ForwardHook(w, r)
	suite.Equal(http.StatusNotAcceptable, w.Result().StatusCode)
}

func (suite *httpAdapterTestSuite) TestForwardSuccessWithPlainBody() {
	a := &Adapter{
		hs: suite.hs,
	}
	w := httptest.NewRecorder()
	buf := bytes.NewBufferString("hello")
	r := httptest.NewRequest(http.MethodPost, "/hook/test", buf)
	suite.hs.EXPECT().Forward(gomock.Any(), gomock.Any(), []byte("hello")).Return(nil)
	a.ForwardHook(w, r)
	suite.Equal(http.StatusOK, w.Result().StatusCode)
}

func (suite *httpAdapterTestSuite) TestForwardSuccessWithForm() {
	a := &Adapter{
		hs: suite.hs,
	}
	w := httptest.NewRecorder()
	buf := bytes.NewBufferString("message=hello")
	r := httptest.NewRequest(http.MethodPost, "/hook/test", buf)
	r.Header.Add("content-type", "application/x-www-form-urlencoded")
	suite.hs.EXPECT().Forward(gomock.Any(), gomock.Any(), []byte("{\"message\":[\"hello\"]}")).Return(nil)
	a.ForwardHook(w, r)
	suite.Equal(http.StatusOK, w.Result().StatusCode)
}

func (suite *httpAdapterTestSuite) TestForwardSuccessWithJSON() {
	a := &Adapter{
		hs: suite.hs,
	}
	w := httptest.NewRecorder()
	buf := bytes.NewBufferString("{\"message\":\"hello\"}")
	r := httptest.NewRequest(http.MethodPost, "/hook/test", buf)
	r.Header.Add("content-type", "application/json")
	suite.hs.EXPECT().Forward(gomock.Any(), gomock.Any(), []byte("{\"message\":\"hello\"}")).Return(nil)
	a.ForwardHook(w, r)
	suite.Equal(http.StatusOK, w.Result().StatusCode)
}
