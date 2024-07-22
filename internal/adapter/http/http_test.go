package http

import (
	"testing"

	"github.com/k1nky/tookhook/internal/adapter/http/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type httpAdapterTestSuite struct {
	suite.Suite
	hookService *mock.MockhookService
}

func TestHTTPAdapter(t *testing.T) {
	suite.Run(t, new(httpAdapterTestSuite))
}

func (suite *httpAdapterTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.hookService = mock.NewMockhookService(ctrl)
}

// func (suite *httpAdapterTestSuite) TestHealth() {
// 	a := &Adapter{
// 		hooker: suite.hookService,
// 	}
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest(http.MethodGet, "/health", nil)
// 	a.Health(w, r)
// 	body, _ := io.ReadAll(w.Body)
// 	suite.Equal([]byte(`{"status": "OK"}`), body)
// 	suite.Equal(http.StatusOK, w.Result().StatusCode)
// }

// a := &Adapter{
// 	auth: suite.authService,
// }
// for _, tt := range tests {
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.payload))
// 	if len(tt.expectLogin) > 0 {
// 		suite.authService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(tt.expectLogin...)
// 	}
// 	a.Login(w, r)
// 	suite.Equal(tt.want.statusCode, w.Code)
// }
