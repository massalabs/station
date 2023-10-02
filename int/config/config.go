package config

type VersionInfo struct {
	Version string
}

//nolint:gochecknoglobals
var Version = &VersionInfo{
	Version: "dev",
}
