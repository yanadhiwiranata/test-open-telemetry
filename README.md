Testing opentelemetry to signoz local


included lib that will used by other project

opentelemetry

```
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc

```


chi router
```
go get -u github.com/go-chi/chi/v5
go get github.com/riandyrn/otelchi

```


echo router
```
go get github.com/labstack/echo/v4
go get go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho
```