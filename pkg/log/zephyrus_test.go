// Copyright © 2020 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package zaprus wraps uber-go/zap functions with a logrus style API.
package log

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	ctx := context.Background()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	newCtx := NewLoggerWithContext(ctx, zapLogger)

	// Fetching logger from Context
	logger := NewLoggerFromContext(newCtx)
	logger = logger.With(zap.Field{Key: "namespace", FieldType: "string", String: "csv"})
	logger.With(zap.Field{Key: "count", Integer: 1}).Info("Logging new value")
}
