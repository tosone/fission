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
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// specializeV2 specialize the pod
func specializeV2(context *gin.Context) {
	var err error
	if isSpecialized {
		context.String(http.StatusBadRequest, "not a generic container")
		return
	}

	if _, err = os.Stat(codePath); err != nil {
		if os.IsNotExist(err) {
			logger.Error("code path does not exist", zap.String("path", codePath))
			context.String(http.StatusInternalServerError, fmt.Sprintf("no such a file or directory: %s", codePath))
			return
		} else {
			logger.Error("unknown error", zap.Error(err), zap.String("path", codePath))
			context.String(http.StatusInternalServerError, fmt.Sprintf("unknown error: %v", err))
			return
		}
	}

	if err = load(); err != nil {
		logger.Error("subprocess with error", zap.Error(err))
		context.String(http.StatusInternalServerError, fmt.Sprintf("subprocess with error: %v", err))
	}
	context.String(http.StatusOK, "specialize is ok")
}