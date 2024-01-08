# Takway

This is simple demo for MAIXSENSE to interact with OpenAI VLM model.

### board

1. Prerequest
Upload the script to board.

- python 3.8

2. Run

```bash
cd /board
python run.py
```

### server

1. Prerequest

- golang 1.18
- ffmpeg

2. Run

```bash
cd /server
go mod tidy
go run main.go
```