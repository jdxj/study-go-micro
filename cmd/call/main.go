package main

import (
	"context"
	"io"
	"log"
	hello "study_go_micro/hello/proto"
	"time"

	wrapperTrace "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/asim/go-micro/v3"
)

func main() {
	t, io, err := NewTracer("tracer-cli", "127.0.0.1:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.WrapClient(
			wrapperTrace.NewClientWrapper(opentracing.GlobalTracer())),
	)

	service.Init()
	ms := hello.NewMeetService("hello", service.Client())
	req := &hello.User{Name: "abc"}
	go func() {
		for {
			time.Sleep(time.Second)
			resp, err := ms.Hello(context.TODO(), req)
			if err != nil {
				log.Fatalf("Hello: %s\n", err)
			}
			log.Printf("resp: %s\n", resp)
		}
	}()
	service.Run()
}

// NewTracer 创建一个jaeger Tracer
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
