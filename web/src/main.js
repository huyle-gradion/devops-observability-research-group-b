// Browser entry for the served page (app/static/index.html).
//
// Stage 11C wires the OpenTelemetry Web SDK into this page so a single click on
// "Run test" produces ONE real trace that is joined to the backend:
//
//   test-click            (this page, internal span)
//     -> HTTP GET /test   (fetch client span, injected with W3C traceparent)
//        -> GET /test     (backend otelhttp server span)
//           -> test-work  (backend child span; logged with the same trace_id)
//
// EXPERIMENTAL: OpenTelemetry browser instrumentation is pre-stable. Its
// maturity and API guarantees are NOT the same as the server SDKs, and this
// controlled single-session page is not a production RUM cohort. The page
// states that limitation in the UI.
//
// Same-origin export: BatchSpanProcessor posts spans to
// window.location.origin + '/otel/v1/traces', which stage 11B proxies to the
// private Collector. The exact export URL is listed in the fetch and XHR
// instrumentations' ignoreUrls so export requests never create recursive spans.
//
// The click flow uses Promise chains (not async/await) so the active browser
// context is preserved by ZoneContextManager without depending on native
// async/await Zone preservation.
//
// esbuild compiles this file to app/static/app.js as an IIFE for ES2015
// browsers (see web/package.json). No framework.

import {
  context,
  SpanKind,
  SpanStatusCode,
  trace,
} from "@opentelemetry/api";
import { resourceFromAttributes } from "@opentelemetry/resources";
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { WebTracerProvider } from "@opentelemetry/sdk-trace-web";
import { ZoneContextManager } from "@opentelemetry/context-zone";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";
import { registerInstrumentations } from "@opentelemetry/instrumentation";
import { DocumentLoadInstrumentation } from "@opentelemetry/instrumentation-document-load";
import { FetchInstrumentation } from "@opentelemetry/instrumentation-fetch";
import { XMLHttpRequestInstrumentation } from "@opentelemetry/instrumentation-xml-http-request";

(function () {
  "use strict";

  var TEST_URL = "/test?delay_ms=100";
  var EXPORT_PATH = "/otel/v1/traces";
  var SERVICE_NAME = "devops-web-huy";
  var TRACER_NAME = "devops-web-huy-page";

  // Same-origin export endpoint. The OTLP HTTP exporter posts the span batch to
  // exactly this URL, and the same value is the ignoreUrls entry below.
  var exportUrl = window.location.origin + EXPORT_PATH;

  var exporter = new OTLPTraceExporter({ url: exportUrl });

  var provider = new WebTracerProvider({
    resource: resourceFromAttributes({ "service.name": SERVICE_NAME }),
    spanProcessors: [new BatchSpanProcessor(exporter)],
  });

  // register() with no explicit propagator installs the default W3C composite
  // propagator (W3CTraceContextPropagator + W3CBaggagePropagator), so the fetch
  // and XHR instrumentations below inject a real `traceparent` header
  // automatically. This page never builds the header or any id by hand.
  provider.register({ contextManager: new ZoneContextManager() });

  registerInstrumentations({
    tracerProvider: provider,
    instrumentations: [
      new DocumentLoadInstrumentation(),
      // ignoreUrls matches an exact string URL, so the OTLP export request to
      // exportUrl is never traced (an instrumented export would recurse).
      new FetchInstrumentation({ ignoreUrls: [exportUrl] }),
      new XMLHttpRequestInstrumentation({ ignoreUrls: [exportUrl] }),
    ],
  });

  var tracer = provider.getTracer(TRACER_NAME);

  function byId(id) {
    return document.getElementById(id);
  }

  function setText(id, text) {
    var el = byId(id);
    if (el) {
      el.textContent = text;
    }
  }

  // One click -> one active `test-click` span wrapping one real request. The
  // reported status/trace-id come only from the actual response; the span is
  // ended on both success and error, then the provider is flushed.
  function runTest() {
    var button = byId("run-test");
    if (button) {
      button.disabled = true;
    }
    setText("response", "requesting " + TEST_URL + " ...");
    setText("trace-id", "pending ...");

    var span = tracer.startSpan("test-click", { kind: SpanKind.INTERNAL });
    span.setAttribute("http.request.method", "GET");
    var traceId = span.spanContext().traceId;

    // Make the click span active for the request so the fetch instrumentation
    // parents its client span under test-click.
    var ctx = trace.setSpan(context.active(), span);

    context.with(ctx, function () {
      fetch(TEST_URL)
        .then(function (response) {
          span.setAttribute("http.response.status_code", response.status);
          return response.text().then(function (body) {
            return { status: response.status, body: body };
          });
        })
        .then(function (result) {
          // Success is shown only after the real response body has arrived.
          setText("response", result.status + " " + result.body);
          setText("trace-id", traceId);
          if (result.status >= 400) {
            span.setStatus({
              code: SpanStatusCode.ERROR,
              message: "HTTP " + result.status,
            });
          } else {
            span.setStatus({ code: SpanStatusCode.OK });
          }
        })
        .catch(function (error) {
          // Network/exception path: show the real failure and record it.
          setText("response", "request failed: " + error);
          setText("trace-id", traceId);
          span.recordException(error);
          span.setStatus({ code: SpanStatusCode.ERROR, message: String(error) });
        })
        .then(function () {
          // End the click span exactly once on BOTH paths, then flush so the
          // batch actually leaves the browser for the same-origin proxy.
          span.end();
          if (button) {
            button.disabled = false;
          }
          provider.forceFlush().catch(function (flushError) {
            // A flush failure only means spans could not reach the Collector; it
            // is not a request failure, so it is logged, never shown as success.
            if (window.console) {
              window.console.error("OTel flush failed", flushError);
            }
          });
        });
    });
  }

  function init() {
    var button = byId("run-test");
    if (button) {
      button.addEventListener("click", runTest);
    }
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
