package apiconfig

import "time"

// GetBaseCfg reads and sets the base api's base environment variable base
func GetBaseCfg(cfg map[string]any) *CfgBase {
	return &CfgBase{
		Port:      cfg["port"].(string),
		Timeout:   time.Duration(cfg["timeout"].(int)) * time.Second,
		URIPrefix: cfg["uri_prefix"].(string),
	}
}
