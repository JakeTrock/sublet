# 2024-12-01

## 1. this is a beta program to install nixos and reconf it live using LLMS

### preferred install

```
curl https://raw.githubusercontent.com/JakeTrock/sublet/refs/heads/master/installer/infect.sh | PROVIDER=digitalocean NIX_CHANNEL=nixos-23.05 bash 2>&1 | tee /tmp/infect.log
```

or cloudconfig(only tested on digitalocean)

```
#cloud-config

runcmd:
  - curl https://raw.githubusercontent.com/JakeTrock/sublet/refs/heads/master/installer/infect.sh | PROVIDER=digitalocean NIX_CHANNEL=nixos-23.05 bash 2>&1 | tee /tmp/infect.log
```

### beta install do not use yet

```
curl https://raw.githubusercontent.com/JakeTrock/sublet/refs/heads/master/installer/installer/infect_v2_beta.sh | PROVIDER=digitalocean NIX_CHANNEL=nixos-24.11 bash 2>&1 | tee /tmp/infect.log
```

or cloudconfig(only tested on digitalocean)

```
#cloud-config

runcmd:
  - curl https://raw.githubusercontent.com/JakeTrock/sublet/refs/heads/master/installer/infect_v2_beta.sh | PROVIDER=digitalocean NIX_CHANNEL=nixos-24.11 bash 2>&1 | tee /tmp/infect.log
```
