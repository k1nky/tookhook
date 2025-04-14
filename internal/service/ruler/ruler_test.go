package ruler

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/k1nky/tookhook/internal/entity/rules"
	"github.com/k1nky/tookhook/internal/service/ruler/mock"
	log "github.com/k1nky/tookhook/pkg/logger"
	"github.com/k1nky/tookhook/pkg/plugin"
	pluginmock "github.com/k1nky/tookhook/pkg/plugin/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestServiceGetIncomeHookByName(t *testing.T) {
	type fields struct {
		rules *rules.Rules
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *rules.Hook
	}{
		{
			name: "not found",
			fields: fields{
				rules: &rules.Rules{
					Hooks: []rules.Hook{
						{
							Income: "first",
						},
					},
				},
			},
			args: args{
				ctx:  context.TODO(),
				name: "not_found",
			},
			want: nil,
		},
		{
			name: "found",
			fields: fields{
				rules: &rules.Rules{
					Hooks: []rules.Hook{
						{Income: "first"},
						{Income: "second"},
					},
				},
			},
			args: args{
				ctx:  context.TODO(),
				name: "second",
			},
			want: &rules.Hook{Income: "second"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				rules: tt.fields.rules,
			}
			if got := svc.GetIncomeHookByName(tt.args.ctx, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetIncomeHookByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

type serviceValidateSuite struct {
	suite.Suite
	pm  *mock.Mockpluginmanager
	svc *Service
}

func TestServiceValidate(t *testing.T) {
	suite.Run(t, new(serviceValidateSuite))
}

func (suite *serviceValidateSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.pm = mock.NewMockpluginmanager(ctrl)
	suite.svc = New(suite.pm, nil, &log.Blackhole{})
}

func (suite *serviceValidateSuite) TestInvalidBaseStructure() {
	rules := &rules.Rules{
		Hooks: []rules.Hook{
			{
				Income: "",
				Handlers: []*rules.Handler{
					{
						Type: "plugin_name",
					},
				},
			},
		},
	}
	suite.pm.EXPECT().Get(gomock.Any()).Times(0)
	err := suite.svc.Validate(context.TODO(), rules)
	suite.Error(err)
}

func (suite *serviceValidateSuite) TestValidate() {
	rules := &rules.Rules{
		Hooks: []rules.Hook{
			{
				Income: "Rule1",
				Handlers: []*rules.Handler{
					{
						Type: "plugin_name",
					},
				},
			},
		},
	}
	tests := []struct {
		name      string
		plugin    plugin.Plugin
		wantError bool
	}{
		{
			name: "valid",
			plugin: &pluginmock.MockPlugin{
				ValidateResult: nil,
			},
			wantError: false,
		},
		{
			name: "plugin return error",
			plugin: &pluginmock.MockPlugin{
				ValidateResult: errors.New("plugin: invalid value"),
			},
			wantError: true,
		},
		{
			name:      "plugin not found",
			plugin:    nil,
			wantError: false,
		},
	}
	for _, tt := range tests {
		suite.pm.EXPECT().Get(gomock.Any()).Return(tt.plugin)
		err := suite.svc.Validate(context.TODO(), rules)
		suite.Equal(tt.wantError, err != nil)
	}
}

type serviceLoadSuite struct {
	suite.Suite
	pm    *mock.Mockpluginmanager
	store *mock.Mockstorage
	svc   *Service
}

func TestServiceLoad(t *testing.T) {
	suite.Run(t, new(serviceLoadSuite))
}

func (suite *serviceLoadSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.pm = mock.NewMockpluginmanager(ctrl)
	suite.store = mock.NewMockstorage(ctrl)
	suite.svc = New(suite.pm, suite.store, &log.Blackhole{})
}

func (suite *serviceLoadSuite) TestGetRulesFailed() {
	before := &rules.Rules{
		Hooks: []rules.Hook{
			{Income: "Rule1"},
		},
	}
	suite.svc.rules = before
	suite.store.EXPECT().GetRules(gomock.Any()).Return(nil, errors.New("failed"))
	err := suite.svc.Load(context.TODO())
	suite.Error(err)
	suite.Equal(before, suite.svc.rules)
}

func (suite *serviceLoadSuite) TestValidateFailed() {
	before := &rules.Rules{
		Hooks: []rules.Hook{
			{Income: "Rule1"},
		},
	}
	suite.svc.rules = before
	suite.store.EXPECT().GetRules(gomock.Any()).Return(&rules.Rules{
		Hooks: []rules.Hook{
			{Income: ""},
		},
	}, nil)
	err := suite.svc.Load(context.TODO())
	suite.Error(err)
	suite.Equal(before, suite.svc.rules)
}

func (suite *serviceLoadSuite) TestSuccess() {
	before := &rules.Rules{
		Hooks: []rules.Hook{
			{Income: "Rule1"},
		},
	}
	after := &rules.Rules{
		Hooks: []rules.Hook{
			{Income: "Rule2"},
		},
	}
	suite.svc.rules = before
	suite.store.EXPECT().GetRules(gomock.Any()).Return(after, nil)
	err := suite.svc.Load(context.TODO())
	suite.NoError(err)
	suite.Equal(after, suite.svc.rules)
}
