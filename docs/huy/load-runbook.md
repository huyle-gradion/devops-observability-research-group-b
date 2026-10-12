# Load runbook — Stage 12 k6 run against the private app forward (huy)

Repeatable procedure for the **controlled** k6 run in Stage 12: a bounded 9:1
success-to-controlled-error mix sent through a private `kubectl port-forward` to our
own `devops-app`, observed on the `huy-workload` Grafana dashboard at the same time.

**Scope of this run.** It exercises our own private forwarded service, proves the
script's checks/thresholds pass against real traffic, and shows the dashboard move.
It is the Step 4/5 rehearsal evidence. It **does not certify the Friday Oct 16 live
demo** — that is a separate, later gate performed live. Do not label a recording of
this run as the Friday demo.

**Target rule.** The load goes only to `BASE_URL`, which you must set to the private
forward. The script has no default host and fails fast if `BASE_URL` is empty. Never
point it at a public host or another team's service.

---

## 0. Preconditions

- `KUBECONFIG="$HOME/.kube/huy-k3s.yaml"` and a reachable `huy-k3s` cluster.
- The app is deployed in namespace `devops-app` (`svc/devops-app`, container port 8080).
- The `huy-monitoring` kube-prometheus-stack release is installed and the ServiceMonitor
  applied (`k8s/monitoring/servicemonitor.yaml`), so the app target is scraped.
- The `huy-workload` dashboard is imported and saved in Grafana.

```sh
export KUBECONFIG="$HOME/.kube/huy-k3s.yaml"
kubectl -n devops-app get pods -o wide
kubectl -n devops-app get svc devops-app
```

## 1. Pinned k6 binary (record the version)

Use the pinned **official** release binary for this workstation (native linux/amd64),
not a rolling `latest` and not a package manager that drifts. Pinned version:
**`v2.3.0`**.

```sh
K6_VERSION=v2.3.0
mkdir -p ~/.local/bin
cd "$(mktemp -d)"

curl -fsSLO "https://github.com/grafana/k6/releases/download/${K6_VERSION}/k6-${K6_VERSION}-linux-amd64.tar.gz"
curl -fsSLO "https://github.com/grafana/k6/releases/download/${K6_VERSION}/k6-${K6_VERSION}-checksums.txt"

# Verify the tarball against the release-published checksum before trusting it.
grep "linux-amd64.tar.gz" "k6-${K6_VERSION}-checksums.txt" | sha256sum -c -

tar -xzf "k6-${K6_VERSION}-linux-amd64.tar.gz"
install -m 0755 "k6-${K6_VERSION}-linux-amd64/k6" ~/.local/bin/k6
```

Expected checksum for `k6-v2.3.0-linux-amd64.tar.gz` (from
`k6-v2.3.0-checksums.txt`):

```
39c3117b6af817592dcd0ce4242105c0a7af10948c2a425306f0be8f7a8a8ab1
```

On this workstation k6 was already present at `~/.local/bin/k6`; if it is, skip the
download and just record what is installed:

```sh
command -v k6
k6 version          # record the exact version string in the run notes
```

## 2. Start the private app forward and set BASE_URL

Forward the app Service to localhost (pick a free port; 8080 is the default used by the
rest of the runbooks).

```sh
kubectl -n devops-app port-forward svc/devops-app 8080:8080
```

Leave that terminal running. In a second terminal, set `BASE_URL` and prove the forward
answers before generating load:

```sh
export BASE_URL=http://127.0.0.1:8080
curl -s http://127.0.0.1:8080/test        # expect: Hello world from container-<pod>\n
```

Note that the Service forward pins one backing pod; that single pod will carry the
traffic while the other stays quieter. That is expected and is why both pods are shown
on the resource panels.

## 3. Confirm the Prometheus target is Up (before you run)

Do not start load until the scrape target is healthy; otherwise the dashboard is empty
for reasons unrelated to the test. The ServiceMonitor scrapes every **15s**.

```sh
kubectl -n monitoring get servicemonitor devops-app
kubectl -n monitoring get svc          # read the actual Prometheus Service name first
# The chart exposes an Operator-created `prometheus-operated` Service (port 9090):
kubectl -n monitoring port-forward svc/prometheus-operated 9090:9090
#   then browse http://127.0.0.1:9090/targets and find the devops-app target = Up
```

Also confirm the app's own metrics endpoint is scraped:

```sh
curl -s http://127.0.0.1:8080/metrics | grep '^http_requests_total'
```

## 4. Open the Grafana `huy-workload` dashboard

```sh
kubectl -n monitoring get svc                       # read the actual Grafana Service name
kubectl -n monitoring port-forward svc/huy-monitoring-grafana 3000:80
```

Browse <http://127.0.0.1:3000>, open the **`huy-workload`** dashboard, and set the time
range to **Last 15 minutes** with the built-in **5s refresh**. Its five panels are:

- **Request rate — GET /test (req/s)**
- **Error ratio — GET /test (percentunit)**
- **p95 latency — GET /test (s)**
- **Per-pod CPU — app container (cores)**
- **Per-pod memory — app container (bytes)**

Leave this visible **while** k6 runs (step 5), so you watch k6 output and the dashboard
move at the same time.

## 5. Run the load test

From the repository root, with `BASE_URL` exported (the script reads `__ENV.BASE_URL`):

```sh
k6 run load-tests/test.js
```

Equivalent explicit form if the environment does not propagate:

```sh
k6 run -e BASE_URL="$BASE_URL" load-tests/test.js
```

The single `ramping-vus` scenario runs 30s ramp to 10 VUs, 60s hold, 30s ramp down, with
a 10s `gracefulStop`. Each iteration sends **nine** `GET /test?delay_ms=50` (expect 200,
body starts `Hello world from container`) then **one** `GET /test?delay_ms=200&fail=1`
(expect 500, exact body `intentional test failure\n`).

## 6. Watch k6 output and the dashboard together

While it runs, keep both in view:

- **k6 stdout**: `checks`, `http_req_duration`/p95, and the case-tagged request/failure
  rates. `http_req_failed{case:success}` and `checks` are evaluated at the end.
- **Grafana `huy-workload`**: request rate and error ratio should rise during ramp/steady
  and fall during cooldown, with the p95 reflecting the added `delay_ms`. The two
  per-pod CPU/memory panels show the sampled resource use of the received pod.

Remember the **scrape/rate lag**: panels use a **1m rate window** over **15s scrapes**,
so a panel can take a minute or two to reflect traffic (and to settle after it stops).
That delay is normal, not a broken dashboard.

## 7. Stop the forward

After the run finishes and you have read the summary and the paired dashboard screenshots,
stop the app forward with `Ctrl-C` in its terminal (and the Grafana/Prometheus forwards
when done). Stopping the forward is enough; the cluster, Services and data stay up.

---

## Expected results / interpretation

- **Mix:** the workload is a real 9:1 — nine success requests to one controlled-error
  request per iteration — so the app's error ratio should settle near ~10% of `/test`
  traffic over the steady interval (not exactly during ramp/window transitions).
- **Thresholds:** `checks: ['rate==1']` (every check passes) and
  `http_req_failed{case:success}: ['rate<0.01']` (unexpected success-case failures near
  zero). There is deliberately **no** global failure-rate gate, because ~10% of the
  requests are intended 500s.
- **The 500s are the controlled-error case**, not an outage: `http.setResponseCallback(
  http.expectedStatuses(200, 500))` marks them expected, and they carry `case:
  controlled-error`. Application counters still correctly classify them as app errors.
- **One pod busy, one quiet:** a Service forward targets one pod; the resource panels
  show both honestly rather than inventing a spike.
- If a threshold fails, diagnose the real cause — status/body mismatch, a dead/rolled
  forward, a scrape/target problem, or genuine resource pressure — before re-running.

## Troubleshooting

**Port-forward lifecycle.**
- *Port already in use* (`unable to listen on port 8080`): something already holds 8080.
  Find it and either stop it or forward a different local port, then set `BASE_URL` to
  match (e.g. `8081:8080` → `BASE_URL=http://127.0.0.1:8081`). The container port stays
  8080.
- *Forward dropped mid-run* (`error: lost connection to pod`, `Connection refused`,
  `EOF`): if the pod rolled or the connection dropped, k6 will start seeing transport
  errors. Stop k6, recreate the forward, re-check step 3, and re-run — do **not** count
  transport errors from an abandoned forward as a valid demonstration.
- *`curl: (7) Failed to connect`* before the run: the forward is not up yet; wait for
  `Forwarding from 127.0.0.1:8080 -> 8080` and retry.

**Prometheus target.**
- If the target is not `Up`, the dashboard stays empty. Check the ServiceMonitor labels
  (`release: huy-monitoring`, namespace selector) and that the app Service carries
  `app: devops-app` with the `http` port the ServiceMonitor references.
- Confirm `/metrics` actually serves `http_requests_total` via the forward (step 3).

**Reading the panels.**
- The **request-rate / error-ratio / p95** panels come from the app's own `/metrics`
  (`http_requests_total`, `http_request_duration_seconds`) — they answer "how much
  traffic and how is it behaving".
- The **per-pod CPU / memory** panels come from kubelet/cAdvisor — they answer "what did
  that traffic cost", not "is traffic flowing". Do not read resource panels as request
  counts, or request panels as resource usage.
- An empty panel in the first minute is scrape/rate-window lag, not a failure. Wait for
  at least one full 15s scrape plus a full 1m rate window after traffic.

**Expected 500s look like errors.**
- The `intentional test failure` 500s are the **controlled-error case**. k6 marks them
  expected (they do not inflate `http_req_failed`) and the app counts them as app errors,
  which is the point. If you see a 500 whose body is **not** exactly
  `intentional test failure` (plus newline), or a 200 on the error request, that is a real
  problem — investigate, do not relabel it.

## Record for the handoff

- k6 version string (`k6 version`).
- `BASE_URL` used and the exact forward command.
- k6 summary: `checks` rate, `http_req_failed{case:success}` rate, request counts by
  `case`, p95.
- Dashboard screenshots at ramp, steady and cooldown (same time range), noting the
  received pod vs. the quieter one.
- Run timestamp and exit status. State clearly that this run is the Step 4/5 rehearsal
  and does **not** certify the Friday live demo.
