package rapididentity

const path = "/api/rest/admin/announcements"
const credsFile = "/.rapididentity/credential"

type Announcement struct {
	Id        string `json:"id,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Acl       *Acl   `json:"acl,omitempty"`
	Message   string `json:"message,omitempty"`
	ShowOnce  bool   `json:"showOnce,omitempty"`
	Read      bool   `json:"read,omitempty"`
}

type Acl struct {
	GroupAclsEnabled   bool        `json:"groupAclsEnabled,omitempty"`
	GroupAcls          []*GroupAcl `json:"groupAcls,omitempty"`
	GroupExclusionAcls []*GroupAcl `json:"groupExclusionAcls,omitempty"`
	FilterAclEnabled   bool        `json:"filterAclEnabled,omitempty"`
	FilterAcl          string      `json:"filterAcl,omitempty"`
}

type GroupAcl struct {
	Dn          string `json:"dn,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
