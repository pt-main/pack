# Pack — file orchestrator

```bash
go install github.com/pt-main/cmd/pack@latest
```

Packs files into a single `.pack` file with optional compression and encryption.

## Commands

**Create an archive:**
```bash
pack create archive.pack
pack create input.txt archive.pack data
```

**Add files:**
```bash
pack push archive.pack --script.py=myscript --config.json=cfg
```

**Extract a file:**
```bash
pack get archive.pack myscript out.py
```

**Show structure:**
```bash
pack struct archive.pack
```

**Show contents:**
```bash
pack show archive.pack myscript
```

## Flags

- `--secure` — AES-GCM encryption
- `--zip` — gzip compression
- `--key="..."` — encryption key (16/24/32 bytes)

Flags specified during creation must be used for all operations.

## Format

```
PACK@0
  → name (once)
  → metadata
  → containers: name → data
```

## License

Apache 2.0

---

By Pt, 2026, using Lc and Tap.