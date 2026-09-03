package response

import "shack/internal/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
