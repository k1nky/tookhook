package monitor

import (
	"context"
	"testing"

	"github.com/k1nky/tookhook/internal/entity"
	log "github.com/k1nky/tookhook/internal/logger"
	"github.com/k1nky/tookhook/internal/service/monitor/mock"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type serviceMonitorSuite struct {
	suite.Suite
	pm    *mock.Mockpluginmanager
	store *mock.Mockstorage
	hs    *mock.MockhookService
	svc   *Service
}

func TestServiceLoad(t *testing.T) {
	suite.Run(t, new(serviceMonitorSuite))
}

func (suite *serviceMonitorSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.pm = mock.NewMockpluginmanager(ctrl)
	suite.store = mock.NewMockstorage(ctrl)
	suite.hs = mock.NewMockhookService(ctrl)
	suite.svc = New(suite.pm, suite.hs, suite.store, &log.Blackhole{})
}

func (suite *serviceMonitorSuite) TestGetRulesFailed() {
	tests := []struct {
		name         string
		expected     entity.ServiceStatus
		hsStatus     entity.Status
		pluginStatus entity.PluginsStatus
	}{
		{
			name: "OK",
			expected: entity.ServiceStatus{Status: entity.StatusOk, Plugins: entity.PluginsStatus{
				"plugin": entity.StatusOk,
			}},
			hsStatus: entity.StatusOk,
			pluginStatus: entity.PluginsStatus{
				"plugin": entity.StatusOk,
			},
		},
		{
			name: "hooker service - failed",
			expected: entity.ServiceStatus{Status: entity.StatusFailed, Plugins: entity.PluginsStatus{
				"plugin": entity.StatusOk,
			}},
			hsStatus: entity.StatusFailed,
			pluginStatus: entity.PluginsStatus{
				"plugin": entity.StatusOk,
			},
		},
		{
			name: "plugin - failed",
			expected: entity.ServiceStatus{Status: entity.StatusFailed, Plugins: entity.PluginsStatus{
				"plugin": entity.StatusFailed,
			}},
			hsStatus: entity.StatusOk,
			pluginStatus: entity.PluginsStatus{
				"plugin": entity.StatusFailed,
			},
		},
	}

	for _, tt := range tests {
		suite.hs.EXPECT().Health(gomock.Any()).Return(tt.hsStatus)
		suite.pm.EXPECT().Health(gomock.Any()).Return(tt.pluginStatus)
		got := suite.svc.Status(context.TODO())
		suite.Equal(tt.expected, got)
	}
}
