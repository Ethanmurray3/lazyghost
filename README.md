# lazyghost

A small Bubble Tea TUI for switching Ghostty custom shaders.

## Usage

```sh
go run .
```

Use `j`/`down` and `k`/`up` to move, then press `enter` or `space` to apply a shader.

lazyghost embeds the shader files from `shaders/` into the binary. When you apply a shader, it installs that shader to:

```text
~/.config/ghostty/lazyghost-shaders/
```

Then it updates:

```text
~/.config/ghostty/config.ghostty
```

with a `custom-shader = ...` line and asks Ghostty to reload.

## Install From GitHub

Recommended:

```sh
curl -fsSL https://raw.githubusercontent.com/Ethanmurray3/lazyghost/main/install.sh | sh
```

This downloads the latest release binary for your OS/CPU and installs it to:

```text
~/.local/bin/lazyghost
```

Then run:

```sh
lazyghost
```

If `~/.local/bin` is not on your `PATH`, add it:

```sh
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Install With Go

```sh
go install github.com/Ethanmurray3/lazyghost@latest
```

Go downloads the source, builds it for your computer, and installs the binary to:

```sh
$(go env GOPATH)/bin
```

For most setups, that is:

```text
~/go/bin/lazyghost
```

Then run:

```sh
~/go/bin/lazyghost
```

To run it as `lazyghost` from anywhere, add Go's bin directory to your `PATH`:

```sh
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

Then run:

```sh
lazyghost
```

To update later:

```sh
go install github.com/Ethanmurray3/lazyghost@latest
```
