# Master Home Lab Setup Guide: Ubuntu Server 24.04 (Headless)

## Project Overview

**Goal:** Transform a spare PC ("city17") into a robust, headless Linux server for hosting services, file sharing, and acting as a central home lab node.

**OS:** Ubuntu Server 24.04 LTS  
**Architecture:** Headless (SSH only), Docker-based services, secure by design.

---

## Naming Theme: Half-Life Universe

All hostnames and service names use Valve's Half-Life universe.

### Established Names
| Name | Assignment |
|------|------------|
| `black mesa` | Home WiFi network |
| `aperture` | Personal PC (main workstation) |
| `citadel` | This git repository |
| `city17` | Main home lab server |

### Future Names (Reserved)
| Name | Suggested Role |
|------|----------------|
| `kleiner` | Secondary server / test node |
| `lambda` | Dedicated NAS |
| `ravenholm` | Security lab / honeypot |
| `eli` | Backup server |

---

## Hardware: city17 Server

| Component | Details |
|-----------|---------|
| OS Drive | NVMe (size TBD) |
| Storage Drive | 1TB HDD → `/srv/storage` |
| Backup Drive | 500GB HDD → `/srv/backups` |
| GPU | NVIDIA (drivers required) |

**Storage Strategy:** Tiered (no RAID - mismatched drive sizes)
- NVMe: Fast OS + Docker containers
- 1TB HDD: Samba shares, media, bulk files
- 500GB HDD: Nightly rsync backups of critical data

---

## Network Configuration

| Setting | Value |
|---------|-------|
| Subnet | `192.168.68.0/24` |
| Router/Gateway | `192.168.68.1` |
| IP Assignment | DHCP with router reservation |
| city17 Reserved IP | `192.168.68.52` |
| Discovery | Avahi mDNS (`city17.local`) |
| Raspberry Pi | `192.168.68.100` (existing) |

---

## Setup Phases Overview

| Phase | Description | Status |
|-------|-------------|--------|
| 0 | Pre-install prep (SSH key, bootable USB) | Completed |
| 1 | Ubuntu Server installation | Completed |
| 2 | Network setup (DHCP reservation + Avahi) | Completed |
| 5 | Docker + Portainer | Completed |
| 6 | Storage setup (tiered + Samba + backups) | Completed |
| 8 | Maintenance (unattended upgrades) | Completed |

**Phases Not Used:**
- **Phase 3:** Security Hardening (deferred - firewall/SSH hardening for future)
- **Phase 4:** NVIDIA Drivers (GPU disabled due to kernel instability)

---

## Quick Reference

### SSH Access
```bash
ssh breen@city17.local
```

### System Updates
```bash
sudo apt update && sudo apt upgrade -y
```

### Service Status
```bash
sudo systemctl status <service>
docker ps
```

### Storage Check
```bash
df -h /srv/storage /srv/backups
```

### Backup Status
```bash
cat /srv/backups/backup.log | tail -10
```

### GPU Status
```bash
nvidia-smi
```

---

## Detailed Instructions

For step-by-step commands and detailed setup instructions, see **[agents.md](agents.md)**.

---

## Future Project Ideas

| Service | Purpose | Suggested Name |
|---------|---------|----------------|
| Pi-hole / AdGuard Home | Network ad-blocking + local DNS | `pihole.local` |
| Jellyfin / Plex | Media server | `alyx.local` |
| Nextcloud | Self-hosted cloud storage | `lambda.local` |
| Nginx Proxy Manager | Reverse proxy | `gate.local` |
| Tailscale | Secure remote access VPN | N/A (overlay) |
| Proxmox (2nd machine) | Virtualization platform | `kleiner.local` |
| Security lab containers | DVWA, Metasploitable | `ravenholm.local` |
