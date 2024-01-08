package config

// Version the version information
type Version struct {
	version string
	commit  string
	date    string
}

// NewVersion creating a new version
func NewVersion() *Version {
	return &Version{}
}

// WithVersion setting the version information fluid
func (v *Version) WithVersion(version string) *Version {
	v.version = version
	return v
}

// WithCommit setting the commit information fluid
func (v *Version) WithCommit(commit string) *Version {
	v.commit = commit
	return v
}

// WithDate setting the date information fluid
func (v *Version) WithDate(date string) *Version {
	v.date = date
	return v
}

// Version return in the version information
func (v *Version) Version() string {
	return v.version
}

// Commit return in the commit information
func (v *Version) Commit() string {
	return v.commit
}

// Date return in the date information
func (v *Version) Date() string {
	return v.date
}
