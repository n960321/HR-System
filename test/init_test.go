package test

import (
	"HRSystem/internal/service"
	"HRSystem/pkg/database"
	"HRSystem/pkg/logger"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"
)

var (
	db               *database.Database
	accountSvc       *service.AccountService
	clockInRecordSvc *service.ClockInRecordService
	stopContainer    func()
)

func TestMain(m *testing.M) {
	logger.SetLogger(true)
	stopContainer = runMysqlImage()

	db = database.NewDatabase(database.Config{
		Host:         "0.0.0.0",
		Port:         "13306",
		User:         "root",
		Password:     "123456test",
		DBName:       "HR-System",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	})
	accountSvc = service.NewAccountService(db)
	clockInRecordSvc = service.NewClockInRecordService(db)

	SeedAccount(db.GetGorm())

	code := m.Run()
	if stopContainer != nil {
		stopContainer()
	}
	os.Exit(code)
}

func runMysqlImage() func() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	log.Info().Msg("Running Mysql image...")
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get current path")
	}
	log.Info().Str("currentPath", currentPath).Msg("currentPath")
	ctx := context.Background()
	containerName := fmt.Sprintf("HR-System-TEST-%v", time.Now().Unix())
	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: "mysql:latest",
			Env:   []string{"MYSQL_ROOT_PASSWORD=123456test"},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"3306/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "13306",
					},
				},
			},
			Binds: []string{currentPath + "/../deploy/mysql:/docker-entrypoint-initdb.d"},
		},
		nil,
		nil,
		containerName,
	)

	if err != nil {
		log.Fatal().Err(err).Str("container name", containerName).Msg("Create container Failed")
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Fatal().Err(err).Str("container name", containerName).Msg("Create start Failed")
	}
	// TODO - 要換個方式偵測mysql 健康狀況，是健康的才往下進行
	time.Sleep(10 * time.Second)

	log.Info().Str("container name", containerName).Msg("Starting Conatiner")


	return func() {
		cli.ContainerStop(ctx, resp.ID, container.StopOptions{})
		log.Info().Str("container name", containerName).Msg("Stop Conatiner")
	}
}
