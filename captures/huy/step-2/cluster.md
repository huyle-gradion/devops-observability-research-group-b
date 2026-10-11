# Stage 06 — k3s cluster evidence (huy)

Produced by Sisyphus. Observed **2026-10-11**. No secrets (kubeconfig kept local, untracked).

## Topology (as built — single node)

| Role | Host | How |
|---|---|---|
| k3s server (control-plane) | `omarchy` (Arch/Omarchy, Linux amd64) | Linux VM under **QEMU + session libvirt** |
| external client (Step 3) | Redmi Note 13 Pro 5G (Android, Tailscale) | phone |
| second physical host (bonus 2) | macOS 26 Intel (Tailscale, no inbound SSH) | **deferred** |

The plan assumed two Ubuntu VMs; the actual second host is macOS-Intel (needs a Linux VM,
owner-run, no inbound SSH), so per the owner's decision we built **single-node first**, then
extend for bonus 2.

## VM

- Ubuntu 24.04.5 cloud image; qcow2 overlay `~/vms/huy-k3s/disk.qcow2` (40 GiB, backed by
  the 597 MiB cloud image); 2 vCPU / 4 GiB; `/dev/kvm` acceleration.
- Networking: QEMU **user-mode** (`-netdev user`) with host forwards
  `127.0.0.1:2222→:22` and `127.0.0.1:6443→:6443`. SSH via a dedicated keyless key
  `~/.ssh/huy_k3s`.
- Started with `qemu-system-x86_64 ... -daemonize -pidfile ~/vms/huy-k3s/qemu.pid`.

## k3s

- Version **v1.37.1+k3s1** (release commit `356c0254`), installed from the **pinned GitHub
  release binary** with `sha256sum -c` verification (the moving `get.k3s.io | sh` path was
  avoided, per the plan).
- systemd unit `k3s.service` runs `k3s server --write-kubeconfig-mode 644`.
- `kubectl` v1.37.1 on the host; kubeconfig `~/.kube/huy-k3s.yaml` (mode 0600, untracked).

## Verified from the host

```
NAME      STATUS   ROLES           AGE   VERSION        INTERNAL-IP   OS-IMAGE             CONTAINER-RUNTIME
huy-k3s   Ready    control-plane   16s   v1.37.1+k3s1   10.0.2.15     Ubuntu 24.04.5 LTS   containerd://2.3.4-k3s1
```

- Namespaces: `default`, `kube-node-lease`, `kube-public`, `kube-system`.
- StorageClass: `local-path` (`rancher.io/local-path`, default).
- `kubectl cluster-info` → control plane at `https://127.0.0.1:6443`.

## Boundaries / deviations

- **Single-node** → satisfies required Step 2. **Bonus 2** (pods spread across two physical
  hosts, cross-node pod-IP/service responses) is **not** done; needs the macOS Linux VM.
- No arm64/k3s-arch claim; the node is amd64.
- The cluster is a VM on the workstation, not a bare-metal host — recorded honestly.
