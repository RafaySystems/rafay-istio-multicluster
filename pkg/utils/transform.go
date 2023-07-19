package utils

// create a maps for upgrade status

func GetUpgradeStatusMap() map[string]string {
	upgradeStatusMap := make(map[string]string)
	upgradeStatusMap["UpgradeEdgeSuccess"] = "Upgrade Successful"
	upgradeStatusMap["UpgradeEdgeNotConfigured"] = "Upgrade Pending"
	upgradeStatusMap["UpgradeEdgeInProgress"] = "Upgrade In Progress"
	upgradeStatusMap["UpgradeEdgeFailed"] = "Upgrade Failed"
	upgradeStatusMap["UpgradeEdgeConfigured"] = "Upgrade Configured"
	return upgradeStatusMap
}
