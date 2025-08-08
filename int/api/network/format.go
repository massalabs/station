package network

import (
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/int/config"
)

func buildAvailableNetworkInfos(cfgManager *config.MSConfigManager) []*models.NetworkInfoItem {
	networks := make([]*models.NetworkInfoItem, 0, len(cfgManager.Network.Networks))
	for i := range cfgManager.Network.Networks {
		nfo := cfgManager.Network.Networks[i]
		status := string(nfo.Status())
		name := nfo.Name
		url := nfo.NodeURL
		version := nfo.Version
		chainID := int64(nfo.ChainID)
		networks = append(networks, &models.NetworkInfoItem{
			Name:    name,
			URL:     url,
			Version: version,
			ChainID: chainID,
			Status:  status,
		})
	}
	return networks
}
