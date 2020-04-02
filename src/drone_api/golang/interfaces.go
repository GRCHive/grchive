package drone

type DroneApi interface {
	// Initialization
	MustInitialize(DroneConfig)

	// Repos
	RepoEnable(owner string, repo string) error
	RepoSync() error
}
