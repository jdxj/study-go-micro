package main

import (
	"io"
	"log"
	"time"

	"github.com/uber/jaeger-client-go"

	"study_go_micro/hello/handler"
	proto "study_go_micro/hello/proto"

	"github.com/asim/go-micro/v3"
	"github.com/opentracing/opentracing-go"

	wrapperTrace "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	t, io, err := NewTracer("tracer-srv", "127.0.0.1:6831")
	if err != nil {
		log.Fatalf("NewTracer: %s\n", err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name("hello"),
		micro.WrapHandler(wrapperTrace.NewHandlerWrapper(opentracing.GlobalTracer())))
	service.Init()

	err = proto.RegisterMeetHandler(service.Server(), &handler.Helloer{})
	if err != nil {
		log.Fatalf("RegisterMeetHandler: %s\n", err)
	}

	if err := service.Run(); err != nil {
		log.Printf("Run: %s\n", err)
	}
}

func NewTracer(servicename string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	return tracer, closer, err
}
