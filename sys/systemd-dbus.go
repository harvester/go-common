package sys

import (
	"context"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/sirupsen/logrus"
)

// RestartService restarts a service.  If the service isn't already running it will be started.
func RestartService(unit string) error {
	ctx := context.Background()
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		logrus.Errorf("Failed to create new connection for systemd. err: %v", err)
		return err
	}
	defer conn.Close()
	responseChan := make(chan string, 1)
	if _, err := conn.RestartUnitContext(ctx, unit, "fail", responseChan); err != nil {
		logrus.Errorf("Failed to restart service %s. err: %v", unit, err)
		return err
	}
	return nil
}

// TryRestartService will restart a service, but only if it's currently running.
// A service that isn't running won't be affected.
func TryRestartService(unit string) error {
	ctx := context.Background()
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		logrus.Errorf("Failed to create new connection for systemd. err: %v", err)
		return err
	}
	defer conn.Close()
	responseChan := make(chan string, 1)
	if _, err := conn.TryRestartUnitContext(ctx, unit, "fail", responseChan); err != nil {
		logrus.Errorf("Failed to restart service %s. err: %v", unit, err)
		return err
	}
	return nil
}
