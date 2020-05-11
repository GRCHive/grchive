package worker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type DockerWorker struct {
	containerName string
	workDir       string
	mavenDir      string
}

func (w DockerWorker) Run() (int, error) {
	return runKotlinContainer(w.containerName)
}

func (w DockerWorker) Logs() (string, error) {
	return readLogsFromContainer(w.containerName)
}

func (w DockerWorker) Cleanup() {
	removeKotlinContainer(w.containerName)
	os.RemoveAll(w.workDir)
	os.RemoveAll(w.mavenDir)
}

type DockerFactory struct {
}

func (f DockerFactory) CreateWorker(opts WorkerOptions) (Worker, error) {
	worker := DockerWorker{
		containerName: fmt.Sprintf("script-runner-%d", opts.RunId),
	}

	var err error
	worker.workDir, err = ioutil.TempDir("", "script-runner")
	if err != nil {
		return nil, err
	}

	worker.mavenDir, err = ioutil.TempDir("", "maven-dir")
	if err != nil {
		return nil, err
	}

	fmt.Printf("Work Dir: %s\n", worker.workDir)
	fmt.Printf("Maven Dir: %s\n", worker.mavenDir)

	err = createKotlinContainer(
		worker.workDir,
		worker.containerName,
		worker.mavenDir,
		"-group",
		opts.ClientLibGroupId,
		"-artifact",
		opts.ClientLibArtifactId,
		"-vers",
		opts.ClientLibVersion,
		"-class",
		opts.RunClassName,
		"-fn",
		opts.RunFunctionName,
		"-meta",
		opts.RunMetadataName,
		"-runId",
		strconv.FormatInt(opts.RunId, 10),
	)
	if err != nil {
		return nil, err
	}

	return worker, nil
}
