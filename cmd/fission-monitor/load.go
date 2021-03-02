/*
Copyright 2021 The Fission Authors.

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
	"io"
	"os"
	"os/exec"
	"path"

	"go.uber.org/zap"
)

func load() (err error) {
	var cmd = exec.Command(codePath)

	cmd.Env = []string{}
	cmd.Dir = path.Dir(codePath)

	var stdoutIn, stderrIn io.ReadCloser
	if stdoutIn, err = cmd.StdoutPipe(); err != nil {
		return
	}
	if stderrIn, err = cmd.StderrPipe(); err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		logger.Error("subprocess start with error", zap.Error(err))
	}

	running = true

	go func() {
		if err = cmd.Wait(); err != nil {
			logger.Error("subprocess running with error", zap.Error(err))
		}
		running = false
	}()

	go func() {
		if _, err = io.Copy(os.Stdout, stdoutIn); err != nil {
			logger.Error("copy stdout with err", zap.Error(err))
		}
	}()

	go func() {
		if _, err = io.Copy(os.Stderr, stderrIn); err != nil {
			logger.Error("copy stdout with err", zap.Error(err))
		}
	}()

	return
}
