package gitea

type GiteaApi interface {
	// Initialization
	MustInitialize(cfg GiteaConfig)

	// Admin
	AdminCreateUser(GiteaUser) error
	AdminCreateOrganization(GiteaUser, GiteaOrganization) error

	// User
	UserCreateAccessToken(GiteaUser) (GiteaToken, error)

	// Repository
	RepositoryCreate(GiteaUser, GiteaRepository) error
	RepositoryTransfer(GiteaUserlike, GiteaUserlike, GiteaRepository) error
}
