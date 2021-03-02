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
	"log"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var app *gin.Engine

// running subprocess is running or not
var running bool

// isSpecialized each pod should specialize once
var isSpecialized bool

const codePath = "/userfunc/user"

func main() {
	var err error

	var config = zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if logger, err = config.Build(); err != nil {
		log.Fatalf("can't initialize zap logger: %v\n", err)
	}

	defer func() {
		if err = logger.Sync(); err != nil {
			log.Fatalf("flush log with error: %v\n", err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	app = gin.Default()

	app.Use(gin.Recovery())
	app.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	app.Use(ginzap.RecoveryWithZap(logger, true))

	app.GET("/healthz", healthz)
	app.GET("/specialize", specialize)
	app.GET("/v2/specialize", specializeV2)

	app.NoRoute(func(context *gin.Context) {
		context.Status(http.StatusNotFound)
	})
	if err = app.Run(":8888"); err != nil {
		logger.Error("server start with error", zap.Error(err))
	}
}
