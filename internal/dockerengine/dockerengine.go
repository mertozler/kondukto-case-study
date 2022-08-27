package dockerengine

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/mertozler/internal/config"
	"github.com/mertozler/internal/models"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Engine interface {
	NewScanResults(projectId string) (error, models.ScanData)
}

type Docker struct {
	Client          *client.Client
	ContainerConfig *container.Config
	imageName       string
}

func NewDockerEngine(config *config.DockerEngine) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}
	return &Docker{
		Client:    cli,
		imageName: config.Imagename,
		ContainerConfig: &container.Config{
			Image:        config.Imagename,
			Cmd:          []string{"-r", "/code", "-f", "json"},
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
		},
	}, nil
}

func (e *Docker) NewScanResults(scanId string) (error, models.ScanData) {
	ctx := context.Background()
	err, targetDirectory := getTargetDirectoryPath(scanId)
	if err != nil {
		return err, models.ScanData{}
	}

	logrus.Infof("%s image is pulling", e.imageName)
	out, err := e.Client.ImagePull(ctx, e.imageName, types.ImagePullOptions{})
	if err != nil {
		return err, models.ScanData{}
	}

	defer out.Close()
	io.Copy(os.Stdout, out)
	binds := []string{
		targetDirectory + ":/code",
	}

	logrus.Info("Container creating in docker")
	resp, err := e.Client.ContainerCreate(ctx, e.ContainerConfig, &container.HostConfig{
		Binds: binds,
	}, nil, nil, "")
	if err != nil {
		return err, models.ScanData{}
	}

	logrus.Info("Container starting in docker")
	if err := e.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err, models.ScanData{}
	}
	statusCh, errCh := e.Client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err, models.ScanData{}
		}
	case <-statusCh:
	}

	logrus.Info("Getting security results from bandit logs for ", scanId)
	reader, _ := e.Client.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	defer reader.Close()

	p := make([]byte, 8)
	reader.Read(p)
	logDataFromBandit, _ := ioutil.ReadAll(reader)

	scanData, err := getScanDataFromBanditLogs(scanId, logDataFromBandit)

	return nil, scanData
}

func getScanDataFromBanditLogs(scanId string, logDataFromBandit []byte) (models.ScanData, error) {
	logData := string(logDataFromBandit)
	firstIndex := strings.Index(logData, "{")
	var scanDataFromBanditLogs models.ScanDatas
	if err := json.Unmarshal([]byte(logData[firstIndex:]), &scanDataFromBanditLogs); err != nil {
		return models.ScanData{}, err
	}
	var scanData models.ScanData
	scanData.ScanID = scanId
	scanData.ScanData = scanDataFromBanditLogs
	return scanData, nil
}

func getTargetDirectoryPath(scanId string) (error, string) {
	dir, err := os.Getwd()
	if err != nil {
		return err, ""
	}
	targetDirectory := dir + "/tmp/src/" + scanId
	return nil, targetDirectory
}
