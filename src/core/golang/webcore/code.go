package webcore

import (
	"errors"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	drone "gitlab.com/grchive/grchive/drone_api"
	gitea "gitlab.com/grchive/grchive/gitea_api"
	"gopkg.in/yaml.v2"
	"strconv"
)

func GetManagedCodeFromGitea(codeId int64, orgId int32, role *core.Role) (string, error) {
	grcRepo, err := database.GetLinkedGiteaRepository(orgId)
	if err != nil {
		return "", err
	}

	code, err := database.GetCode(codeId, orgId, role)
	if err != nil {
		return "", err
	}

	content, _, err := gitea.GlobalGiteaApi.RepositoryGetFile(gitea.GiteaRepository{
		Owner: grcRepo.GiteaOrg,
		Name:  grcRepo.GiteaRepo,
	}, code.GitPath, code.GitHash)
	return content, err
}

func DeleteManagedCodeFromGitea(code *core.ManagedCode) error {
	grcRepo, err := database.GetLinkedGiteaRepository(code.OrgId)
	if err != nil {
		return err
	}

	giteaRepo := gitea.GiteaRepository{
		Owner: grcRepo.GiteaOrg,
		Name:  grcRepo.GiteaRepo,
	}

	err = gitea.GlobalGiteaApi.RepositoryDeleteFile(giteaRepo, code.GitPath, gitea.GiteaDeleteFileOptions{
		Message: fmt.Sprintf("Delete file - %s", code.GitPath),
		Sha:     code.GiteaFileSha,
	})

	return err
}

func StoreManagedCodeToGitea(code *core.ManagedCode, script string, role *core.Role, msg string) error {
	grcRepo, err := database.GetLinkedGiteaRepository(code.OrgId)
	if err != nil {
		return err
	}

	err = UpdateGiteaRepositoryTemplate(code.OrgId)
	if err != nil {
		return err
	}

	giteaRepo := gitea.GiteaRepository{
		Owner: grcRepo.GiteaOrg,
		Name:  grcRepo.GiteaRepo,
	}

	// We need to check if the file exists or not in Gitea irrespective of what our database thinks
	// since it's possible for our code to get into a state where the file exists in Gitea but we do not
	// have it tracked in the database.
	_, sha, err := gitea.GlobalGiteaApi.RepositoryGetFile(giteaRepo, code.GitPath, "master")

	// Assume a non-nil error means that the file doesn't exist. If it turns out we're wrong,
	// we should get another error when we try to create a new file.
	if err != nil {
		code.GitHash, code.GiteaFileSha, err = gitea.GlobalGiteaApi.RepositoryCreateFile(
			giteaRepo,
			code.GitPath,
			gitea.GiteaCreateFileOptions{
				Content: script,
				Message: msg,
			},
		)
	} else {
		code.GitHash, code.GiteaFileSha, err = gitea.GlobalGiteaApi.RepositoryUpdateFile(
			giteaRepo,
			code.GitPath,
			gitea.GiteaCreateFileOptions{
				Content: script,
				Message: msg,
			},
			sha,
		)
	}

	if err != nil {
		return err
	}

	// Store information in the database after since it'll be easier to handle cases where
	// files are in Gitea but not in the database than if the file were in the database but not in Gitea.
	if role != nil {
		code.UserId = role.UserId
		return database.InsertManagedCode(code, role)
	}
	return nil
}

func GenerateScriptMetadataYaml(params []*core.CodeParameter, clientDataId []int64) (string, error) {
	type ScriptMetadataData struct {
		Params       []*core.CodeParameter `yaml:"params"`
		ClientDataId []int64               `yaml:"clientDataId"`
	}

	metadata := ScriptMetadataData{
		Params:       params,
		ClientDataId: clientDataId,
	}

	data, err := yaml.Marshal(&metadata)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func RunAuthorizedScriptImmediate(runId int64, approval core.GenericApproval) error {
	if !approval.Response {
		return errors.New("Script run not approved.")
	}

	run, err := database.GetScriptRun(runId, core.ServerRole)
	if err != nil {
		return err
	}

	code, err := database.GetCodeFromScriptCodeLink(run.LinkId, core.ServerRole)
	if err != nil {
		return err
	}

	if run.RequiresBuild {
		repo, err := database.GetLinkedGiteaRepository(code.OrgId)
		if err != nil {
			return err
		}

		// In this case, we need to fire off a Drone CI job to compile the latest code + the current script revision.
		// We must specify branch/commit here due to a Gitea issue with the /repos/{owner}/{repo}/commits/{ref} API endpoint
		// that Drone uses to find the latest commit of the branch. This endpoint doesn't work in Gitea so we need to specify the
		// commit directly.
		commitSha, err := gitea.GlobalGiteaApi.RepositoryGitGetRefSha(
			gitea.GiteaRepository{
				Owner: repo.GiteaOrg,
				Name:  repo.GiteaRepo,
			},
			"refs/heads/master",
		)

		if err != nil {
			return err
		}

		err = drone.GlobalDroneApi.BuildCreate(repo.GiteaOrg, repo.GiteaRepo, map[string]string{
			"branch":     "master",
			"commit":     commitSha,
			"SCRIPT_RUN": strconv.FormatInt(run.Id, 10),
		})

		if err != nil {
			return err
		}
	} else {
		// Grab JAR path from drone CI since that's the only place we store it.
		// This is hitting the same DB table as the GetCodeBuildStatus call earlier, can we merge it somehow?
		jar, err := database.GetCodeJar(code.Id, code.OrgId, core.ServerRole)
		if err != nil {
			return err
		}

		// In this case, we can directly send off a request to make the script runner run this script.
		DefaultRabbitMQ.SendMessage(PublishMessage{
			Exchange: DEFAULT_EXCHANGE,
			Queue:    SCRIPT_RUNNER_QUEUE,
			Body: ScriptRunnerMessage{
				RunId: run.Id,
				Jar:   jar,
			},
		})
	}
	return nil
}

// A bit of a misnomer, should just be called when the scheduled script run
// gets approved.
func RunAuthorizedScriptScheduled(taskId int64, approval core.GenericApproval) error {
	if !approval.Response {
		return errors.New("Script run not approved.")
	}

	DefaultRabbitMQ.SendMessage(PublishMessage{
		Exchange: DEFAULT_EXCHANGE,
		Queue:    TASK_MANAGER_QUEUE,
		Body: TaskManagerMessage{
			Action: "Add",
			TaskId: taskId,
		},
	})
	return nil
}
