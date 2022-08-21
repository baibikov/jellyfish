/*
Copyright 2022 Jellyfish message broker
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"

	"jellifish/internal/broker"
	"jellifish/internal/config"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

const configPath = "./configs/broker.yaml"

func main() {
	if err := run(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func run() (err error) {
	logrus.Info("init jellyfish ctx")
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer cancel()

	logrus.Info("init and parse jellyfish config")
	cnf, err := config.New(configPath)
	if err != nil {
		return errors.Wrapf(err, "jellyfish: parse config by path %s", configPath)
	}

	logrus.Infof("init jellyfish broker at addr %s", cnf.Addr)
	listener, err := broker.New(cnf)
	if err != nil {
		return errors.Wrapf(err, "jellyfish: run broker on host %s", cnf.Addr)
	}
	defer multierr.AppendInvoke(&err, multierr.Close(listener))
	go func() {
		logrus.Info("start jellyfish broadcast")
		multierr.AppendInto(&err, errors.Wrap(listener.Broadcast(ctx), "jellyfish"))
	}()

	<-ctx.Done()
	logrus.Info("stop jellyfish broker app")
	return nil
}
