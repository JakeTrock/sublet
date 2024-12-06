FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu

# Install basic requirements and build tools for node-sqlite3
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    xz-utils \
    sudo \
    git \
    python3 \
    make \
    g++ \
    sqlite3 \
    libsqlite3-dev

# Create the vscode user
# RUN useradd -m -s /bin/bash vscode \
#     && usermod -aG sudo vscode \
#     && echo "vscode ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/vscode

# Switch to vscode user for remaining operations, make sure node_modules are not copied to the container
USER vscode
WORKDIR /home/vscode

# Install Nix
RUN curl -L https://nixos.org/nix/install | sh -s -- --no-daemon

# Add Nix to PATH and configure
ENV USER=vscode \
    PATH="/home/vscode/.nix-profile/bin:${PATH}"

# Enable flakes
RUN mkdir -p ~/.config/nix \
    && echo "experimental-features = nix-command flakes" >> ~/.config/nix/nix.conf