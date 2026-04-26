#!/bin/sh
set -eu

repo="Ethanmurray3/lazyghost"
bin_name="lazyghost"
install_dir="${INSTALL_DIR:-$HOME/.local/bin}"

os="$(uname -s)"
arch="$(uname -m)"

case "$os" in
  Darwin) os="darwin" ;;
  Linux) os="linux" ;;
  *)
    echo "Unsupported OS: $os" >&2
    exit 1
    ;;
esac

case "$arch" in
  arm64|aarch64) arch="arm64" ;;
  x86_64|amd64) arch="amd64" ;;
  *)
    echo "Unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

asset="${bin_name}-${os}-${arch}"
url="https://github.com/${repo}/releases/latest/download/${asset}"
target="${install_dir}/${bin_name}"

mkdir -p "$install_dir"

echo "Installing ${asset} to ${target}"
curl -fsSL "$url" -o "$target"
chmod +x "$target"

case ":$PATH:" in
  *":$install_dir:"*)
    echo "Installed. Run: ${bin_name}"
    ;;
  *)
    echo "Installed to ${target}"
    echo "Add this to your shell config to run '${bin_name}' from anywhere:"
    echo "  export PATH=\"${install_dir}:\$PATH\""
    ;;
esac
