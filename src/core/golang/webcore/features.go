package webcore

func EnableAutomationFeature(orgId int32) error {
	// This functioln needs to setup organization specific things
	// in Gitea and Drone CI.
	// 	1) Create an organization and repository for holding all the
	// 	   org's Kotlin code in Gitea.
	// 	2) Create an organization specific user for us to assume the identity of.
	//  3) Put in default template code for the project.
	//  4) Enable the repository in Drone CI.

	return nil
}
