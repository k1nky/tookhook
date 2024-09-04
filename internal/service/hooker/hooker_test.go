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

func (suite *serviceHookerSuite) TestForward_NotFound() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(nil)
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.ErrorIs(err, entity.ErrNotFound)
}

func (suite *serviceHookerSuite) TestForward_RuleDisabled() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&entity.Hook{
		Income:   "test",
		Disabled: true,
	})
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_NoPlugin() {
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

func (suite *serviceHookerSuite) TestForward_Success() {
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
	suite.tq.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(nil)
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_EnququeFailed() {
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
		ForwardResultError: nil,
	})
	suite.tq.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(errors.New("enquque failed"))
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_MultiplePlugins() {
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

func (suite *serviceHookerSuite) TestForwardMultiple_PluginsDisabled() {
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
	suite.tq.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	err := suite.svc.Forward(context.TODO(), "test", nil)
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_HandlerNotMatch() {
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
