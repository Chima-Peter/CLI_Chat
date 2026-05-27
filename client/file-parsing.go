package client

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const maxZipSourceSize = 100 << 20 // 100 MiB

// blockedExtensions are never zipped (security / policy).
var blockedExtensions = map[string]string{
	".bat": "batch scripts are not allowed",
	".cmd": "command scripts are not allowed",
	".com": "executables are not allowed",
	".cpl": "control panel applets are not allowed",
	".dll": "executables are not allowed",
	".exe": "executables are not allowed",
	".msc": "Microsoft Management Console files are not allowed",
	".msi": "installers are not allowed",
	".ps1": "PowerShell scripts are not allowed",
	".scr": "screensaver executables are not allowed",
	".vbs": "VBScript files are not allowed",
}

// unsuitableExtensions are already compressed or do not deflate meaningfully.
var unsuitableExtensions = map[string]string{
	".7z":   "archive is already compressed",
	".aac":  "audio is already compressed",
	".avif": "image is already compressed",
	".avi":  "video is already compressed",
	".bz2":  "archive is already compressed",
	".docx": "Office document is already compressed (ZIP-based)",
	".flac": "audio is already compressed",
	".gif":  "image is already compressed",
	".gz":   "archive is already compressed",
	".heic": "image is already compressed",
	".heif": "image is already compressed",
	".ico":  "image is already compressed",
	".jpeg": "image is already compressed",
	".jpg":  "image is already compressed",
	".m4a":  "audio is already compressed",
	".m4v":  "video is already compressed",
	".mkv":  "video is already compressed",
	".mov":  "video is already compressed",
	".mp3":  "audio is already compressed",
	".mp4":  "video is already compressed",
	".ogg":  "media is already compressed",
	".opus": "audio is already compressed",
	".pdf":  "PDF is already compressed",
	".png":  "image is already compressed",
	".pptx": "Office document is already compressed (ZIP-based)",
	".rar":  "archive is already compressed",
	".svg":  "SVG is typically not worth re-compressing in ZIP",
	".tar":  "archive containers are not meaningfully compressed again",
	".tgz":  "archive is already compressed",
	".tif":  "image is already compressed",
	".tiff": "image is already compressed",
	".wasm": "WebAssembly binaries are already compressed",
	".webm": "video is already compressed",
	".webp": "image is already compressed",
	".wmv":  "video is already compressed",
	".xlsx": "Office document is already compressed (ZIP-based)",
	".xz":   "archive is already compressed",
	".zip":  "archive is already compressed",
	".zst":  "archive is already compressed",
}

func ZipFile(filename string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find your Downloads folder: %w", err)
	}
	downloadsPath := filepath.Join(home, "Downloads")
	safeFileName := filepath.Base(filename)
	filePath := filepath.Join(downloadsPath, safeFileName)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("could not find %q in your Downloads folder", safeFileName)
		}
		return "", fmt.Errorf("could not access %q in Downloads: %w", safeFileName, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("cannot zip %s: folders are not supported", safeFileName)
	}

	ext := strings.ToLower(filepath.Ext(info.Name()))
	if reason, blocked := blockedExtensions[ext]; blocked {
		return "", fmt.Errorf("cannot zip %s: %s", safeFileName, reason)
	}
	if reason, unsuitable := unsuitableExtensions[ext]; unsuitable {
		return "", fmt.Errorf("cannot zip %s: %s (not suitable for compression)", safeFileName, reason)
	}
	if info.Size() > maxZipSourceSize {
		return "", fmt.Errorf("cannot zip %s: file exceeds 100 MB limit (%.1f MB)",
			safeFileName, float64(info.Size())/(1<<20))
	}

	base := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
	zipPath := filepath.Join(downloadsPath, base+".zip")

	out, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("create %q: %w", zipPath, err)
	}

	zipWriter := zip.NewWriter(out)

	sourceFile, err := os.Open(filePath)
	if err != nil {
		_ = out.Close()
		return "", fmt.Errorf("open %q: %w", filePath, err)
	}
	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		_ = out.Close()
		return "", fmt.Errorf("stat open file %q: %w", filePath, err)
	}

	header, err := zip.FileInfoHeader(sourceInfo)
	if err != nil {
		_ = out.Close()
		return "", fmt.Errorf("build zip header for %q: %w", info.Name(), err)
	}
	header.Name = info.Name()
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		_ = out.Close()
		return "", fmt.Errorf("create zip entry for %q: %w", info.Name(), err)
	}

	if _, err = io.Copy(writer, sourceFile); err != nil {
		_ = zipWriter.Close()
		_ = out.Close()
		return "", fmt.Errorf("write %q into archive: %w", info.Name(), err)
	}

	if err := zipWriter.Close(); err != nil {
		_ = out.Close()
		return "", fmt.Errorf("finalize archive %q: %w", zipPath, err)
	}
	if err := out.Close(); err != nil {
		return "", fmt.Errorf("close %q: %w", zipPath, err)
	}

	return zipPath, nil
}
