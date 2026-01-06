
# Master Home Lab Setup Guide: Ubuntu Server 24.04 (Headless)

## Project Overview
Goal: Transform a spare PC into a robust, headless Linux server for hosting services, file sharing, and acting as a central lab node.  
OS: Ubuntu Server 24.04 LTS[web:141][web:153]  
Architecture: Headless (SSH only), Docker-based services, secure by design.

Naming Theme: Use Valve’s The Orange Box (Half‑Life 2, Portal, TF2) for hostnames and internal domains, e.g.:
- glados.aperture.lab (Portal)
- citadel.combine.lab (Half‑Life 2)
- manco.manco.lab (TF2)[web:197][web:199][web:212]

---

## Phase 1: Installation & Base Config

### 1. Installation Method

Option A – Interactive install (easier first time)
1. Download Ubuntu Server 24.04 ISO.[web:141][web:153]
2. Create bootable USB (Rufus on Windows or dd on Linux).[web:161]
3. Connect a monitor/keyboard temporarily.
4. During install:
   - Enable OpenSSH Server.
   - Set a static IP (e.g., 192.168.1.100/24).
   - Use default LVM for storage.
5. Reboot, then remove monitor/keyboard and use SSH only.

Option B – Cloud-init Autoinstall (advanced)
1. Create user-data with autoinstall YAML for Ubuntu 24.04.[web:142][web:68]
2. Include:
   - Your SSH public key.
   - Static IP configuration.
   - User admin and hostname (e.g., glados).
3. Boot with kernel parameter like:
   autoinstall ds=nocloud;s=/cdrom/nocloud/
4. Wait for unattended install, then SSH in.

### 2. First Boot Checklist

ssh admin@192.168.1.100
sudo apt update && sudo apt upgrade -y
sudo apt install -y htop git curl vim net-tools

---

## Phase 2: Security Hardening

### 1. SSH Hardening

Edit /etc/ssh/sshd_config:[web:101][web:104]

PasswordAuthentication no
PubkeyAuthentication yes
PermitRootLogin no
PermitEmptyPasswords no

sudo systemctl restart ssh

Test SSH again from a new terminal using your key.

### 2. Firewall (UFW)

sudo ufw allow ssh
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw enable
sudo ufw status verbose
[web:103][web:109][web:178][web:181]

### 3. Fail2Ban

sudo apt install fail2ban -y
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
sudo nano /etc/fail2ban/jail.local

In [sshd]:[web:111][web:114]

enabled = true
port = ssh
maxretry = 3
bantime = 3600
findtime = 600

sudo systemctl restart fail2ban
sudo fail2ban-client status sshd

---

## Phase 3: Docker Stack

### 1. Install Docker Engine

Use the official Docker repository on Ubuntu 24.04:[web:171][web:169]

sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo systemctl enable --now docker
sudo usermod -aG docker $USER
newgrp docker
docker run hello-world

### 2. Portainer (Docker GUI)

docker volume create portainer_data

docker run -d \
  -p 9443:9443 \
  --name portainer \
  --restart=always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v portainer_data:/data \
  portainer/portainer-ce:latest
[web:25][web:22]

sudo ufw allow 9443/tcp

Access: https://<hostname>:9443

---

## Phase 4: File Sharing (NAS)

### 1. Base Directories

sudo mkdir -p /srv/nas/public
sudo mkdir -p /srv/nas/media
sudo chown -R nobody:nogroup /srv/nas/public

### 2. Samba (Windows/macOS)

sudo apt install samba -y
sudo nano /etc/samba/smb.conf

Add:[web:23][web:26]

[public]
   path = /srv/nas/public
   browseable = yes
   read only = no
   guest ok = yes

sudo systemctl restart smbd
sudo ufw allow samba

Access from Windows: \\<hostname>\public

### 3. NFS (Linux clients)

sudo apt install nfs-kernel-server -y
sudo mkdir -p /srv/nfs/shared
sudo chown nobody:nogroup /srv/nfs/shared
sudo chmod 777 /srv/nfs/shared
sudo nano /etc/exports

Add:[web:23][web:32]
/srv/nfs/shared 192.168.1.0/24(rw,sync,no_subtree_check,no_root_squash)

sudo exportfs -ra
sudo systemctl restart nfs-kernel-server
sudo ufw allow from 192.168.1.0/24 to any port nfs

---

## Phase 5: Monitoring

### 1. Netdata

sudo apt install netdata -y
sudo systemctl enable --now netdata
sudo ufw allow from 192.168.1.0/24 to any port 19999

Access: http://<hostname>:19999[web:83][web:86]

### 2. Prometheus + Grafana (Docker)

Create ~/monitoring/docker-compose.yml:

version: '3.8'
services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    restart: unless-stopped
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    restart: unless-stopped
    ports:
      - "3000:3000"
    volumes:
      - grafana-data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
volumes:
  prometheus-data:
  grafana-data:

Create ~/monitoring/prometheus/prometheus.yml:[web:83][web:95]

global:
  scrape_interval: 15s
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'netdata'
    metrics_path: '/api/v1/allmetrics'
    params:
      format: [prometheus]
    static_configs:
      - targets: ['<hostname>:19999']

cd ~/monitoring
docker compose up -d
sudo ufw allow 9090/tcp
sudo ufw allow 3000/tcp

Access:
- Prometheus: http://<hostname>:9090
- Grafana: http://<hostname>:3000

---

## Phase 6: Backups (Restic)

sudo apt install restic -y
sudo mkdir -p /mnt/backups/restic-repo
sudo chown $USER:$USER /mnt/backups
restic init --repo /mnt/backups/restic-repo
[web:84][web:90]

Create ~/backup.sh:

#!/bin/bash
REPO="/mnt/backups/restic-repo"
PASSWORD_FILE="$HOME/.restic-password"
restic -r "$REPO" --password-file "$PASSWORD_FILE" \
  backup /home /etc /var/lib/docker/volumes \
  --exclude=/home/*/.cache
restic -r "$REPO" --password-file "$PASSWORD_FILE" \
  forget --keep-daily 7 --keep-weekly 4 --keep-monthly 6 --prune

echo "YOUR_RESTIC_PASSWORD" > ~/.restic-password
chmod 600 ~/.restic-password
chmod +x ~/backup.sh
crontab -e

Add:
0 2 * * * /home/admin/backup.sh >> /var/log/restic-backup.log 2>&1

---

## 🚀 Future Project Ideas (Orange Box Themed)

- DNS / Ad-blocking: Pi-hole → pihole.aperture.lab or pihole.manco.lab.
- Media Server: Jellyfin/Plex → alyx.whiteforest.lab or heavy.manco.lab.[web:199][web:211]
- Personal Cloud: Nextcloud → blackmesa.combine.lab or companioncube.aperture.lab.[web:197][web:199]
- Reverse Proxy: Nginx Proxy Manager → gate.aperture.lab or checkpoint.combine.lab.
- Remote Access: Tailscale on host, exposing only internal services over VPN.

---

## 🛠️ CLI Quick Reference

ssh admin@<hostname>
sudo apt update && sudo apt upgrade -y
docker ps
docker logs -f <container>
docker restart <container>
sudo ufw status
sudo ufw allow <port>/tcp
sudo systemctl status <service>
sudo systemctl restart <service>
