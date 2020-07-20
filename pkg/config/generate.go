package config

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
)

const (
	RootDir               = ".gebug"
	Path                  = "config.yaml"
	DockerfileName        = "Dockerfile"
	DockerComposeFileName = "docker-compose.yml"
)

func FilePath(workDir string, fileName string) string {
	return path.Join(workDir, RootDir, fileName)
}

func createConfigFile(fileName string, workDir string, renderFunc func(io.Writer) error) error {
	filePath := FilePath(workDir, fileName)
	zap.L().Debug("Generating config file", zap.String("path", filePath))
	file, err := os.Create(filePath)
	if err != nil {
		return errors.WithMessagef(err, "create file '%s'", fileName)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			zap.L().Error("Failed to close file", zap.String("path", filePath))
		}
	}()

	err = renderFunc(file)
	if err != nil {
		return errors.WithMessagef(err, "generate file content: '%s'", fileName)
	}

	return nil
}

func (c *Config) Generate(workDir string) error {
	for fileName, renderFunc := range map[string]func(io.Writer) error{
		DockerComposeFileName: c.RenderDockerComposeFile,
		DockerfileName:        c.RenderDockerfile,
	} {
		if err := createConfigFile(fileName, workDir, renderFunc); err != nil {
			return err
		}
	}
	return nil
}
