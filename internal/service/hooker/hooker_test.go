package hooker

import (
	"context"
	"errors"
	"testing"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/internal/service/hooker/mock"
	log "github.com/k1nky/tookhook/pkg/logger"
	pluginmock "github.com/k1nky/tookhook/pkg/plugin/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type serviceHookerSuite struct {
	suite.Suite
	store *mock.MockrulesStore
	pm    *mock.Mockpluginmanager
	tq    *mock.Mocktaskqueue
	svc   *Service
}

func TestService(t *testing.T) {
	suite.Run(t, new(serviceHookerSuite))
}

func (suite *serviceHookerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.pm = mock.NewMockpluginmanager(ctrl)
	suite.store = mock.NewMockrulesStore(ctrl)
	suite.tq = mock.NewMocktaskqueue(ctrl)
	suite.svc = New(suite.store, suite.pm, &log.Blackhole{}, suite.tq)
}

func (suite *serviceHookerSuite) TestForwardNotFound() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(nil)
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.ErrorIs(err, entity.ErrNotFound)
}

func (suite *serviceHookerSuite) TestForwardRuleDisabled() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&entity.Hook{
		Income:   "test",
		Disabled: true,
	})
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForwardNoPlugin() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&entity.Hook{
		Income: "test",
		Handlers: []*entity.Handler{
			{
				Type: "plugin1",
			},
		},
	})
	suite.pm.EXPECT().Get(gomock.Any()).Return(nil)
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForwardSuccess() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&entity.Hook{
		Income: "test",
		Handlers: []*entity.Handler{
			{
				Type: "plugin1",
			},
		},
	})
	suite.pm.EXPECT().Get(gomock.Any()).Return(&pluginmock.MockPlugin{
		ForwardResultData:  []byte("success"),
		ForwardResultError: nil,
	})
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForwardPluginFailed() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&entity.Hook{
		Income: "test",
		Handlers: []*entity.Handler{
			{
				Type: "plugin1",
			},
		},
	})
	suite.pm.EXPECT().Get(gomock.Any()).Return(&pluginmock.MockPlugin{
		ForwardResultData:  nil,
		ForwardResultError: errors.New("unexpected error"),
	})
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.Error(err)
}

func (suite *serviceHookerSuite) TestForwardMultiplePlugins() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&entity.Hook{
		Income: "test",
		Handlers: []*entity.Handler{
			{
				Type: "plugin1",
			},
			{
				Type: "plugin2",
			},
		},
	})
	suite.pm.EXPECT().Get(gomock.Any()).Return(&pluginmock.MockPlugin{
		ForwardResultData:  nil,
		ForwardResultError: nil,
	}).Times(2)
	suite.tq.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForwardMultiplePluginsDisabled() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&entity.Hook{
		Income: "test",
		Handlers: []*entity.Handler{
			{
				Type:     "plugin1",
				Disabled: true,
			},
			{
				Type: "plugin2",
			},
		},
	})
	suite.pm.EXPECT().Get(gomock.Any()).Return(&pluginmock.MockPlugin{
		ForwardResultData:  nil,
		ForwardResultError: nil,
	}).Times(1)
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForwardHandlerNotMatch() {
	h := &entity.Handler{
		Type: "plugin1",
		On:   "123",
	}
	h.Compile()
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&entity.Hook{
		Income:   "test",
		Handlers: []*entity.Handler{h},
	})
	suite.pm.EXPECT().Get(gomock.Any()).Times(0)
	err := suite.svc.Forward(context.TODO(), "test", []byte("abc"))
	suite.NoError(err)
}
