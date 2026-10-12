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
import { onCLS, onINP, onLCP } from "web-vitals";

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

  // ---------------------------------------------------------------------------
  // 11D: real Core Web Vitals from THIS page (EXPERIMENTAL, controlled demo)
  //
  // Where the numbers come from:
  //   * Navigation Timing  -> one PerformanceNavigationTiming entry describes
  //     this document load (start, DOMContentLoaded, load).
  //   * Resource Timing    -> a PerformanceResourceTiming entry per fetched
  //     resource. The same-origin /test response sets a Server-Timing header,
  //     which is visible in DevTools Network (Timing tab) and readable here as
  //     PerformanceResourceTiming.serverTiming (see reportServerTiming below).
  //   * PerformanceObserver-> the browser pushes largest-contentful-paint,
  //     layout-shift and event-timing entries as they happen; web-vitals
  //     subscribes to them and calls these callbacks with the REAL value and a
  //     measurement id. Nothing below is invented or copied from another site.
  //
  // Lab vs field: the DevTools Performance panel and Trace Viewer replay one
  // reproducible session (lab). Real-user monitoring aggregates many real
  // users' sessions (field). This page is a single controlled local/lab demo.
  //
  // EXPERIMENTAL / limits: web-vitals 6.2.3 and the underlying browser
  // PerformanceObserver APIs are pre-stable and vary by browser. This is a
  // controlled single-session demo, NOT a production RUM cohort, and it is NOT
  // a CrUX statistic (CrUX is aggregated field data over many real users behind
  // a public origin). A CLS of 0 is a real measured value for a load with no
  // layout shift, not a placeholder.
  //
  // Units: LCP and INP are DOMHighResTimeStamp millisecond durations; CLS is a
  // unitless layout-shift score.
  var VITAL_UNITS = { LCP: "ms", INP: "ms", CLS: "score (unitless)" };
  var VITAL_ORDER = ["LCP", "INP", "CLS"];
  var vitalState = {};

  function browserContext() {
    var ua = (navigator && navigator.userAgent) || "unknown user agent";
    return ua + " | controlled local demo, not field RUM, not CrUX";
  }

  function captureTime() {
    return new Date().toISOString();
  }

  // Called by web-vitals with the real metric object. Displays AND logs the
  // real name, value, measurement id, unit, page URL, context and capture time.
  function reportVital(metric) {
    var unit = VITAL_UNITS[metric.name] || "unit";
    var record = {
      name: metric.name,
      value: metric.value,
      unit: unit,
      id: metric.id,
      rating: metric.rating,
      url: window.location.href,
      context: browserContext(),
      capturedAt: captureTime(),
    };
    vitalState[metric.name] = record;
    if (window.console && window.console.log) {
      window.console.log(
        "[web-vitals] " + record.name + " = " + record.value + " " + unit,
        record
      );
    }
    renderVitals();
  }

  function renderVitals() {
    var el = byId("web-vital");
    if (!el) {
      return;
    }
    var blocks = [];
    for (var i = 0; i < VITAL_ORDER.length; i++) {
      var name = VITAL_ORDER[i];
      var m = vitalState[name];
      if (!m) {
        blocks.push(name + ": no real report yet (waiting for the browser)");
        continue;
      }
      blocks.push(
        m.name + ": " + m.value + " " + m.unit + " (rating: " + m.rating + ")\n" +
          "  id: " + m.id + "\n" +
          "  url: " + m.url + "\n" +
          "  context: " + m.context + "\n" +
          "  captured: " + m.capturedAt
      );
    }
    el.textContent = blocks.join("\n\n");
  }

  // Reads the same-origin /test Server-Timing header back through Resource
  // Timing. Same-origin, so serverTiming is exposed (no Timing-Allow-Origin
  // needed). Never fabricates a value: if the entry/header is absent it says so.
  function reportServerTiming(requestUrl) {
    var note = "no PerformanceResourceTiming entry yet for " + requestUrl +
      " (see DevTools Network > Timing).";
    if (window.performance && performance.getEntriesByName) {
      var entries = performance.getEntriesByName(requestUrl, "resource");
      if (entries.length) {
        var entry = entries[entries.length - 1];
        var serverTiming = entry.serverTiming;
        if (serverTiming && serverTiming.length) {
          var parts = [];
          for (var i = 0; i < serverTiming.length; i++) {
            parts.push(serverTiming[i].name + " dur=" + serverTiming[i].duration + "ms");
          }
          note = "PerformanceResourceTiming.serverTiming (" + requestUrl + "): " +
            parts.join(", ");
        } else {
          note = "PerformanceResourceTiming for " + requestUrl +
            " has no Server-Timing entry (see DevTools Network).";
        }
      }
    }
    setText("server-timing", note);
    if (window.console && window.console.log) {
      window.console.log("[Server-Timing] " + note);
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
          // The Resource Timing entry for /test is finalized with the response,
          // so read it just after this task settles (clearer than racing the
          // promise microtasks).
          window.setTimeout(function () {
            reportServerTiming(new URL(TEST_URL, window.location.href).href);
          }, 0);
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

  // Registers all three real Core Web Vitals with reportAllChanges so the page
  // shows/logs every real update (including a genuine CLS of 0). LCP and INP
  // are millisecond durations; CLS is a unitless score.
  function registerWebVitals() {
    var options = { reportAllChanges: true };
    onLCP(reportVital, options);
    onINP(reportVital, options);
    onCLS(reportVital, options);
  }

  function init() {
    var button = byId("run-test");
    if (button) {
      button.addEventListener("click", runTest);
    }
    renderVitals();
    registerWebVitals();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
