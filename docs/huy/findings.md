# Research findings (Assignment 3.1a–3.1f) — personal project

Owner: huy. Created: 2026-10-11 (Stage 01, substep 01C). Registry compiled: 2026-10-13 (Stage 13A).

This file is the **claim registry** for research topics 3.1a–3.1f. Each row records one
claim in the owner's words, tied to a **primary page that was actually checked** (with a
section anchor where one exists), the date it was checked, the version/status it depends on,
what the check confirmed or corrected, and the evidence from this project's own system where
the claim is about something we ran.

## How to read this registry

| Column | Meaning |
|---|---|
| Claim (owner words) | The finding stated plainly, not as vendor marketing. |
| Primary source | The exact documentation/spec page checked, never a landing page alone. |
| Checked | The date the page was checked. All registry rows: **2026-10-13**. |
| Version/status | The version, GA/beta/experimental label or dated changelog entry the page shows. |
| Confirmation / correction / uncertainty | What the check confirmed, what it corrected, or what remains uncertain. |
| Own-system evidence | The artifact in this repo that shows the claim applied to our system (`—` when the claim is general). |

Rules kept while writing this:

- A literal claim from the assignment brief is treated as a **verification checklist**, not
  something to repeat. Where the brief's wording is wrong or overstated, the row says so.
- No invented "stable since" dates, no quoted third-party endorsements, no undocumented vendor
  guarantees. Uncertain claims are marked **uncertain** in the row.
- Moving claims (Immutable Actions rollout; OpenTelemetry signal maturity) carry explicit
  **16 Oct 2026** and **30 Oct 2026** recheck entries, both `pending`, in the
  "Moving-claim rechecks" section below.
- Rows are written by an agent and do **not** certify the owner's understanding; the owner's
  reproduction and explanation is the learning gate (spec §"Agent and learner gates").

> **Compilation status.** Compiled from bounded primary-source checks on 2026-10-13. Topics
> whose source check had not returned at write time are marked `pending` with the checking
> lane id and are not filled by guesswork.

## a) CI/CD fundamentals & one real system

pending — source check in progress (lane `ses_ece6e8658ffeOzMdDEIk9tGfmq`). Scope: pipeline
definition and standard stages and what each gates; GitHub Actions workflow YAML anatomy
(triggers, jobs, steps, runners, `needs`, events); CI vs continuous delivery vs deployment;
SLSA/in-toto/Sigstore and real Artifact Attestations; Immutable Actions rollout and the
OCI-on-GHCR claim, checked independently; short GitLab CI and Jenkins comparison.

| Claim (owner words) | Primary source | Checked | Version/status | Confirmation / correction / uncertainty | Own-system evidence |
|---|---|---|---|---|---|
| _pending — awaiting lane result; no row invented._ | — | — | — | — | — |

## b) GitOps & continuous delivery

Scope: push vs pull/GitOps models; Argo CD and Flux as the standard agents (one deep);
git as single source of truth; reviewable/rollbackable desired state; drift
detection/correction; manual vs automatic sync; how the image tag is updated after a CI
build; cloudflared/ngrok/Tailscale for reaching a service from outside and which fits a
public demo vs a team-internal service (free tiers).

| Claim (owner words) | Primary source | Checked | Version/status | Confirmation / correction / uncertainty | Own-system evidence |
|---|---|---|---|---|---|
| GitOps is defined by four principles — declarative, versioned/immutable, **pulled** automatically, continuously reconciled — not by a tool. The pull model is the point: an agent in the cluster pulls desired state from git, so the CI job never needs cluster credentials. | https://opengitops.dev/ (four GitOps principles) | 2026-10-13 | GitOps Principles v1.0.0 | Confirms the brief: "push" (CI holds cluster creds and runs `kubectl apply`) is the anti-pattern; pull is a principle, not a vendor choice. Uncertainty: opengitops.dev is the project's home page — the four principles are the page's own content, not a deeper doc. | CI `propose-deploy` holds no cluster credential (no kubeconfig/SSH/SA token) and never writes `env/huy`; Argo CD inside the cluster pulls `env/huy`. `docs/huy/gitops-runbook.md`, `.github/workflows/ci.yml`. |
| Argo CD is a **pull-based reconciler**: with automated sync the CI pipeline only commits to git and the in-cluster Application controller polls and applies. It does not let CI "push" to the cluster. | https://argo-cd.readthedocs.io/en/stable/user-guide/auto_sync/ | 2026-10-13 | GA (stable docs) | Corrects the common "Argo CD pushes changes" wording. The controller reconciles at `timeout.reconciliation` (default 120s + jitter). | `k8s/gitops/application.yaml`; adoption observed `OutOfSync/Healthy` → manual sync → `Synced/Healthy`. `captures/huy/bonus-1/adoption.md`. |
| An Argo CD **Application** maps a git source (`repoURL`, `targetRevision`, `path`) to a destination cluster+namespace, and `syncPolicy.automated` carries `enabled`, `prune`, `selfHeal`, `allowEmpty`. Reviewed desired state lives in the git repo, not in the cluster. | https://argo-cd.readthedocs.io/en/stable/user-guide/application-specification/ | 2026-10-13 | GA; `allowEmpty` added in v1.8+ | Confirms `env/huy` is the single source of truth for the app. Note: `automated.enabled: null` is treated as enabled (backward compatible). | Application `devops-app-huy` := repo `…research-group-b`, `targetRevision: env/huy`, `path: k8s/app`, namespace `devops-app`; `syncPolicy.automated` sets `enabled/selfHeal/prune: true`, `allowEmpty: false`. `k8s/gitops/application.yaml`. |
| **Drift** is expected, not a bug: an app can be `OutOfSync` right after a sync because controllers/webhooks mutate fields or Helm/Kustomize render non-deterministically. Drift is corrected by ignoring differences at JSON-pointer/JQ paths or via `managedFieldsManagers` — or by re-syncing with self-heal. | https://argo-cd.readthedocs.io/en/stable/user-guide/diffing/ | 2026-10-13 | GA (stable) | Corrects "drift means the tool is broken". The `status` field is a frequent source of perpetual OutOfSync and is commonly ignored. | Drift procedure: `kubectl scale` by hand while Git says `2`; `selfHeal` restores `2`. `docs/huy/gitops-runbook.md`. |
| Sync can be **manual** (a human clicks Sync / merges the reviewed change) or **automatic** (`syncPolicy.automated`, with self-heal and prune). **Rollback cannot be performed while automated sync is enabled** — that is a hard constraint, not a workaround. | https://argo-cd.readthedocs.io/en/stable/user-guide/auto_sync/#rollback | 2026-10-13 | GA; `allowEmpty` v1.8+ | Corrects the assumption that you can `rollout undo` an auto-synced app: the controller would restore Git's value. Roll back by reverting the commit in git. | Manual sync proved the loop (08C), then automated sync was enabled (08D). Rollback is documented as a git revert, never `kubectl rollout undo`. `captures/huy/bonus-1/adoption.md`, `docs/huy/gitops-runbook.md`. |
| Argo CD's `AppProject` scopes the allowed `sourceRepos`, `destinations` and (optionally) sync windows; Applications/AppProjects are ordinary Kubernetes CRDs applied declaratively (`kubectl apply`), which does not bypass Argo's reconciliation. | https://argo-cd.readthedocs.io/en/stable/operator-manual/declarative-setup/ | 2026-10-13 | GA (stable) | Confirms scoping is the AppProject's job. Warning from the page: a project deploying into the Argo CD namespace grants admin — ours does not. | AppProject `huy`: sourceRepo = the project repo, destination = in-cluster `devops-app`, `clusterResourceWhitelist: []`. `k8s/gitops/project.yaml`. |
| Flux is the other standard agent but a different shape: a **pipeline of specialised controllers** (Source → Kustomize/Helm/Notification) with per-CRD reconcile intervals, supporting OCI sources ("Gitless GitOps") and image-automation CRDs. Not a single Application controller. | https://fluxcd.io/flux/concepts/ ; https://fluxcd.io/flux/components/ | 2026-10-13 | GA (v2, stable API) | Confirms "pick one and go deep": Flux is toolkit-centric, Argo CD is application-centric. Uncertainty: this is a synthesis from both vendors' docs; both are CNCF projects. | We chose Argo CD (AppProject + Application CRDs). `k8s/gitops/*.yaml`; `docs/huy/gitops-runbook.md`. |
| Flux reconciliation is per `Kustomization` with a required `.spec.interval` (min 60s), `prune`, `healthChecks` and `dependsOn` ordering; drift is re-checked every interval without a Git change. | https://fluxcd.io/flux/components/kustomize/kustomizations/ | 2026-10-13 | GA (v2) | Confirms Flux detects drift continuously, analogously to Argo's self-heal but tied to the Kustomization cycle. | Not used here (comparison only). |
| The **source → image-digest → PR** connection: an image-automation component watches a registry and writes the new reference back to git, and in PR mode opens a pull request that a human reviews before merge — the review is the promotion gate from build artifact to desired state. | https://argocd-image-updater.readthedocs.io/en/stable/basics/update-methods/ | 2026-10-13 | GA (v1.x); PR branch `image-updater-<targetKey>-<sha256>` | Confirms the digest-vs-tag distinction: the `digest` strategy takes a **tag** as input and tracks its SHA256 to detect change. Uncertainty: our pipeline implements the PR half with plain `git`/`gh` shell steps, not Argo CD Image Updater. | `propose-deploy` opens `deploy/huy-<sha>` changing exactly one line — the app image — to the full index digest; a human merges. `.github/workflows/ci.yml`, `captures/huy/bonus-1/propose-deploy.md`. |
| Argo CD Image Updater only manages resources Argo CD already manages, and its PR mode needs a token credential (PAT/GitHub App; SSH keys do not work). | https://argocd-image-updater.readthedocs.io/en/stable/ , https://argocd-image-updater.readthedocs.io/en/stable/basics/update-strategies/ | 2026-10-13 | GA (v1.x) | Uncertainty/scope note: this is the tool's constraint; we did not adopt Image Updater. | Our PR is created by the project PAT in CI (the CI secret is the same project credential, not a new identity). `docs/huy/gitops-runbook.md`. |
| **cloudflared quick tunnels** give an ephemeral public `*.trycloudflare.com` URL over an outbound-only connection with **no account or domain**, but they are for testing: no uptime guarantee, ~200 in-flight request cap, URL changes each run. A **production** Cloudflare Tunnel needs an account and a Cloudflare-managed domain. | https://developers.cloudflare.com/tunnel/get-started/quick-tunnels/ ; https://developers.cloudflare.com/tunnel/ | 2026-10-13 | GA (quick-tunnels page last updated 2026-09-30) | Confirms "public tunnel" + free-but-ephemeral. Correction: "available on all plans" ≠ free production tunnels; production routing needs a domain. | `cloudflared tunnel --url http://127.0.0.1:8080` → `https://debian-performances-truly-facing.trycloudflare.com`; owner's phone on **mobile data** got `Hello world from container-…-scwv5`. `captures/huy/step-3/tunnel.md`. |
| **ngrok** is a public tunnel (reverse proxy) exposing services publicly; its **free tier** is limited — 1 static dev domain (not customisable), up to 3 online endpoints, ~20k HTTP requests/month and ~1 GB/month, with an interstitial page on browser traffic. | https://ngrok.com/docs/pricing-limits/ | 2026-10-13 | Free plan (page modified 2026-10-12) | Correction: free ngrok is not unlimited — 3 endpoints / 20k HTTP req/month. Free endpoints do not time out, unlike a common assumption. | Not used here (comparison only); we used a cloudflared quick tunnel. |
| **Tailscale** is a **private mesh** (WireGuard, end-to-end encrypted), not a public tunnel: every device must run the client or be reached through a subnet router — there is no public "tunnel URL". The free **Personal** plan is up to 6 users with unlimited devices, and is **non-commercial use only**. | https://tailscale.com/pricing ; https://tailscale.com/docs/features/subnet-routers | 2026-10-13 | Free Personal plan (current pricing page) | Confirms public-demo (cloudflared/ngrok) vs team-internal (Tailscale) split. Correction: "unlimited devices" is per user, capped at 6 users free; free is non-commercial. | Tailscale `1.102.3` observed; used for the cluster/phone/worker transport, never as the public demo path. `docs/huy/versions.md`; `captures/huy/step-2/cluster.md`. |

## c) Containers & production-ready images

pending — source check in progress (lane `ses_ece6e865affeWYvpo35HrCXH7U`). Scope: container
isolation/runtime/layers/cache; non-root/minimal/distroless/no shell; secret-layer
persistence; external health vs probes; the installed Buildx SBOM syntax and inventory;
Trivy/Grype/Scout gates and scan limits; Docker Hardened Images provenance/SBOM/VEX checked in
Docker's own docs.

| Claim (owner words) | Primary source | Checked | Version/status | Confirmation / correction / uncertainty | Own-system evidence |
|---|---|---|---|---|---|
| _pending — awaiting lane result; no row invented._ | — | — | — | — | — |

## d) Orchestration — Kubernetes core concepts

pending — source check in progress (lane `ses_ece6e865affd3gJdVVwAmGRaHK`). Scope:
Pod/ReplicaSet/Deployment/Service/Namespace + control loops; scheduling/scaling/rolling
update/self-heal; ConfigMap/Secret env+file injection; base64-not-encryption + git danger;
emptyDir/pod-fs/PV/PVC + RWO vs RWX; local tool storage/default provisioner/version limits;
control plane/worker single/multi-node; the actual extra objects (ServiceMonitor/CRDs/
AppProject/Application) and why.

| Claim (owner words) | Primary source | Checked | Version/status | Confirmation / correction / uncertainty | Own-system evidence |
|---|---|---|---|---|---|
| _pending — awaiting lane result; no row invented._ | — | — | — | — | — |

## e) Observability with OpenTelemetry, Prometheus & Grafana

pending — source check in progress (lane `ses_ece6e8659ffe8JYUjlQ9Sxaiik`). Scope:
traces/metrics/logs vs plain logs; context correlation/W3C propagation/vendor neutrality;
each selected language SDK's signal maturity; Collector current major/release + changelog with
verified config stabilization; batching/routing/fan-out; versioned Semantic Conventions;
Prometheus pull/discovery/time-series/PromQL; Grafana workload + kubelet/cAdvisor; k6
traffic/threshold behavior.

| Claim (owner words) | Primary source | Checked | Version/status | Confirmation / correction / uncertainty | Own-system evidence |
|---|---|---|---|---|---|
| _pending — awaiting lane result; no row invented._ | — | — | — | — | — |

## f) Frontend & browser observability

Scope: RUM vs synthetic strengths; PerformanceObserver/Navigation/Resource Timing;
LCP/INP/CLS definitions/units + web-vitals; CrUX field aggregation vs DevTools/Trace Viewer
lab; current Web SDK/fetch/XHR/document-load APIs + maturity (EXPERIMENTAL label); the actual
`traceparent` joined path; Server-Timing/PerformanceServerTiming vs full SDK.

| Claim (owner words) | Primary source | Checked | Version/status | Confirmation / correction / uncertainty | Own-system evidence |
|---|---|---|---|---|---|
| **RUM (field)** data is what real users actually experienced and is reported as a distribution (75th percentile); **synthetic/lab** data is reproducible but controlled, so it misses real-world variance. Use field data to prioritise, lab data to debug. | https://web.dev/articles/lab-and-field-data-differences | 2026-10-13 | web.dev article (content dated 2022-07-18; vitals definitions it references are current) | Confirms the brief. Correction the page makes: lab can't simulate bfcache restores, mobile tap-delay or personalised content shifts — the browser SDK's single controlled session is **not** a RUM cohort. | Our page is a single controlled local/lab demo, labelled as such in the UI and in `web/src/main.js`; no field cohort. |
| `PerformanceObserver` is the entry point for the timeline: you subscribe to entry types (`navigation`, `resource`, `largest-contentful-paint`, `layout-shift`, `event`) instead of polling. | https://developer.mozilla.org/en-US/docs/Web/API/PerformanceObserver | 2026-10-13 | MDN Baseline / widely available (since Jan 2020); W3C Performance Timeline | Confirms this is the substrate web-vitals builds on. Uncertainty: MDN is normative-adjacent reference; the W3C Performance Timeline spec is the normative one. | `web-vitals` (installed in `web/package.json`) subscribes to these entry types; we do not hand-roll a PerformanceObserver. |
| **Navigation Timing** exposes one `PerformanceNavigationTiming` entry per document load with high-resolution timestamps (`startTime`, `domContentLoaded*`, `loadEvent*`, etc.). **Resource Timing** exposes a `PerformanceResourceTiming` per fetched resource, including a `serverTiming` array when the response sets `Server-Timing`. | https://developer.mozilla.org/en-US/docs/Web/API/Performance_API/Navigation_timing ; https://developer.mozilla.org/en-US/docs/Web/API/Performance_API/Resource_timing | 2026-10-13 | MDN; Navigation Timing L2 current; Resource Timing default buffer 250 entries | Confirms timing data is browser-provided. Correction: Navigation Timing Level 1 is deprecated — use `PerformanceNavigationTiming`. Cross-origin timing needs `Timing-Allow-Origin` (ours is same-origin). | `web/src/main.js` reads `performance.getEntriesByName(url,"resource")[0].serverTiming` for the same-origin `/test` response. |
| **Core Web Vitals units:** LCP and INP are **millisecond durations** (`DOMHighResTimeStamp`); CLS is a **unitless score** (impact fraction × distance fraction, max session window). Thresholds: LCP good ≤ 2.5 s, INP good ≤ 200 ms, CLS good ≤ 0.1. | https://web.dev/articles/lcp ; https://web.dev/articles/inp ; https://web.dev/articles/cls | 2026-10-13 | web.dev articles, current (LCP/INP updated Sep 2025; CLS Apr 2023) | Confirms the brief's LCP/INP/CLS. Correction: INP replaced FID (March 2024) and ignores scrolling/hover; CLS only counts shifts of **existing visible** elements. | UI encodes `VITAL_UNITS = { LCP: "ms", INP: "ms", CLS: "score (unitless)" }`; a real CLS of `0` is shown as a measured value, not a placeholder. `web/src/main.js`. |
| The **`web-vitals` library** is the easiest way to read the vitals and matches how Chrome reports them to CrUX/PSI/Search Console; it exposes `onLCP`/`onINP`/`onCLS` (plus FCP/TTFB) and an attribution build. | https://github.com/GoogleChrome/web-vitals | 2026-10-13 | npm `web-vitals` **6.2.3** installed here (README CDN examples still show v5) | Confirms it is the practical reader. Uncertainty: the repo README's CDN snippet lags the published latest — we pin the installed version explicitly. | `web/package.json` pins `web-vitals: 6.2.3`; `web/src/main.js` calls `onLCP/onINP/onCLS` with `reportAllChanges`. |
| **CrUX is field data** — 28-day rolling aggregation of real Chrome users, at origin/page level, 75th percentile; it excludes Chrome on iOS, Android WebView and non-Chromium browsers. **DevTools Performance panel / Trace Viewer is lab data** — one controlled recorded session. | https://developer.chrome.com/docs/crux ; https://developer.chrome.com/docs/devtools/performance | 2026-10-13 | CrUX GA; DevTools page updated 2024-10-23 | Confirms the field-vs-lab split the brief warns teams confuse. Correction: Trace Viewer is the same lab capture visualised, not a field source. | Our vitals are a lab measurement from one controlled session; the page states "not a CrUX statistic". `web/src/main.js`; `app/static/index.html`. |
| The **OpenTelemetry Web SDK** (`@opentelemetry/sdk-trace-web` + `instrumentation-fetch` / `-xml-http-request` / `-document-load`) emits the same trace/span model as the backend, but **browser instrumentation is explicitly EXPERIMENTAL and "mostly unspecified"** per the OTel docs. | https://opentelemetry.io/docs/languages/js/ | 2026-10-13 | OTel JS: Traces Stable, Metrics Stable, **Logs Development**; **browser instrumentation experimental** | Confirms the brief to verify current names, not blogs. Correction: the maturity gap is real — installing the packages does not give a stable API guarantee. Uncertainty: latest npm numbers drift; we record what we pinned. | Packages pinned in `web/package.json` (`sdk-trace-web` 2.12.0, `instrumentation-fetch` 0.223.0, etc.); `web/src/main.js` labels the setup EXPERIMENTAL. |
| **W3C Trace Context** defines `traceparent` (`version-trace-id-parent-id-flags`); OTel propagation **injects** it on the outgoing fetch/XHR and the backend **extracts** it, so the backend server span inherits the same `traceId` and parents under the browser span. Without it the backend starts an unrelated root span. | https://w3.org/TR/trace-context/ ; https://opentelemetry.io/docs/languages/js/propagation/ | 2026-10-13 | W3C Recommendation/standard; OTel JS propagation docs (SDK 2.x) | Confirms the joined-path mechanism. Our page never builds the header by hand — `provider.register()` installs the default W3C composite propagator. | One real trace: `test-click` → fetch client span (`traceparent` injected) → backend `otelhttp` server span → `test-work` child. `web/src/main.js`, `app/telemetry.go`, `app/main.go`; commits 73802bd/55a9989. |
| **Server-Timing** is the no-SDK option: the backend sets a `Server-Timing` response header and the browser reads it from `PerformanceResourceTiming.serverTiming` (`PerformanceServerTiming` = `name`/`duration`/`description`, same-origin). It is complementary to — not a replacement for — the Web SDK: no trace context or span relationships. | https://w3.org/TR/server-timing/ ; https://developer.mozilla.org/en-US/docs/Web/API/PerformanceServerTiming | 2026-10-13 | W3C Working Draft (IANA-registered header); MDN Baseline since Mar 2023 | Confirms "simpler first option before the full SDK". Correction: it carries server processing time only — it does not join a distributed trace. | `app/main.go` sets `Server-Timing: app;dur=<ms>`; `web/src/main.js` reads it back and shows it in the page; `app/main_test.go` asserts the header. |

## Moving-claim rechecks

Two claims are moving targets. They are **not** marked done ahead of time. Each carries one
row per recheck date; both are `pending` until the check is actually performed.

| Claim | Source | 16 Oct 2026 (pending) | 30 Oct 2026 (pending) |
|---|---|---|---|
| Immutable Actions rollout status and whether actions are published as OCI packages on `ghcr.io` | https://github.blog/changelog/ (see §a) | `pending` | `pending` |
| OpenTelemetry per-signal SDK maturity (traces/metrics/logs per selected language) | https://opentelemetry.io/status/ (see §e) | `pending` | `pending` |

**Rule:** a recheck row stays `pending` until the check is run; a future check is never
recorded as passed in advance.
