# imanage

**imanage** is a tool for managing your photos.

## Installation

1 . Install `libaom-dev` (imanage transcoding)

```bash
$ sudo apt-get install libaom-dev
```

2 . Install libwebp

```bash
$ brew install webp
```

3 . Install libavif

```bash
$ brew install libavif
```

## Availability

|       | to_webp | to_avif | to_webp(metadata) | to_avif(metadata) |
| ----- | ------- | ------- | --- | --- |
|imanage| O       |  O      | X   | X   |
|libwebp| O       |  X      | O   | X   |
|libavif| X       |  O      | X   | O   |
