# Stage 06: Two physical hosts, private k3s cluster

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 06](../../../.scratch/assignment-3-personal/issues/06-cluster.md), Step 2/bonus 2/research d, overview Approach 6, [execution contract](00-execution-contract.md).

## Entry gate and routing

Stage 01 accepted. Local tool/access preparation may overlap CI; no unverified image deployment. Preliminary complex-critical: host installation, network/security boundaries, admin kubeconfig and remote join credentials. BEFORE remote work obtain actual allocated VM access, distinct physical host confirmation, resources, installation and tailnet permission. Unavailable facts are human prerequisites, not an invitation to guess another host or provision a local substitute.

Read [Lima Linux installation](https://github.com/lima-vm/lima/blob/master/website/content/en/docs/installation/_index.md), [k3s requirements](https://docs.k3s.io/installation/requirements), [network options](https://docs.k3s.io/networking/basic-network-options), server/agent CLI references and Tailscale current installation/ACL guidance. Official public release downloads can be anonymous; any authenticated GitHub release call uses the project .env policy.

## Allowed scope and branch

Source branch stage/huy/06-cluster from accepted origin/personal/huy, matching worktree. Reproducibility files docs/huy/cluster-runbook.md, versions/lessons and captures/huy/step-2/, bonus-2/. Permitted operations only owner's named huy-server VM, explicitly allocated huy-worker VM and authorized workstation client. Private files /etc/rancher/k3s/config.yaml and token file on allocated hosts, $HOME/.kube/huy-k3s.yaml locally; never committed. No existing cluster replacement, unrelated host/group/firewall changes, paid service, public API/VXLAN or production access. Commit docs/pins/scrubbed captures only.

Commit scope: only cluster runbook, verified version/config descriptions, lessons and scrubbed Step 2/bonus 2 captures. Live kubeconfig, join token and machine credential files are excluded.

## Bounded substeps

### 06A: Verify tools/allocation and create server VM

Install stable Lima/QEMU, kubectl and Helm as needed from official packages/releases with verified checksum; record versions/architectures. Inspect CLI help before VM creation. Initial local budget 4 vCPU, 4 GiB, 40 GiB; remote worker 2 vCPU, 4 GiB, 30 GiB. These are choices, not guarantees. Existing VM name collision blocks automatic reuse/deletion.

```sh
limactl start --help
limactl start --name=huy-server --vm-type=qemu --cpus=4 --memory=4 --disk=40 template:ubuntu
limactl shell huy-server
```

Inspect actual guest OS/architecture/resources/clock. If installed Lima changes Ubuntu template syntax, use shipped Ubuntu template with SAME topology/resources and record actual command. If virtualization is unavailable, report; native k3s on own host requires explicit revised host-install authorization, not automatic fallback.

Exit gate: actual permitted isolated server guest runs; worker allocation facts recorded before 06B/06C remote operations.

### 06B: Private connectivity

Authenticate Tailscale inside both VMs AND on the workstation that will run kubectl against the private API. OAuth/join-key handling is owner's operation; no key in prompts/git/captures. Record actual SERVER_IP, WORKER_IP and WORKSTATION_IP from tailscale ip -4 and host mapping. Confirm bidirectional VM reachability and workstation-to-server reachability, time synchronization and permitted free-tier tailnet ownership.

Allow TCP 6443 between worker/workstation and server, TCP 10250 and UDP 8472 between the two VMs, restricted to their actual private addresses/authorized identities. Tailnet ACL and host firewall must agree. Do not open VXLAN/API on public interfaces. Do not advertise Kubernetes Pod CIDRs through Tailscale unnecessarily; Flannel carries cross-node pod traffic.

Exit gate: authenticated permitted mesh connectivity and API-client route established, without assuming local CLI installation means login/service is active.

### 06C: Pinned same-version server and agent

Resolve one current supported stable k3s release from official docs/release checksum and verify compatibility with planned Helm/Argo versions. Pin installer source/version and review before authorized execution. Do not pipe moving remote script into root. Write root-owned config before starting the service.

Server config values: node-name huy-server; node-ip and advertise-address SERVER_IP; tls-san array containing SERVER_IP; flannel-iface tailscale0; disable traefik and servicelb. Single-server default SQLite. No HA/embedded-etcd claim or unnecessary cluster-init.

Worker config: node-name huy-worker; server https://SERVER_IP:6443; node-ip WORKER_IP; flannel-iface tailscale0; token-file path to root-readable privately transferred server join token. Install explicit agent role with the SAME version. Transfer the token through owner-approved private file channel, not console/chat/capture; permissions 0600. Server and worker are separate OS service instances.

Copy kubeconfig securely outside Git to $HOME/.kube/huy-k3s.yaml, mode 0600, replace API host with SERVER_IP and rename its context huy-k3s without overwriting unrelated kubeconfig. Private client certificates/keys stay private; never kubectl config view --raw into evidence.

## Verification and evidence

On workstation:

```sh
export KUBECONFIG="$HOME/.kube/huy-k3s.yaml"
kubectl config current-context
kubectl get nodes -o wide
kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{" "}{.status.nodeInfo.kubeletVersion}{" "}{.status.nodeInfo.architecture}{"\n"}{end}'
kubectl get storageclass
kubectl -n kube-system get pods -o wide
```

Sisyphus verifies context, actual two Ready nodes, same version, default local-path provisioner and system pods. Capture separate physical-host/VM mapping with non-secret owner evidence, not merely node names. Stage 07 then proves pod spread and actual cross-node Pod/Service traffic; Ready nodes alone cannot close bonus 2. Node/server logs may be captured only after token scrub; don't disclose kubeconfig or node-token.

Save actual installation/version/resource/network failures in lessons.md. If encapsulation/MTU fails, inspect private reachability, firewall, Flannel interface/routes and payload path before declaring networking passed. After Tailscale interface recreation, restart only affected authorized k3s service and repeat acceptance; do not reinstall or erase data.

## Learning gate

Owner identifies control plane versus agent, SQLite single-server limits, actual node IPs, private API client route, cross-machine scheduling/network implications and local-path node binding. Explain why a second local VM and two names do not prove bonus 2, and why root kubeconfig/join token must not be public.

## Completion and blockers

Two-host node/system acceptance and private boundaries pass before Stage 07. If worker allocation is missing, local preparation can complete but this stage remains BLOCKED. Resource issues require observed scheduling/OOM diagnostics and authorized budget adjustment; no dropping required components. No claim of HA/node-loss durability.

## Sisyphus handoff prompt

```text
Coordinate Stage 06/ticket 06 after Stage 01. Preliminary complex-critical. Do not
implement or install yourself. Obtain exact host/allocation/tailnet authorization,
prepare the docs stage worktree, and separate owner OAuth/private-token transfer
from bounded authorized host operations. Dispatch only named machines/configs.
Require same pinned k3s version, single-server SQLite, tailscale0 private routing,
workstation API access and private kubeconfig. Independently observe two Ready
nodes/system storage and physical-host mapping; retain missing worker as BLOCKED.
Require owner's control-plane/network/security learning evidence. Commit only
scrubbed runbook/pins/captures, never join tokens or client certificates.
```
