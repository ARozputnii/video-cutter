# 🎬 video-cutter

A simple local web UI to upload, preview, trim multiple segments of a video, and merge them into one.  
No need to install Go or FFmpeg manually — everything runs out of the box.

---

## 💡 Features

- Web interface to cut video into multiple ranges
- FFmpeg-based trimming and merging
- Works fully offline
- No external dependencies required (FFmpeg included)
- Cross-platform: Windows & Linux

---

## 🚀 Quick Start
---

### 🪟 Windows

First run
```bash
  double-click windows-launcher.bat
```
Then open your browser and navigate to:

```bash
  http://localhost:3000
```
For next runs, just double-click [video-cutter.exe](video-cutter.exe)
---

### 🧑‍💻 Developer Setup

Prerequisites
 - Go 1.24+ installed 
 - FFmpeg available in system PATH

```bash
    go mod tidy
    go run cmd/main.go
```
---

## ⚠️ Notes

- ✂️ Final output file overwrites the original
- 🧹 Temporary files are cleaned up automatically
- 🎞 Works best with `.mp4` input files
---

## 🧾 License

MIT © 2025
