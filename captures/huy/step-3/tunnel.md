# Stage 07 — off-machine (Step 3) evidence (huy)

Owner-performed, not simulated. Observed **2026-10-12**.

- **Device:** owner's Xiaomi Redmi Note 13 Pro 5G (Android 16), on **mobile data** (not the
  cluster machine, not workstation localhost).
- **URL:** `cloudflared` quick tunnel → `https://debian-performances-truly-facing.trycloudflare.com`
  (ephemeral), forwarding to the `devops-app` Service (`:8080`) in the `huy-k3s` cluster.
- **Request:** `GET /test` from the phone.
- **Response body (as seen on the phone):**
  `Hello world from container-devops-app-6f77b9498b-scwv5`
  (the ClusterIP Service selected pod `…scwv5`; the app name is the Downward-API `POD_NAME`.)

No personal phone/account identifiers, tokens, join keys, or kubeconfig are recorded here.
Only the app was exposed; dashboards/GitOps/trace viewer/API endpoints stayed private.

The quick tunnel is ephemeral and was stopped after this capture; the runbook in
`docs/huy/deploy-runbook.md` recreates it and re-reads the new URL on restart.
