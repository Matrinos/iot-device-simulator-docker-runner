package main

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

/**
 * Sample activities used by file processing sample workflow.
 */
const (
	runSimulationActivityName = "runSimulationActivityName"
	// processFileActivityName  = "processFileActivity"
	// uploadFileActivityName   = "uploadFileActivity"
)

func runSimulationActivity(ctx context.Context,
	userName string, password string, port string,
	containerName string) (*container.ContainerCreateCreatedBody, error) {
	// logger := activity.GetLogger(ctx)
	// logger.Info("Running docker", zap.String("FileID", fileID))
	// data := RunContainer(fileID)

	// tmpFile, err := saveToTmpFile(data)
	// if err != nil {
	// 	logger.Error("downloadFileActivity failed to save tmp file.", zap.Error(err))
	// 	return nil, err
	// }

	// fileInfo := &container.ContainerCreateCreatedBody{FileName: tmpFile.Name(), HostID: HostID}
	// logger.Info("downloadFileActivity succeed.", zap.String("SavedFilePath", fileInfo.FileName))
	return &container.ContainerCreateCreatedBody{}, nil
}
