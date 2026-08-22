# Pack — файловый оркестратор

```bash
go install github.com/pt-main/pack/cmd/pack@latest
```

Упаковывает файлы в один `.pack` файл с опциональным сжатием и шифрованием.

## Команды

**Создать архив:**
```bash
pack create archive.pack
pack create input.txt archive.pack data
```

**Добавить файлы:**
```bash
pack push archive.pack --script.py=myscript --config.json=cfg
```

**Извлечь файл:**
```bash
pack get archive.pack myscript out.py
```

**Показать структуру:**
```bash
pack struct archive.pack
```

**Показать содержимое:**
```bash
pack show archive.pack myscript
```

## Флаги

- `--secure` — шифрование AES-GCM
- `--zip` — сжатие gzip
- `--key="..."` — ключ шифрования (16/24/32 байта)

Флаги, указанные при создании, нужно использовать при всех операциях.

## Формат

```
PACK@0
  → имя (один раз)
  → метаданные
  → контейнеры: имя → данные
```

## Лицензия

Apache 2.0

---

By Pt, 2026, using Lc and Tap.