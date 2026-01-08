# AI Agent Instruction Protocol: Headless Ubuntu Server Setup

## Agent Role & Context
**Role:** Senior DevOps Engineer / Linux System Administrator  
**Objective:** Guide the user through setting up a "Master" Headless Home Lab Server on a spare PC.

**OS Target:** Ubuntu Server 24.04 LTS  
**Constraints:**
- **Headless Environment:** No monitor/keyboard available after initial install.
- **Remote Access:** SSH is the primary access method.
- **Infrastructure as Code:** Prefer declarative configs (Docker Compose) over manual commands where possible.

---

## Naming Convention Theme (Half-Life Universe)

All hostnames, service names, and internal domains are derived from **Valve's Half-Life universe**.

### Established Names (Already In Use)
| Name | Assignment |
|------|------------|
| `black mesa` | Home WiFi network |
| `aperture` | Personal PC (the machine running this agent) |
| `citadel` | This git repository |

### Server Naming
| Hostname | Role | Access |
|----------|------|--------|
| `city17` | Main home lab server | `city17.local` via Avahi/mDNS |

### Future Expansion Ideas
- `kleiner` - Secondary server / test node
- `lambda` - Storage / NAS dedicated box
- `ravenholm` - Honeypot / security lab (fitting name!)
- `eli` - Backup server

### Naming Rules
- Prefer short, memorable names easy to type in SSH commands.
- Keep names stable once assigned; do not rename without explicit user confirmation.
- All machines accessible via `<hostname>.local` using Avahi mDNS.

---

## Network Configuration

| Setting | Value |
|---------|-------|
| Subnet | `192.168.68.0/24` |
| Router/Gateway | `192.168.68.1` (confirm with user) |
| IP Assignment | DHCP with router reservation |
| city17 IP | `192.168.68.52` (DHCP reserved) |
| Discovery | Avahi mDNS (`city17.local`) |
| Raspberry Pi | `192.168.68.100` (existing) |

---

## Hardware Notes

**city17 Server:**
- NVMe drive (primary OS drive)
- 1TB HDD → `/srv/storage` (main bulk data storage)
- 500GB HDD → `/srv/backups` (local backup drive)
- NVIDIA GPU (drivers required)

**Storage Strategy:** Tiered storage (no RAID - drives are mismatched sizes)
- NVMe: Fast OS + Docker containers
- 1TB HDD: Samba shares, media, large files
- 500GB HDD: Nightly rsync backups of critical data from 1TB + configs

---

## Phase 0: Pre-Install Preparation (On Personal PC - "aperture")

### Step 0.1: Generate SSH Keypair
**Location:** Run on your personal PC (aperture)

```bash
# Generate Ed25519 key (modern, secure, fast)
ssh-keygen -t ed25519 -C "breen@city17"

# When prompted for location, press Enter for default (~/.ssh/id_ed25519)
# When prompted for passphrase, enter a strong passphrase (recommended) or leave blank
```

**Validation:** 
```bash
# Verify key was created
ls -la ~/.ssh/id_ed25519*

# Display public key (you'll need this during install)
cat ~/.ssh/id_ed25519.pub
```

### Step 0.2: Download Ubuntu Server 24.04 LTS
**Download URL:** https://ubuntu.com/download/server

Direct link (as of writing): Ubuntu Server 24.04.x LTS ISO

### Step 0.3: Create Bootable USB
**Recommended Tool:** Balena Etcher (cross-platform, simple)
- Download: https://etcher.balena.io/
- Select ISO, select USB drive, flash.

**Alternative (Linux command line):**
```bash
# Find your USB device (BE CAREFUL - wrong device = data loss)
lsblk

# Write ISO to USB (replace /dev/sdX with your USB device)
sudo dd if=ubuntu-24.04-live-server-amd64.iso of=/dev/sdX bs=4M status=progress conv=fsync
```

### Step 0.4: Note Router Admin Access
- Router admin URL: `http://192.168.68.1` (typical, verify)
- Have credentials ready for DHCP reservation setup later

---

## Phase 1: OS Installation (Temporary Monitor/Keyboard Required)

### Step 1.1: Boot and Start Installer
1. Insert USB into city17
2. Power on, enter BIOS/boot menu (usually F2, F12, Del, or Esc)
3. Select USB as boot device
4. Choose "Install Ubuntu Server"

### Step 1.2: Installation Options

| Setting | Value |
|---------|-------|
| Language | English |
| Keyboard | Your preference |
| Installation type | Ubuntu Server |
| Network | DHCP (auto) - we'll set reservation later |
| Proxy | Leave blank (unless you have one) |
| Mirror | Default |
| Storage | Use entire NVMe disk for OS (guided) |
| **Do NOT touch the HDDs during install** | Leave them for later RAID setup |

### Step 1.3: Profile Setup

| Field | Value |
|-------|-------|
| Your name | Breen |
| Server name | `city17` |
| Username | `breen` |
| Password | Choose a strong password |

### Step 1.4: SSH Setup
- [x] Install OpenSSH server
- Import SSH identity: **from GitHub**
  - Enter your GitHub username when prompted
  - Installer will fetch your public keys automatically

### Step 1.5: Featured Server Snaps
- Skip all snaps for now (we'll install Docker properly later)

### Step 1.6: Third-Party Drivers
- If prompted about third-party/proprietary drivers: **Accept/Enable**
- This helps with NVIDIA GPU detection

### Step 1.7: Complete Installation
1. Wait for installation to complete
2. Select "Reboot Now"
3. Remove USB when prompted
4. **Disconnect monitor/keyboard** - you're going headless!

### Step 1.8: Find city17's IP Address
From your personal PC (aperture):

```bash
# Option A: Check router's DHCP client list
# Log into http://192.168.68.1 and look for "city17"

# Option B: Scan network (if you have nmap)
nmap -sn 192.168.68.0/24

# Option C: Try mDNS (may work immediately)
ping city17.local
```

### Step 1.9: First SSH Connection
```bash
ssh breen@<city17-ip-address>
# OR if Avahi is already working:
ssh breen@city17.local
```

**Validation:** You should connect without password prompt (using your SSH key).

---

## Phase 2: Network Setup (DHCP Reservation + Avahi)

### Step 2.1: Get city17's MAC Address
On city17 (via SSH):
```bash
ip link show
# Look for the MAC address (link/ether xx:xx:xx:xx:xx:xx) of your primary interface (usually enp*, eth*, or eno*)
```

### Step 2.2: Create DHCP Reservation in Router
1. Log into router admin panel (`http://192.168.68.1`)
2. Find DHCP settings / Address Reservation / Static Leases
3. Add reservation:
   - MAC Address: (from step 2.1)
   - IP Address: `192.168.68.52`
   - Hostname: `city17`
4. Save/Apply

### Step 2.3: Install and Configure Avahi
On city17:
```bash
# Update package lists
sudo apt update

# Install Avahi daemon
sudo apt install avahi-daemon -y

# Enable and start service
sudo systemctl enable avahi-daemon
sudo systemctl start avahi-daemon

# Verify it's running
sudo systemctl status avahi-daemon
```

### Step 2.4: Reboot to Apply DHCP Reservation
```bash
sudo reboot
```

### Step 2.5: Validation
From your personal PC (aperture):
```bash
# Should resolve to 192.168.68.52
ping city17.local

# SSH using .local name
ssh breen@city17.local
```

---

## Phase 4: NVIDIA Drivers

### Step 4.1: Check Current GPU Status
```bash
# See if GPU is detected
lspci | grep -i nvidia

# Check if nouveau (open source driver) is loaded
lsmod | grep nouveau
```

### Step 4.2: Install NVIDIA Drivers
```bash
# Update package lists
sudo apt update

# Install recommended driver automatically
sudo ubuntu-drivers install

# OR see available drivers and choose:
ubuntu-drivers devices
sudo apt install nvidia-driver-550  # (or whatever version is recommended)
```

### Step 4.3: Reboot and Verify
```bash
sudo reboot

# After reboot, SSH back in and verify:
nvidia-smi
```

**Expected output:** GPU info table showing driver version, GPU name, memory, etc.

### Step 4.4: (Optional) NVIDIA Container Toolkit
For GPU-accelerated Docker containers in the future:
```bash
# Add NVIDIA container toolkit repo
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg

curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
  sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
  sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

sudo apt update
sudo apt install -y nvidia-container-toolkit

# Configure Docker to use NVIDIA runtime
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker
```

---

## Phase 5: Docker + Portainer

### Step 5.1: Install Docker Engine
```bash
# Remove any old versions
sudo apt remove docker docker-engine docker.io containerd runc 2>/dev/null

# Install prerequisites
sudo apt update
sudo apt install -y ca-certificates curl gnupg

# Add Docker's official GPG key
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Add Docker repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Add user to docker group (no sudo needed for docker commands)
sudo usermod -aG docker breen

# Apply group change (or log out and back in)
newgrp docker

# Verify
docker run hello-world
```

### Step 5.2: Deploy Portainer
```bash
# Create volume for Portainer data
docker volume create portainer_data

# Run Portainer
docker run -d \
  -p 9443:9443 \
  --name portainer \
  --restart=always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v portainer_data:/data \
  portainer/portainer-ce:latest
```

### Step 5.3: Access Portainer
1. Open browser: `https://city17.local:9443`
2. Accept self-signed certificate warning
3. Create admin user on first visit
4. Select "Get Started" → Manage local Docker environment

---

## Phase 6: Storage Setup (Tiered Storage + Samba)

### Step 6.1: Identify HDDs
```bash
# List all block devices
lsblk

# Get more details (note sizes to identify which is 1TB vs 500GB)
sudo fdisk -l

# Example output:
# /dev/nvme0n1 - NVMe (OS drive)
# /dev/sda - 1TB HDD (main storage)
# /dev/sdb - 500GB HDD (backups)
```

### Step 6.2: Partition and Format HDDs
```bash
# Format the 1TB drive (replace /dev/sdX with your 1TB drive)
sudo mkfs.ext4 -L storage /dev/sdX

# Format the 500GB drive (replace /dev/sdY with your 500GB drive)
sudo mkfs.ext4 -L backups /dev/sdY
```

### Step 6.3: Create Mount Points and Configure fstab
```bash
# Create mount points
sudo mkdir -p /srv/storage
sudo mkdir -p /srv/backups

# Get UUIDs
sudo blkid /dev/sdX  # 1TB drive
sudo blkid /dev/sdY  # 500GB drive

# Add to fstab for persistent mounts
echo 'UUID=<1tb-uuid-here> /srv/storage ext4 defaults 0 2' | sudo tee -a /etc/fstab
echo 'UUID=<500gb-uuid-here> /srv/backups ext4 defaults 0 2' | sudo tee -a /etc/fstab

# Mount now
sudo mount -a

# Verify
df -h /srv/storage /srv/backups
```

### Step 6.4: Create Directory Structure
```bash
# Main storage directories
sudo mkdir -p /srv/storage/{docker,shares}
sudo mkdir -p /srv/storage/shares/{public,private}

# Backup directories
sudo mkdir -p /srv/backups/{configs,shares,snapshots}

# Set ownership
sudo chown -R breen:breen /srv/storage
sudo chown -R breen:breen /srv/backups
```

### Step 6.5: Install and Configure Samba
```bash
# Install Samba
sudo apt install samba -y

# Backup original config
sudo cp /etc/samba/smb.conf /etc/samba/smb.conf.bak

# Create Samba password for breen
sudo smbpasswd -a breen
```

Edit Samba config:
```bash
sudo nano /etc/samba/smb.conf
```

Add at the end of the file:
```ini
[public]
   path = /srv/storage/shares/public
   browseable = yes
   read only = no
   guest ok = yes
   create mask = 0664
   directory mask = 0775

[private]
   path = /srv/storage/shares/private
   browseable = yes
   read only = no
   guest ok = no
   valid users = breen
   create mask = 0600
   directory mask = 0700
```

Apply:
```bash
# Restart Samba
sudo systemctl restart smbd nmbd
```

### Step 6.6: Access Shares
- **Windows:** `\\city17.local\public` or `\\city17.local\private`
- **macOS:** Finder → Go → Connect to Server → `smb://city17.local/public`
- **Linux:** `smb://city17.local/public` or mount via fstab

#### Linux Client Setup (for SMB access in file manager)

**Arch Linux:**
```bash
sudo pacman -S gvfs-smb
# Log out and back in, or restart file manager
thunar -q && thunar
```

**Ubuntu/Debian:**
```bash
sudo apt install gvfs-backends gvfs-fuse
# Log out and back in for changes to take effect
```

**Manual mount (any distro):**
```bash
sudo apt install cifs-utils  # or: sudo pacman -S cifs-utils
sudo mkdir -p /mnt/city17-public
sudo mount -t cifs //city17.local/public /mnt/city17-public -o guest
```

### Step 6.7: Set Up Automated Backups (rsync)
Create a backup script:
```bash
sudo nano /usr/local/bin/backup-storage.sh
```

Add the following content:
```bash
#!/bin/bash
# Backup critical data from /srv/storage to /srv/backups

TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)
LOG_FILE="/srv/backups/backup.log"

echo "[$TIMESTAMP] Starting backup..." >> "$LOG_FILE"

# Sync shares to backup drive (preserves permissions, deletes removed files)
rsync -av --delete /srv/storage/shares/ /srv/backups/shares/ >> "$LOG_FILE" 2>&1

# Backup important system configs
rsync -av /etc/samba/ /srv/backups/configs/samba/ >> "$LOG_FILE" 2>&1
rsync -av /etc/ssh/ /srv/backups/configs/ssh/ >> "$LOG_FILE" 2>&1
rsync -av /srv/storage/docker/ /srv/backups/configs/docker/ >> "$LOG_FILE" 2>&1

echo "[$TIMESTAMP] Backup completed." >> "$LOG_FILE"
```

Make executable and schedule:
```bash
# Make executable
sudo chmod +x /usr/local/bin/backup-storage.sh

# Add cron job for nightly backup at 3 AM
(crontab -l 2>/dev/null; echo "0 3 * * * /usr/local/bin/backup-storage.sh") | crontab -

# Verify cron job
crontab -l
```

**Manual backup anytime:**
```bash
/usr/local/bin/backup-storage.sh
```

---

## Phase 8: Maintenance

### Step 8.1: Enable Unattended Security Upgrades
```bash
# Usually installed by default, but ensure it's there
sudo apt install unattended-upgrades -y

# Enable automatic security updates
sudo dpkg-reconfigure -plow unattended-upgrades
# Select "Yes" when prompted
```

### Step 8.2: Verify Configuration
```bash
# Check status
sudo systemctl status unattended-upgrades

# View configuration
cat /etc/apt/apt.conf.d/20auto-upgrades
```

Expected content:
```
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
```

---

## Troubleshooting Heuristics

### Connection Issues
| Symptom | Check |
|---------|-------|
| Can't ping city17.local | Is Avahi running? `systemctl status avahi-daemon` |
| SSH refused | Service: `systemctl status ssh` |
| Can't access Portainer | Check if Portainer container is running: `docker ps` |
| Can't access Samba shares | Samba service: `systemctl status smbd` |

### Permission Issues
```bash
# Check file ownership
ls -la /path/to/file

# Check user groups
id breen

# Check if user in docker group
groups breen
```

### Container Issues
```bash
# View running containers
docker ps

# View all containers (including stopped)
docker ps -a

# View container logs
docker logs <container_name>

# Restart a container
docker restart <container_name>
```

### Storage/Backup Issues
```bash
# Check drive mounts
df -h /srv/storage /srv/backups

# Check if drives are healthy
sudo smartctl -a /dev/sda  # 1TB drive
sudo smartctl -a /dev/sdb  # 500GB drive

# View backup log
cat /srv/backups/backup.log

# Run backup manually
/usr/local/bin/backup-storage.sh

# Check cron jobs
crontab -l
```

---

## Quick Reference Commands

```bash
# SSH into city17
ssh breen@city17.local

# System updates
sudo apt update && sudo apt upgrade -y

# Reboot
sudo reboot

# Check disk space
df -h

# Check memory
free -h

# Check running services
systemctl list-units --type=service --state=running

# Docker status
docker ps

# Backup status
cat /srv/backups/backup.log | tail -10

# GPU status
nvidia-smi
```

---

## Future Expansion Ideas

Once city17 is stable, consider adding:

| Service | Purpose | Port |
|---------|---------|------|
| Pi-hole / AdGuard Home | Network-wide ad blocking + local DNS | 53, 80 |
| Tailscale | Secure remote access from anywhere | N/A (outbound) |
| Jellyfin / Plex | Media server | 8096 / 32400 |
| Home Assistant | Home automation | 8123 |
| Nextcloud | Self-hosted cloud storage | 443 |
| Gitea | Self-hosted Git server | 3000 |
| Proxmox (2nd machine) | Virtualization platform | 8006 |

---

## Second Server Planning (Future)

When ready to set up the second spare PC:
- Suggested hostname: `kleiner` or `lambda`
- Role options:
  - **Proxmox host:** Run VMs for learning/testing
  - **Dedicated NAS:** Offload storage from city17
  - **Security lab:** Isolated network for pentesting practice
