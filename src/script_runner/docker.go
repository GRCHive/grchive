package main

import (
	"archive/tar"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"gitlab.com/grchive/grchive/core"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const dockerWorkspaceDir string = "/data"
const dockerOutputDir string = "/output"

var dockerImageMap map[string]bool = map[string]bool{}
var dockerImageMutex sync.RWMutex = sync.RWMutex{}

func mustDockerCreateClient(c *client.Client, err error) *client.Client {
	if err != nil {
		core.Error("Failed to create client: " + err.Error())
	}

	// Cache a map of every image we have locally.
	imgList, err := c.ImageList(
		context.Background(),
		types.ImageListOptions{
			All: true,
		},
	)

	if err != nil {
		core.Error("Failed to pull image list: " + err.Error())
	}

	for _, img := range imgList {
		if len(img.RepoTags) == 0 {
			continue
		}
		dockerImageMap[img.RepoTags[0]] = true
	}
	return c
}

var dockerClient *client.Client = mustDockerCreateClient(client.NewEnvClient())

const defaultCPUPeriod int64 = 100000

func createKotlinContainer(workspaceDir string, containerName string) error {
	entrypoint := append(strslice.StrSlice{}, "/data/run.sh")
	args := append(strslice.StrSlice{}, "POOOOOOOOP")

	runnerNetwork := os.Getenv("SCRIPT_RUNNER_NETWORK")
	_, err := dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:      core.EnvConfig.Drone.RunnerImage,
			Entrypoint: entrypoint,
			Cmd:        args,
			WorkingDir: "/data",
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				mount.Mount{
					Type:     mount.TypeBind,
					Source:   workspaceDir,
					Target:   "/data",
					ReadOnly: false,
				},
			},
			Resources: container.Resources{
				CPUPeriod: defaultCPUPeriod,
				CPUQuota:  defaultCPUPeriod,
			},
			NetworkMode: container.NetworkMode(runnerNetwork),
		},
		&network.NetworkingConfig{},
		containerName,
	)

	if err != nil {
		return err
	}

	return nil
}

func runKotlinContainer(containerName string) (int, error) {
	err := dockerClient.ContainerStart(
		context.Background(),
		containerName,
		types.ContainerStartOptions{},
	)

	if err != nil {
		return -1, err
	}

	// Block until container is done running.
	for {
		json, err := dockerClient.ContainerInspect(
			context.Background(),
			containerName,
		)

		if err != nil {
			return -1, err
		}

		if !json.ContainerJSONBase.State.Running {
			return json.ContainerJSONBase.State.ExitCode, nil
		}

		time.Sleep(1 * time.Second)
	}
	return -1, nil
}

func copyDataFromContainer(containerName string, containerSrc string, hostDst string) error {
	rc, _, err := dockerClient.CopyFromContainer(
		context.Background(),
		containerName,
		containerSrc)

	if err != nil {
		return err
	}

	defer rc.Close()

	// rc is a reader that points to a TAR archive. Read, output, and untar.
	tarReader := tar.NewReader(rc)
	foundFile := false

	for !foundFile {
		h, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil && err != io.EOF {
			return err
		}

		if h.Name != filepath.Base(containerSrc) {
			continue
		}

		foundFile = true

		data := make([]byte, h.Size)
		totalRead := 0
		for {
			numBytes, err := tarReader.Read(data[totalRead:h.Size])
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			totalRead += numBytes
		}

		err = ioutil.WriteFile(hostDst, data, os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	if !foundFile {
		return errors.New("Failed to find file to copy.")
	}

	return nil
}

func readLogsFromContainer(containerName string) (string, error) {
	rc, err := dockerClient.ContainerLogs(
		context.Background(),
		containerName,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		},
	)

	if err != nil {
		return "", err
	}

	defer rc.Close()

	logBuilder := strings.Builder{}

	for {
		header := make([]byte, 8)
		_, err = rc.Read(header)
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		var prefix string
		if header[0] == 1 {
			prefix = "STDOUT"
		} else if header[0] == 2 {
			prefix = "STDERR"
		}

		lenMessageBytes := binary.BigEndian.Uint32(header[4:len(header)])

		data := make([]byte, lenMessageBytes)
		_, err = rc.Read(data)
		if err != nil && err != io.EOF {
			return "", err
		}

		logBuilder.WriteString(fmt.Sprintf("[%s] %s", prefix, string(data)))
	}

	return logBuilder.String(), nil
}

func removeKotlinContainer(containerName string) error {
	err := dockerClient.ContainerRemove(
		context.Background(),
		containerName,
		types.ContainerRemoveOptions{},
	)
	if err != nil {
		return err
	}

	return nil
}

func pullKotlinImage(imageName string) error {
	dockerImageMutex.RLock()
	_, ok := dockerImageMap[imageName]
	dockerImageMutex.RUnlock()

	if !ok {
		auth := types.AuthConfig{
			Username: core.EnvConfig.GitlabRegistryAuth.Username,
			Password: core.EnvConfig.GitlabRegistryAuth.Password,
		}

		authBytes, err := json.Marshal(auth)
		if err != nil {
			return err
		}

		authBase64 := base64.URLEncoding.EncodeToString(authBytes)

		rc, err := dockerClient.ImagePull(
			context.Background(),
			imageName,
			types.ImagePullOptions{
				RegistryAuth: authBase64,
			},
		)

		if err != nil {
			return err
		}

		defer rc.Close()

		resp, err := dockerClient.ImageLoad(
			context.Background(),
			rc,
			false,
		)

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		dockerImageMutex.Lock()
		dockerImageMap[imageName] = true
		dockerImageMutex.Unlock()
	}

	return nil
}
