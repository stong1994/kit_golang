// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"fmt"
	"github.com/uber/jaeger-client-go"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/stong1994/kit_golang/log"
	jagerCfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// Init creates a new instance of Jaeger tracer.
func Init(serverName string, logFactory log.Factory, samplerConfig *jagerCfg.SamplerConfig, logSpans bool) (opentracing.Tracer, io.Closer, error) {
	cfg := jagerCfg.Configuration{
		ServiceName: serverName,
		Sampler:     samplerConfig,
		Reporter: &jagerCfg.ReporterConfig{
			QueueSize: 500,
			LogSpans:  logSpans,
		},
		Headers: &jaeger.HeadersConfig{
			TraceContextHeaderName:   "x-request-id",
			TraceBaggageHeaderPrefix: "trace-ctx",
		},
	}

	jaegerLogger := jaegerLoggerAdapter{logFactory.Bg()}

	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jagerCfg.Logger(jaegerLogger),
		jagerCfg.Metrics(jMetricsFactory),
		jagerCfg.MaxTagValueLength(2048),
	)
	if err != nil {
		logFactory.Bg().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	}
	return tracer, closer, err
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func (l jaegerLoggerAdapter) Debugf(msg string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(msg, args...))
}
