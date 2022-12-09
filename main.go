package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/labstack/echo/v4"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var tracer = otel.Tracer("echo-server")

func main() {
	serverChi()
	// serverEcho()
}

func initTracer() *sdktrace.TracerProvider {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("localhost:4317")),
	)
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "yan-adhi-chi"),
			attribute.String("library.language", "go"),
		),
	)

	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	return tp
}

func serverChi() {
	r := chi.NewRouter()

	r.Use(otelchi.Middleware("yan-adhi-chi"))
	chiRouter(r)

	tp := initTracer()
	defer func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	otel.SetTracerProvider(tp)
	http.ListenAndServe(":4000", r)
}

func serverEcho() {
	r := echo.New()

	tp := initTracer()
	defer func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	otel.SetTracerProvider(tp)
	r.Use(otelecho.Middleware("yan-adhi-chi"))
	echoRouter(r)

	http.ListenAndServe(":4000", r)
}

func chiRouter(m *chi.Mux) {
	m.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trytracer(r.Context())
		w.Write([]byte("index"))
	}))
}

func echoRouter(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		trytracer(c.Request().Context())
		return c.String(http.StatusOK, "index")
	})
}

func trytracer(ctx context.Context) {
	_, span := tracer.Start(context.Background(), "hello-1st-span")
	defer span.End()

	sleep(ctx)
}

func sleep(ctx context.Context) {
	time.Sleep(2 * time.Second)
	_, span := tracer.Start(context.Background(), "hello-2nd-span")
	defer span.End()
	time.Sleep(2 * time.Second)
}
