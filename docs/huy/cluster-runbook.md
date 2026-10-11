# Cluster runbook — k3s in a QEMU VM (huy)

Fast rebuild of the Stage 06 single-node cluster on omarchy, no root for the VM itself
(session libvirt / direct QEMU). Owner runs any `sudo pacman` install once.

## 0. One-time host packages (owner, sudo)

```sh
sudo pacman -S --needed qemu-full libvirt virt-install edk2-ovmf dnsmasq
```

## 1. Base image + overlay disk

```sh
mkdir -p ~/vms/huy-k3s && cd ~/vms/huy-k3s
curl -fL -o noble-server-cloudimg-amd64.img \
  https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img
qemu-img create -f qcow2 -F qcow2 -b noble-server-cloudimg-amd64.img disk.qcow2 40G
```

## 2. Cloud-init seed (keyless key!)

```sh
ssh-keygen -t ed25519 -N "" -C huy-k3s-vm -f ~/.ssh/huy_k3s   # passphrase-less, VM-only
cat > user-data <<EOF
#cloud-config
hostname: huy-k3s
users:
  - name: huy
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: [sudo, adm]
    shell: /bin/bash
    ssh_authorized_keys:
      - $(cat ~/.ssh/huy_k3s.pub)
ssh_pwauth: false
EOF
printf 'instance-id: huy-k3s-01\nlocal-hostname: huy-k3s\n' > meta-data
xorriso -as mkisofs -output seed.iso -volid cidata -joliet -rock user-data meta-data
```

**Trap:** an interactive/passphrase-protected key breaks non-interactive SSH. Use a dedicated
passphrase-less key for the VM.

## 3. Boot the VM (user-mode networking + port forwards)

```sh
qemu-system-x86_64 -enable-kvm -cpu host -m 4096 -smp 2 \
  -drive file=disk.qcow2,if=virtio,format=qcow2 \
  -drive file=seed.iso,if=virtio,format=raw,readonly=on \
  -netdev user,id=n0,hostfwd=tcp:127.0.0.1:2222-:22,hostfwd=tcp:127.0.0.1:6443-:6443 \
  -device virtio-net-pci,netdev=n0 \
  -display none -serial file:serial.log -daemonize -pidfile qemu.pid
```

## 4. Install k3s (pinned binary + checksum)

```sh
TAG=v1.37.1+k3s1
curl -fsSL -o /tmp/k3s "https://github.com/k3s-io/k3s/releases/download/${TAG}/k3s"
curl -fsSL -o /tmp/k3s.sha256 "https://github.com/k3s-io/k3s/releases/download/${TAG}/sha256sum-amd64.txt"
( cd /tmp && grep ' k3s$' k3s.sha256 > k3s.expected && sha256sum -c k3s.expected )
sudo install -m0755 /tmp/k3s /usr/local/bin/k3s
# systemd unit runs: /usr/local/bin/k3s server --write-kubeconfig-mode 644
```

## 5. Kubeconfig to the host

```sh
ssh -p 2222 -i ~/.ssh/huy_k3s huy@127.0.0.1 'sudo cat /etc/rancher/k3s/k3s.yaml' \
  > ~/.kube/huy-k3s.yaml && chmod 600 ~/.kube/huy-k3s.yaml
export KUBECONFIG=$HOME/.kube/huy-k3s.yaml
kubectl get nodes -o wide
```

## Traps learned

- Passphrase-protected SSH key → non-interactive auth fails; use a keyless VM key.
- `get.k3s.io | sh` is a moving script piped to root — pin the release binary + checksum.
- Cloud-init `runcmd` k3s install didn't take; install explicitly over SSH and verify.
- User-mode networking hides the guest IP; reach the cluster via the `hostfwd` ports.
- Upgrading to two hosts (bonus 2) needs a Linux VM on the macOS host (owner-run; no inbound SSH).
