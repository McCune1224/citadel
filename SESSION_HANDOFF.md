# Session Handoff: city17 Home Lab Setup

This document captures the current state of the home lab setup session for continuation on another machine.

---

## Current Status

**Phase 1 (OS Installation): IN PROGRESS**

The user is currently running through the Ubuntu Server 24.04 installer on the spare PC (city17).

---

## Decisions Made

| Decision | Value |
|----------|-------|
| Server hostname | `city17` |
| Username | `breen` |
| Theme | Half-Life universe |
| Network approach | DHCP reservation + Avahi mDNS |
| Reserved IP | `192.168.68.101` |
| Discovery | `city17.local` via Avahi |
| Storage strategy | Tiered (no RAID) |
| OS drive | NVMe |
| Data drive | 1TB HDD -> `/srv/storage` |
| Backup drive | 500GB HDD -> `/srv/backups` |
| GPU | NVIDIA (drivers needed) |
| SSH key import | From GitHub username `McCune1224` |

---

## Names Already In Use

| Name | Assignment |
|------|------------|
| `black mesa` | Home WiFi network |
| `aperture` | Personal PC |
| `citadel` | This git repository |
| `city17` | Main home lab server (being set up now) |

---

## Network Info

| Device | IP |
|--------|-----|
| Router/Gateway | `192.168.68.1` |
| Raspberry Pi | `192.168.68.100` |
| city17 (reserved) | `192.168.68.101` |

---

## What's Completed

- [x] Phase 0: Pre-install prep
  - [x] SSH keypair exists on personal PC (`~/.ssh/id_ed25519`)
  - [x] Public key is on GitHub (username: `McCune1224`)
  - [x] Ubuntu Server 24.04 ISO downloaded
  - [x] Bootable USB created (using Popsicle after wiping with `wipefs`)
  - [x] Documentation updated (`agents.md`, `master_lab_setup_guide.md`, `README.md`)

- [ ] Phase 1: OS Installation - **IN PROGRESS**
  - User is running through the Ubuntu installer now
  - Settings to use:
    - Hostname: `city17`
    - Username: `breen`
    - Enable OpenSSH server
    - Import SSH keys from GitHub: `McCune1224`
    - Skip all featured snaps
    - Accept third-party drivers if prompted
    - Install to NVMe only, leave HDDs alone

---

## Next Steps After Install Completes

1. **Find city17's IP address** (from router admin panel or network scan)
2. **First SSH connection**: `ssh breen@<ip-address>`
3. **Phase 2: Network Setup**
   - Get MAC address from city17
   - Create DHCP reservation in router for `192.168.68.101`
   - Install Avahi: `sudo apt install avahi-daemon -y`
   - Verify: `ping city17.local`
4. **Phase 3: Security Hardening** (SSH, UFW, Fail2ban)
5. **Phase 4: NVIDIA Drivers**
6. **Phase 5: Docker + Portainer**
7. **Phase 6: Storage Setup** (mount HDDs, Samba, backup script)
8. **Phase 7: Monitoring** (Netdata)
9. **Phase 8: Maintenance** (unattended upgrades)

---

## Key Documentation Files

All detailed step-by-step commands are in:

- `agents.md` - Full detailed instructions for all phases
- `master_lab_setup_guide.md` - High-level overview and quick reference
- `README.md` - Repo overview

---

## User Preferences

- **NO EMOJIS** in any documentation or responses
- Half-Life naming theme (not Portal, not TF2)
- Prefers explanations of options before making decisions
- Using laptop to continue session (personal PC is "aperture")

---

## How to Continue This Session

1. Pull this repo on the laptop
2. Open OpenCode in the citadel directory
3. Tell OpenCode: "Continue the city17 home lab setup - the Ubuntu install should be done. See SESSION_HANDOFF.md for context."
4. OpenCode will read this file and the other docs to pick up where we left off

---

## Useful Commands for Continuation

```bash
# Find city17 on the network (run from laptop)
nmap -sn 192.168.68.0/24

# Or check router's DHCP client list at:
# http://192.168.68.1

# First SSH (before Avahi is set up)
ssh breen@<city17-ip>

# After Avahi is configured
ssh breen@city17.local
```
