package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractDOCX(t *testing.T) {
	path := filepath.Join(t.TempDir(), "resume.docx")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	archive := zip.NewWriter(file)
	document, err := archive.Create("word/document.xml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = document.Write([]byte(`<?xml version="1.0"?><w:document xmlns:w="x"><w:body><w:p><w:r><w:t>Ray Zhang</w:t></w:r></w:p><w:p><w:r><w:t>Go Developer</w:t></w:r></w:p></w:body></w:document>`))
	if err = archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err = file.Close(); err != nil {
		t.Fatal(err)
	}

	content, err := extractDOCX(path)
	if err != nil {
		t.Fatal(err)
	}
	if normalized := normalizeExtractedText(content); normalized != "Ray Zhang\nGo Developer" {
		t.Fatalf("unexpected content: %q", normalized)
	}
}

func TestNormalizeExtractedText(t *testing.T) {
	got := normalizeExtractedText(" Name \r\n\r\n\r\n Skills \x00 ")
	if got != "Name\n\nSkills" {
		t.Fatalf("unexpected normalized text: %q", got)
	}
}

func TestExtractUTF8Text(t *testing.T) {
	file := &memoryMultipartFile{Reader: strings.NewReader("Education\nComputer Science")}
	got, err := extractResumeText(file, ".txt")
	if err != nil || got != "Education\nComputer Science" {
		t.Fatalf("got %q, err %v", got, err)
	}
}

type memoryMultipartFile struct{ *strings.Reader }

func (file *memoryMultipartFile) Close() error { return nil }
