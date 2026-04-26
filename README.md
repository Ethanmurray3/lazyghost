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

After this repo is pushed to GitHub:

```sh
go install github.com/Ethanmurray3/lazyghost@latest
```

Then run:

```sh
lazyghost
```
