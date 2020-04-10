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
	RepositoryTransfer(GiteaUserlike, GiteaUserlike, *GiteaRepository) error
	RepositoryAddCollaborator(GiteaRepository, GiteaUserlike) error

	RepositoryCreateFile(GiteaRepository, string, GiteaCreateFileOptions) (string, error)
	RepositoryUpdateFile(GiteaRepository, string, GiteaCreateFileOptions, string) (string, error)
	RepositoryDeleteFile(GiteaRepository, string, GiteaDeleteFileOptions) error
	RepositoryGetFile(GiteaRepository, string, string) (string, string, error)
}
