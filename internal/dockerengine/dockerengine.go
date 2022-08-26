package dockerengine

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/mertozler/internal/models"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func NewScanResults(projectId string) (error, models.ScanData) {
	ctx := context.Background()
	err, targetDirectory := getTargetDirectoryPath(projectId)
	if err != nil {
		return err, models.ScanData{}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err, models.ScanData{}
	}
	imageName := "opensorcery/bandit"
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err, models.ScanData{}
	}

	defer out.Close()
	io.Copy(os.Stdout, out)
	binds := []string{
		targetDirectory + ":/code",
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		Cmd:          []string{"-r", "/code", "-f", "json"},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Binds: binds,
	}, nil, nil, "")
	if err != nil {
		return err, models.ScanData{}
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err, models.ScanData{}
	}
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err, models.ScanData{}
		}
	case <-statusCh:
	}

	reader, _ := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	defer reader.Close()

	//read the first 8 bytes to ignore the HEADER part from docker container logs
	p := make([]byte, 8)
	reader.Read(p)
	content, _ := ioutil.ReadAll(reader)

	scanData, err := getScanDataFromBanditLogs(projectId, content)

	return nil, scanData
}

func getScanDataFromBanditLogs(projectId string, content []byte) (models.ScanData, error) {
	logData := string(content)
	firstIndex := strings.Index(logData, "{")
	data := cleanStringData(logData[firstIndex:])
	var input interface{}
	if err := json.Unmarshal([]byte(data), &input); err != nil {
		return models.ScanData{}, err
	}
	var scanData models.ScanData
	scanData.ScanID = projectId
	scanData.ScanData = input
	return scanData, nil
}

func cleanStringData(data string) string {
	var string_b string = data

	string_b = strings.ReplaceAll(string_b, "\r\n ", "")
	string_b = strings.TrimSpace(string_b)
	return string_b
}

func getTargetDirectoryPath(projectId string) (error, string) {
	dir, err := os.Getwd()
	if err != nil {
		return err, ""
	}
	targetDirectory := dir + "/tmp/src/" + projectId
	return nil, targetDirectory
}
