package audio

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestResolver(t *testing.T, dir string) Resolver {
	t.Helper()

	r, err := NewResolver(dir)
	if err != nil {
		t.Fatalf("NewResolver: %v", err)
	}
	return r
}

func TestResolver_AudioFilePath_EmptyDir(t *testing.T) {
	r := newTestResolver(t, "")

	if path := r.AudioFilePath("10", "10100001"); path != "" {
		t.Errorf("AudioFilePath with empty dir: got %q, want empty", path)
	}
}

func TestResolver_AudioFilePath_NonexistentFile(t *testing.T) {
	r := newTestResolver(t, "/nonexistent/audio/dir")

	if path := r.AudioFilePath("10", "10100001"); path != "" {
		t.Errorf("AudioFilePath with nonexistent file: got %q, want empty", path)
	}
}

func TestResolver_AudioFilePath_SubtitleSuffixStrip(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "00")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "end_all00.ogg"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	r := newTestResolver(t, dir)
	expected := filepath.Join(dir, "00", "end_all00.ogg")

	if path := r.AudioFilePath("00", "end_all00_s0"); path != expected {
		t.Errorf("AudioFilePath _s0: got %q, want %q", path, expected)
	}
	if path := r.AudioFilePath("00", "end_all00_s15"); path != expected {
		t.Errorf("AudioFilePath _s15: got %q, want %q", path, expected)
	}
}

func TestResolver_AudioFilePath_SubtitleSuffixNoBase(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "00"), 0755); err != nil {
		t.Fatal(err)
	}

	r := newTestResolver(t, dir)

	if path := r.AudioFilePath("00", "end_all00_s5"); path != "" {
		t.Errorf("AudioFilePath should return empty when base file also missing: got %q", path)
	}
}

func TestResolver_AudioFilePath_ExactMatchPreferred(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "00")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "end_all00_s0.ogg"), []byte("exact"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "end_all00.ogg"), []byte("base"), 0644); err != nil {
		t.Fatal(err)
	}

	r := newTestResolver(t, dir)
	expected := filepath.Join(dir, "00", "end_all00_s0.ogg")

	if path := r.AudioFilePath("00", "end_all00_s0"); path != expected {
		t.Errorf("AudioFilePath should prefer exact match: got %q, want %q", path, expected)
	}
}

func TestResolver_HasAudio_EmptyDir(t *testing.T) {
	if newTestResolver(t, "").HasAudio() {
		t.Error("HasAudio with empty dir string: expected false")
	}
}

func TestResolver_HasAudio_NonexistentDir(t *testing.T) {
	if newTestResolver(t, "/nonexistent/audio/dir").HasAudio() {
		t.Error("HasAudio with nonexistent dir: expected false")
	}
}

func TestResolver_HasAudio_EmptyExistingDir(t *testing.T) {
	if newTestResolver(t, t.TempDir()).HasAudio() {
		t.Error("HasAudio with empty existing dir: expected false")
	}
}

func TestResolver_HasAudio_WithFiles(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "10")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "10100001.ogg"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	if !newTestResolver(t, dir).HasAudio() {
		t.Error("HasAudio with files: expected true")
	}
}
