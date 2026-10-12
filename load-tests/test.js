// Stage 12 load test: a bounded, deterministic 9:1 success-to-controlled-error mix
// against an explicitly supplied own target.
//
// The traffic model is fixed and honest: nine real `GET /test?delay_ms=50` requests
// that must return 200, then one deliberate `GET /test?delay_ms=200&fail=1` request
// that must return 500. There is no random error injection, no fabricated backend
// timing and no token in any header.
//
// The app terminates the controlled-error body with a newline
// (`failBody = "intentional test failure\n"` in app/main.go), so the exact-body
// check matches that byte-for-byte.

import http from 'k6/http';
import { check } from 'k6';

// Own-target requirement: BASE_URL must be supplied explicitly. There is no default
// host. A missing value makes the run fail fast during initialisation, before any VU
// starts, so a load test can never silently target the wrong (or a public) service.
const BASE_URL = __ENV.BASE_URL;
if (!BASE_URL) {
  throw new Error(
    'BASE_URL is required and must be non-empty, e.g. ' +
      'BASE_URL=http://127.0.0.1:8080. Refusing to default to any other host.'
  );
}

// The deliberate 500 from /test?fail=1 is a controlled response, so mark it expected.
// Without this, k6 would count the intended error as a transport-level failure in
// http_req_failed; with it, http_req_failed only reflects unexpected statuses.
http.setResponseCallback(http.expectedStatuses(200, 500));

const SUCCESS_BODY_PREFIX = 'Hello world from container';
const CONTROLLED_ERROR_BODY = 'intentional test failure\n';

// One scenario keeps the nine-to-one ratio intact. `ramping-vus` ramps to 10 VUs over
// 30s, holds 10 VUs for 60s, then ramps back to 0 over 30s, with a bounded graceful
// stop so in-flight iterations finish instead of being killed.
export const options = {
  scenarios: {
    workload: {
      executor: 'ramping-vus',
      startVUs: 0,
      gracefulStop: '10s',
      stages: [
        { duration: '30s', target: 10 },
        { duration: '60s', target: 10 },
        { duration: '30s', target: 0 },
      ],
    },
  },
  thresholds: {
    // Every check in every iteration must pass.
    checks: ['rate==1'],
    // Only the success case is gated: unexpected failures there must stay near zero.
    // The deliberate ~10% of 500s is the controlled-error case and is NOT a failure,
    // so there is no global failure-rate threshold over the whole run.
    'http_req_failed{case:success}': ['rate<0.01'],
  },
};

export default function () {
  // Nine success requests.
  for (let i = 0; i < 9; i++) {
    const successRes = http.get(`${BASE_URL}/test?delay_ms=50`, {
      tags: { case: 'success' },
    });
    check(successRes, {
      'success: status is 200': (r) => r.status === 200,
      'success: body starts with Hello world from container': (r) =>
        typeof r.body === 'string' && r.body.startsWith(SUCCESS_BODY_PREFIX),
    });
  }

  // One controlled-error request.
  const errorRes = http.get(`${BASE_URL}/test?delay_ms=200&fail=1`, {
    tags: { case: 'controlled-error' },
  });
  check(errorRes, {
    'controlled-error: status is 500': (r) => r.status === 500,
    'controlled-error: body is intentional test failure': (r) =>
      r.body === CONTROLLED_ERROR_BODY,
  });
}
