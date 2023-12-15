package config

type Version struct {
	version string
	commit  string
	date    string
}

func NewVersion() *Version {
	return &Version{}
}

func (v *Version) WithVersion(version string) *Version {
	v.version = version
	return v
}

func (v *Version) WithCommit(commit string) *Version {
	v.commit = commit
	return v
}

func (v *Version) WithDate(date string) *Version {
	v.date = date
	return v
}

func (v *Version) Version() string {
	return v.version
}

func (v *Version) Commit() string {
	return v.commit
}

func (v *Version) Date() string {
	return v.date
}
