package drone

type DroneApi interface {
	// Initialization
	MustInitialize(DroneConfig)

	// Repos
	RepoEnable(owner string, repo string) error
	RepoSync() error

	// Builds
	BuildCreate(owner string, repo string, params map[string]string) error
}
