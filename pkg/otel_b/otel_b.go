package otel_b

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

type OtelB struct {
	Tracer trace.Tracer
	Client *http.Client
}

func (o *OtelB) InitTraceProvider(serviceName string, collectorUrl string) (func(context.Context) error, error) {
	ctx := context.Background()
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(serviceName),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	conn, err := grpc.DialContext(ctx, collectorUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %W", err)
	}

	traceExport, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))

	if err != nil {
		return nil, fmt.Errorf("failed to create export trace: %w", err)
	}

	// create span para envio em batch
	bsp := sdktrace.NewBatchSpanProcessor(traceExport)

	//create tracer provider with span bsp
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	o.Tracer = otel.Tracer("request service psi-mgnt")

	//insert automatic transport, any request with this client will be track
	//there are a generic interceptor
	o.Client = &http.Client{
		Transport: otelhttp.NewTransport(Interceptor{RoundTripper: http.DefaultTransport}),
	}

	return tracerProvider.Shutdown, nil
}

type Interceptor struct {
	http.RoundTripper
}

func (Interceptor) ModifyRequest(r *http.Request) *http.Request {
	// otel_b.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(r.Header))
	fmt.Println("LOG - Host: " + r.URL.Host)
	return r
}

func (i Interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	// modify before the request is sent
	newReq := i.ModifyRequest(r)

	// send the request using the DefaultTransport
	return i.RoundTripper.RoundTrip(newReq)
}

func (o *OtelB) TelemetryMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if o.Tracer != nil {
			var span trace.Span
			carrier := propagation.HeaderCarrier(r.Header)
			ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
			ctx, span = o.Tracer.Start(ctx, r.URL.Path)
			defer span.End()
		}
		*r = *r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	}
}

func (o *OtelB) WithRouteTag(route string, h http.HandlerFunc) http.Handler {
	return otelhttp.WithRouteTag(route, o.TelemetryMiddleware(h))
}
