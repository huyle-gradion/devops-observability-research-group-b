// Browser entry for the served page (app/static/index.html).
//
// Stage 11A keeps this deliberately minimal but real: it wires the "Run test"
// button to a live, same-origin GET /test request and renders the actual
// response. Stage 11C extends this same entry with the OpenTelemetry
// click/fetch trace and trace-id display; Stage 11D adds the real Web Vital
// measurement.
//
// esbuild compiles this file to app/static/app.js as an IIFE for ES2015
// browsers (see web/package.json). No framework and no runtime dependencies.
// Promise chains (not async/await) are used so later browser-context
// propagation does not depend on transient Zone async hooks.

(function () {
  "use strict";

  var TEST_URL = "/test?delay_ms=100";

  function byId(id) {
    return document.getElementById(id);
  }

  function setText(id, text) {
    var el = byId(id);
    if (el) {
      el.textContent = text;
    }
  }

  // A single click performs exactly one real request and shows the real status
  // and body. The trace-id and Web Vital areas stay untouched here; stages 11C
  // and 11D populate them from the actual trace and measurement.
  function runTest() {
    var button = byId("run-test");
    if (button) {
      button.disabled = true;
    }
    setText("response", "requesting " + TEST_URL + " ...");

    fetch(TEST_URL)
      .then(function (response) {
        return response.text().then(function (body) {
          return { status: response.status, body: body };
        });
      })
      .then(function (result) {
        setText("response", result.status + " " + result.body);
      })
      .catch(function (error) {
        setText("response", "request failed: " + error);
      })
      .then(function () {
        if (button) {
          button.disabled = false;
        }
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
