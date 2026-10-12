# k6 timed run 3 (the valid synchronized run)
- Date: 2026-10-12 19:15-19:17 UTC. k6 v2.3.0, script load-tests/test.js at stage commit 4478007.
- Target: private forward 127.0.0.1:8080 to svc/devops-app (bind verified 200 before start).
- App digest sha256:4b46f576, 2 pods 698d766b86. Dashboard huy-workload, 5s refresh.
- Result: 26,420 checks, 100 percent pass. Thresholds held: checks rate 1, success-case failures under 0.01.
- Mid-run panel snapshot: rate 105.2/s, error ratio 0.0997, p95 0.175s, per-pod CPU split honest through single-pod forward.
- Runs 1 and 2 same evening are INVALID for dashboard sync: a stale local 11A-lane binary held host :8080, so the forward never bound and traffic never reached the cluster. Their client-side checks stand, their dashboard numbers do not. Stale process killed, forward bind verified before run 3.
