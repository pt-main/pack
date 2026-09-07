# Pack — file orchestrator

```bash
go install github.com/pt-main/pack/cmd/pack@latest
```

A minimalistic file archiver that works in the CLI.  
It packs files into a single `.pack` file with optional compression and encryption.

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

**Show the structure:**
```bash
pack struct archive.pack
```

**Show file contents:**
```bash
pack show archive.pack myscript
```

### Working with directories

**Pack a directory:**
```bash
pack files zip archive.pack ./myproject
```
All files from `./myproject` are added to the archive with relative paths as container names.

**Unpack a directory:**
```bash
pack files unzip archive.pack ./output
```
Extracts all files from the archive into `./output`, preserving the folder structure.

## Flags

- `--secure` – AES‑GCM encryption
- `--zip` – gzip compression
- `--key="..."` – encryption key (16/24/32 bytes)

Flags specified at creation time must be used for all subsequent operations.

## Format

```
PACK@0
  -> name (once)
  -> metadata
  -> containers: name -> data
```

## License

Apache 2.0

---

By Pt, 2026, using Lc and Tap.