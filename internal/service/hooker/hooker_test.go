package hooker

import (
	"context"
	"errors"
	"testing"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/internal/entity/hooks"
	"github.com/k1nky/tookhook/internal/entity/rules"
	"github.com/k1nky/tookhook/internal/entity/tasks"
	"github.com/k1nky/tookhook/internal/service/hooker/mock"
	log "github.com/k1nky/tookhook/pkg/logger"
	pluginmock "github.com/k1nky/tookhook/pkg/plugin/mock"
	"github.com/k1nky/tookhook/pkg/thstrings/restrings"
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
	err := suite.svc.Forward(context.TODO(), &hooks.HookRequest{
		Meta: hooks.HookRequestMeta{Name: "test"},
	})
	suite.ErrorIs(err, entity.ErrNotFound)
}

func (suite *serviceHookerSuite) TestForward_RuleDisabled() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&rules.Hook{
		Income:   "test",
		Disabled: true,
	})
	err := suite.svc.Forward(context.TODO(), &hooks.HookRequest{
		Meta: hooks.HookRequestMeta{Name: "test"},
	})
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_NoPlugin() {
	suite.pm.EXPECT().Get(gomock.Any()).Return(nil)
	err := suite.svc.processForward(context.TODO(), &tasks.ForwardTask{})
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_Success() {
	suite.pm.EXPECT().Get(gomock.Any()).Return(&pluginmock.MockPlugin{
		ForwardResultData:  []byte("success"),
		ForwardResultError: nil,
	})
	err := suite.svc.processForward(context.TODO(), &tasks.ForwardTask{})
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_EnququeFailed() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&rules.Hook{
		Income: "test",
		Handlers: []*rules.Handler{
			{
				Type: "plugin1",
			},
		},
	})
	suite.tq.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(errors.New("enquque failed"))
	err := suite.svc.processHook(context.TODO(), &tasks.HookTask{})
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_MultipleHandlers() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&rules.Hook{
		Income: "test",
		Handlers: []*rules.Handler{
			{
				Type: "plugin1",
			},
			{
				Type: "plugin2",
			},
		},
	})
	suite.tq.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	err := suite.svc.processHook(context.TODO(), &tasks.HookTask{})
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_MultipleHandlersDisabled() {
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&rules.Hook{
		Income: "test",
		Handlers: []*rules.Handler{
			{
				Type:     "plugin1",
				Disabled: true,
			},
			{
				Type: "plugin2",
			},
		},
	})
	suite.tq.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	err := suite.svc.processHook(context.TODO(), &tasks.HookTask{})
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_HandlerNotMatch() {
	h := &rules.Handler{
		Type: "plugin1",
		On:   *restrings.New("123"),
	}
	h.Compile()
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&rules.Hook{
		Income:   "test",
		Handlers: []*rules.Handler{h},
	})
	err := suite.svc.processHook(context.TODO(), &tasks.HookTask{
		Data: []byte("abc"),
	})
	suite.NoError(err)
}

func (suite *serviceHookerSuite) TestForward_Discard() {
	h := &rules.Handler{
		Type: "plugin1",
		PreActions: rules.Actions{
			&rules.Action{
				Type: rules.ActionTypeDiscard,
			},
		},
	}
	h.Compile()
	suite.store.EXPECT().GetIncomeHookByName(gomock.Any(), gomock.Any()).Return(&rules.Hook{
		Income:   "test",
		Handlers: []*rules.Handler{h},
	})
	err := suite.svc.processHook(context.TODO(), &tasks.HookTask{
		Data: []byte("abc"),
	})
	suite.NoError(err)
}
