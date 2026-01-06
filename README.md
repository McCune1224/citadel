# Citadel

Home lab infrastructure repository for managing server configurations, documentation, and automation.

## Overview

This repo contains setup guides and configuration for a headless Ubuntu Server home lab environment.

## Servers

| Hostname | Role | IP | Status |
|----------|------|-----|--------|
| city17 | Main home lab server | `192.168.68.101` | Setup in progress |

## Documentation

| File | Description |
|------|-------------|
| [agents.md](agents.md) | Detailed step-by-step setup instructions |
| [master_lab_setup_guide.md](master_lab_setup_guide.md) | High-level overview and quick reference |

## Naming Theme

All hostnames follow Valve's Half-Life universe naming convention.

## Network

- **Subnet:** `192.168.68.0/24`
- **Discovery:** Avahi mDNS (`.local` domains)
- **Access:** SSH key-based authentication only
