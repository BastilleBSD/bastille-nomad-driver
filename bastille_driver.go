package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/hashicorp/nomad/plugins/drivers/base"
	"github.com/hashicorp/go-hclog"
)

type BastilleDriver struct {
	config *base.Config
	logger hclog.Logger
}

func NewBastilleDriver() (*BastilleDriver, error) {
	return &BastilleDriver{
		logger: hclog.New(&hclog.LoggerOptions{
			Name:  "bastille",
			Level: hclog.Debug,
		}),
	}, nil
}

func (d *BastilleDriver) Name() string {
	return "bastille"
}

func (d *BastilleDriver) Type() drivers.Type {
	return drivers.Type("bastille")
}

func (d *BastilleDriver) Capabilities() *drivers.Capabilities {
	return &drivers.Capabilities{TaskSignals: true}
}

func (d *BastilleDriver) StartTask(ctx context.Context, req *drivers.StartTaskRequest, stream drivers.TaskHandle) (*drivers.StartTaskResponse, string, error) {
	jailName := req.TaskName
	release := req.Config["release"]
	ip := req.Config["ip"]
	iface := req.Config["iface"]

	args := []string{"create", jailName, release, ip}
	if iface != "" {
		args = append(args, iface)
	}

	if err := runBastille(args...); err != nil {
		return nil, "", fmt.Errorf("failed to create jail: %w", err)
	}

	if err := runBastille("start", jailName); err != nil {
		return nil, "", fmt.Errorf("failed to start jail: %w", err)
	}

	return &drivers.StartTaskResponse{
		TaskID: jailName,
	}, "", nil
}

func (d *BastilleDriver) StopTask(ctx context.Context, req *drivers.StopTaskRequest) (*drivers.StopTaskResponse, error) {
	if err := runBastille("stop", req.TaskID); err != nil {
		return nil, err
	}
	if err := runBastille("destroy", "-y", req.TaskID); err != nil {
		return nil, err
	}
	return &drivers.StopTaskResponse{}, nil
}

func (d *BastilleDriver) RecoverTask(ctx context.Context, req *drivers.RecoverTaskRequest) (*drivers.RecoverTaskResponse, error) {
	// Assume jail is still running
	return &drivers.RecoverTaskResponse{
		TaskID: req.TaskID,
	}, nil
}

