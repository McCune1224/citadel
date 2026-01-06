# Session Handoff: city17 Home Lab Setup

This document captures the current state of the home lab setup session for continuation.

---

## Current Status

**Phase 6 (Storage Setup): COMPLETED**
**Phase 3 (Security Hardening): SKIPPED FOR NOW**
**Phase 7 (Monitoring - Netdata): PENDING - START HERE NEXT SESSION**

---

## Decisions Made

| Decision | Value |
|----------|-------|
| Server hostname | `city17` |
| Username | `breen` |
| Theme | Half-Life universe |
| Network approach | DHCP reservation + Avahi mDNS |
| Reserved IP | `192.168.68.52` |
| Discovery | `city17.local` via Avahi |
| Storage strategy | Tiered (no RAID) |
| OS drive | NVMe (238.5GB Samsung) |
| Data drive | 931.5GB Toshiba HDD (`/dev/sdb`) -> `/srv/storage` |
| Backup drive | 465.8GB Toshiba HDD (`/dev/sda`) -> `/srv/backups` |
| GPU | NVIDIA GeForce GTX 1070 |
| SSH key import | From GitHub username `McCune1224` |

---

## Names Already In Use

| Name | Assignment |
|------|------------|
| `black mesa` | Home WiFi network |
| `aperture` | Personal PC |
| `citadel` | This git repository |
| `city17` | Main home lab server |

---

## Network Info

| Device | IP |
|--------|-----|
| Router/Gateway | `192.168.68.1` |
| Raspberry Pi | `192.168.68.100` |
| city17 | `192.168.68.52` (DHCP reserved) |

---

## Hardware: city17

| Component | Details |
|-----------|---------|
| Ethernet MAC | `c8:d3:ff:41:af:89` (interface: `enp5s0`) |
| GPU | NVIDIA GeForce GTX 1070 |
| NVMe | 238.5GB Samsung (OS installed) |
| HDD 1 | 465.8GB Toshiba (`/dev/sda`) - backup drive |
| HDD 2 | 931.5GB Toshiba (`/dev/sdb`) - main storage |

Note: HDDs have been wiped and formatted as ext4.

---

## What's Completed

- [x] Phase 0: Pre-install prep
- [x] Phase 1: OS Installation (Ubuntu Server 24.04)
- [x] Phase 2: Network Setup
  - [x] DHCP reservation created (192.168.68.52)
  - [x] Avahi installed and working
  - [x] `city17.local` resolves correctly
  - [x] kitty-terminfo installed (optional)
- [x] Phase 4: NVIDIA Drivers
  - [x] nvidia-driver-570 installed
  - [x] nvidia-smi working (GTX 1070, CUDA 12.8)
- [x] Phase 5: Docker + Portainer
  - [x] Docker CE installed
  - [x] breen added to docker group
  - [x] Portainer running on port 9443
- [x] Phase 6: Storage Setup
  - [x] /dev/sda formatted as ext4 (backups)
  - [x] /dev/sdb formatted as ext4 (storage)
  - [x] Mounted at /srv/backups and /srv/storage
  - [x] Directory structure created
  - [x] Samba installed and configured (public + private shares)
  - [x] Backup script scheduled (3 AM nightly)

---

## What's Next

- [ ] Phase 3: Security Hardening (SSH, UFW, Fail2ban) - skipped for now
- [ ] Phase 7: Monitoring (Netdata)
- [ ] Phase 8: Maintenance (unattended upgrades)

---

## How to Connect

```bash
ssh breen@city17.local
```

---

## Services Running

| Service | Access | Notes |
|---------|--------|-------|
| Portainer | `https://city17.local:9443` | Docker web UI, self-signed cert |
| Samba (public) | `smb://city17.local/public` | Guest access, no password |
| Samba (private) | `smb://city17.local/private` | Requires breen + samba password |

---

## Technologies Installed

| Technology | Version | Purpose |
|------------|---------|---------|
| Ubuntu Server | 24.04 LTS | Operating system |
| Docker CE | latest | Container runtime |
| Portainer CE | latest | Docker web management |
| Samba | apt default | SMB file sharing |
| NVIDIA Driver | 570.195.03 | GPU driver |
| CUDA | 12.8 | GPU compute (available) |
| Avahi | apt default | mDNS (.local discovery) |
| rsync | apt default | Backup script |

---

## Phase 3 Commands (Ready to Run)

When resuming, SSH into city17 and run:

### Step 1: System Updates
```bash
sudo apt update && sudo apt upgrade -y
```

### Step 2: SSH Hardening
```bash
sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo sed -i 's/PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config
sudo systemctl restart ssh
```

### Step 3: Firewall (UFW)
```bash
sudo ufw allow ssh
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw enable
```

### Step 4: Fail2ban
```bash
sudo apt install -y fail2ban
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
sudo sed -i '/^\[sshd\]/a enabled = true' /etc/fail2ban/jail.local
sudo systemctl enable fail2ban
sudo systemctl restart fail2ban
```

### Step 5: Verify
```bash
sudo ufw status
sudo fail2ban-client status sshd
```

---

## Storage UUIDs

| Drive | Label | UUID | Mount |
|-------|-------|------|-------|
| /dev/sda | backups | c3aa6648-ee8c-4b4d-8e59-f700e635c8c0 | /srv/backups |
| /dev/sdb | storage | b89fc0e9-b482-43be-9154-28db15de750e | /srv/storage |

---

## Key Documentation Files

- `agents.md` - Full detailed instructions for all phases
- `master_lab_setup_guide.md` - High-level overview and quick reference
- `README.md` - Repo overview

---

## User Preferences

- **NO EMOJIS** in any documentation or responses
- Half-Life naming theme
- Prefers explanations before making decisions
