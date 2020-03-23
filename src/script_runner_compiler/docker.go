package main

import (
	"archive/tar"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"gitlab.com/grchive/grchive/core"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const dockerWorkspaceDir string = "/data"
const dockerOutputDir string = "/output"

func mustDockerCreateClient(c *client.Client, err error) *client.Client {
	if err != nil {
		core.Error("Failed to create client: " + err.Error())
	}
	return c
}

var dockerClient *client.Client = mustDockerCreateClient(client.NewEnvClient())

func createKotlinContainer(workspaceDir string, containerName string, workspaceVolumeName string, settings core.ScriptRunSettings) error {
	_, err := dockerClient.VolumeCreate(context.Background(), volume.VolumeCreateBody{
		Driver: "local",
		DriverOpts: map[string]string{
			"type":   "tmpfs",
			"device": "tmpfs",
			"o":      fmt.Sprintf("size=%d", settings.DiskSizeBytes),
		},
		Name: workspaceVolumeName,
	})

	if err != nil {
		return err
	}

	inputDir := "/input"

	args := strslice.StrSlice{}
	if settings.CompileOnly {
		args = append(args, "--compile_only")
	}

	args = append(args, "--script", getExpectedScriptFname(inputDir))
	if settings.ScriptChecksum != "" {
		args = append(args, "--checksum", settings.ScriptChecksum)
	}

	args = append(args, "--jar", getExpectedJarFname(inputDir))

	if settings.GrchiveCoreVersion != "" {
		args = append(args, "--version", settings.GrchiveCoreVersion)
	}

	args = append(args, "--output", getExpectedCompiledJarFname(dockerOutputDir))

	const defaultCPUPeriod int64 = 100000

	_, err = dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: fmt.Sprintf("registry.gitlab.com/grchive/grchive/kotlin_runner:%s", settings.KotlinContainerVersion),
			Cmd:   args,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				mount.Mount{
					Type:     mount.TypeBind,
					Source:   workspaceDir,
					Target:   inputDir,
					ReadOnly: true,
				},
				mount.Mount{
					Type:   mount.TypeVolume,
					Source: workspaceVolumeName,
					Target: dockerWorkspaceDir,
				},
			},
			Resources: container.Resources{
				Memory:    settings.MemBytesAllocation,
				CPUPeriod: defaultCPUPeriod,
				CPUQuota:  int64(math.Round(float64(defaultCPUPeriod) * settings.CpuAllocation)),
			},
		},
		&network.NetworkingConfig{},
		containerName,
	)

	if err != nil {
		return err
	}

	return nil
}

func runKotlinContainer(containerName string, settings core.ScriptRunSettings) error {
	err := dockerClient.ContainerStart(
		context.Background(),
		containerName,
		types.ContainerStartOptions{},
	)

	if err != nil {
		return err
	}

	// Block until container is done running.
	for {
		json, err := dockerClient.ContainerInspect(
			context.Background(),
			containerName,
		)

		if err != nil {
			return err
		}

		if !json.ContainerJSONBase.State.Running {
			break
		}

		time.Sleep(1 * time.Second)
	}
	return nil
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
		} else if header[1] == 2 {
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

func removeKotlinContainer(containerName string, workspaceVolumeName string) error {
	err := dockerClient.ContainerRemove(
		context.Background(),
		containerName,
		types.ContainerRemoveOptions{},
	)
	if err != nil {
		return err
	}

	return dockerClient.VolumeRemove(context.Background(), workspaceVolumeName, true)
}
