package entity

const (
	StatusOk     Status = "OK"
	StatusFailed Status = "Failed"
)

type Status string

type PluginsStatus map[string]Status

type ServiceStatus struct {
	Status  Status        `json:"status"`
	Plugins PluginsStatus `json:"plugins"`
}
