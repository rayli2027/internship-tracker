package main

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/dslipak/pdf"
	"github.com/gin-gonic/gin"
)

const maxResumeUpload = 10 << 20

func (app *App) extractResumeHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxResumeUpload+(1<<20))
	header, err := c.FormFile("file")
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "resume file must be 10 MB or smaller"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "select a resume file"})
		return
	}
	if header.Size > maxResumeUpload {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "resume file must be 10 MB or smaller"})
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" && ext != ".docx" && ext != ".txt" && ext != ".md" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "supported formats are PDF, DOCX, TXT, and MD"})
		return
	}
	file, err := header.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to open resume file"})
		return
	}
	defer file.Close()
	content, err := extractResumeText(file, ext)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	content = normalizeExtractedText(content)
	if len([]rune(content)) < 30 {
		message := "no readable text was found"
		if ext == ".pdf" {
			message = "no readable text was found; this PDF may be scanned or password-protected"
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": content, "source_name": filepath.Base(header.Filename)})
}

func extractResumeText(file multipart.File, ext string) (string, error) {
	if ext == ".txt" || ext == ".md" {
		raw, err := io.ReadAll(io.LimitReader(file, maxResumeUpload+1))
		if err != nil {
			return "", fmt.Errorf("unable to read text file")
		}
		if !utf8.Valid(raw) {
			return "", fmt.Errorf("text file must use UTF-8 encoding")
		}
		return string(raw), nil
	}
	temp, err := os.CreateTemp("", "internship-resume-*"+ext)
	if err != nil {
		return "", fmt.Errorf("unable to process resume")
	}
	path := temp.Name()
	defer os.Remove(path)
	if _, err = io.Copy(temp, io.LimitReader(file, maxResumeUpload+1)); err != nil {
		temp.Close()
		return "", fmt.Errorf("unable to process resume")
	}
	if err = temp.Close(); err != nil {
		return "", fmt.Errorf("unable to process resume")
	}
	if ext == ".pdf" {
		return extractPDF(path)
	}
	return extractDOCX(path)
}

func extractPDF(path string) (string, error) {
	reader, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("unable to read PDF; it may be damaged or password-protected")
	}
	plain, err := reader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("unable to extract PDF text; it may be scanned or password-protected")
	}
	raw, err := io.ReadAll(io.LimitReader(plain, 2<<20))
	if err != nil {
		return "", fmt.Errorf("unable to extract PDF text")
	}
	return string(raw), nil
}

func extractDOCX(path string) (string, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return "", fmt.Errorf("unable to read DOCX; the file may be damaged")
	}
	defer archive.Close()
	for _, entry := range archive.File {
		if entry.Name != "word/document.xml" {
			continue
		}
		stream, err := entry.Open()
		if err != nil {
			return "", fmt.Errorf("unable to read DOCX content")
		}
		defer stream.Close()
		decoder := xml.NewDecoder(stream)
		var output strings.Builder
		for {
			token, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", fmt.Errorf("unable to parse DOCX content")
			}
			switch value := token.(type) {
			case xml.StartElement:
				switch value.Name.Local {
				case "t":
					var text string
					if decoder.DecodeElement(&text, &value) == nil {
						output.WriteString(text)
					}
				case "tab":
					output.WriteByte('\t')
				case "br", "cr":
					output.WriteByte('\n')
				}
			case xml.EndElement:
				if value.Name.Local == "p" {
					output.WriteByte('\n')
				}
			}
		}
		return output.String(), nil
	}
	return "", fmt.Errorf("DOCX does not contain a readable document")
}

func normalizeExtractedText(value string) string {
	value = strings.ReplaceAll(value, "\x00", "")
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	lines := strings.Split(value, "\n")
	cleaned := make([]string, 0, len(lines))
	blank := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if !blank && len(cleaned) > 0 {
				cleaned = append(cleaned, "")
			}
			blank = true
			continue
		}
		cleaned = append(cleaned, line)
		blank = false
	}
	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}
