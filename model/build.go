// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package model

type Build struct {
	InstallationID int
	SHA            string
	Owner          string
	Repository     string
	Branch         string
}

func (b *Build) Serialize() map[string]interface{} {
	m := make(map[string]interface{})
	m["installation_id"] = b.InstallationID
	m["sha"] = b.SHA
	m["owner"] = b.Owner
	m["repository"] = b.Repository
	m["branch"] = b.Branch
	return m
}
