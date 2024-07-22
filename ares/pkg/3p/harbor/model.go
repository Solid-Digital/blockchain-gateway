package harbor

import (
	"encoding/json"
	"strconv"
)

type ProjectMetadata struct {
	AutoScan           bool   `json:"auto_scan"`
	EnableContentTrust bool   `json:"enable_content_trust"`
	PreventVul         bool   `json:"prevent_vul"`
	Public             bool   `json:"public"`
	Severity           string `json:"severity"`
}

func (m *ProjectMetadata) UnmarshalJSON(data []byte) error {
	var err error

	dataMap := make(map[string]string)

	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	if public, ok := dataMap["public"]; ok {
		m.Public, err = strconv.ParseBool(public)

		if err != nil {
			return err
		}
	}

	if autoScan, ok := dataMap["auto_scan"]; ok {
		m.AutoScan, err = strconv.ParseBool(autoScan)

		if err != nil {
			return err
		}
	}

	if enableContentTrust, ok := dataMap["enable_content_trust"]; ok {
		m.EnableContentTrust, err = strconv.ParseBool(enableContentTrust)
		if err != nil {
			return err
		}
	}

	if preventVul, ok := dataMap["prevent_vul"]; ok {
		m.PreventVul, err = strconv.ParseBool(preventVul)
		if err != nil {
			return err
		}
	}

	m.Severity = dataMap["severity"]

	return nil
}

type Project struct {
	ID                int    `json:"project_id"`
	OwnerID           int    `json:"owner_id"`
	Name              string `json:"name"`
	CreationTime      string `json:"creation_time"`
	UpdateTime        string `json:"update_time"`
	OwnerName         string `json:"owner_name"`
	Togglable         bool   `json:"togglable"`
	CurrentUserRoleID int    `json:"current_user_role_id"`
	RepoCount         int    `json:"repo_count"`

	Metadata ProjectMetadata `json:"metadata"`
}
