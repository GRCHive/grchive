package main

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
)

type RuntimeEnvironment struct {
	ContainerPath string
	CoreLibPath   string
}

func pullRuntimeEnvironment(settings core.ScriptRunSettings) (*RuntimeEnvironment, error) {
	containerUri := fmt.Sprintf("registry.gitlab.com/grchive/grchive/kotlin_runner:%s", settings.KotlinContainerVersion)
	err := pullKotlinImage(containerUri)
	if err != nil {
		return nil, err
	}

	coreLibPath, err := pullCoreLibraryFromVersion(settings.GrchiveCoreVersion)
	if err != nil {
		return nil, err
	}

	env := RuntimeEnvironment{
		ContainerPath: containerUri,
		CoreLibPath:   coreLibPath,
	}
	return &env, nil
}
