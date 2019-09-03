package helper

import (
	"github.com/RichardKnop/machinery/v1/tracing"
	"github.com/opentracing/opentracing-go"
	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func SetupTracer(serviceName string) (func(), error) {

	tracing.MachineryTag = opentracing.Tag{Key: string(opentracing_ext.Component), Value: "NTESB"}
	tracing.WorkflowGroupTag = opentracing.Tag{Key: "NTESB.workflow", Value: "group"}
	tracing.WorkflowChordTag = opentracing.Tag{Key: "NTESB.workflow", Value: "chord"}
	tracing.WorkflowChainTag = opentracing.Tag{Key: "NTESB.workflow", Value: "chain"}
	// Jaeger setup code
	config := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
	}

	closer, err := config.InitGlobalTracer(serviceName)
	if err != nil {
		return nil, err
	}

	cleanupFunc := func() {
		closer.Close()
	}

	return cleanupFunc, nil
}
