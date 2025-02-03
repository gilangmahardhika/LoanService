package helpers

// DetectFileExtension determines the file extension based on the file's magic numbers
func DetectFileExtension(fileBytes []byte) string {
	switch {
	case len(fileBytes) > 2 && fileBytes[0] == 0xFF && fileBytes[1] == 0xD8:
		return ".jpg"
	case len(fileBytes) > 8 && string(fileBytes[0:8]) == "\x89PNG\r\n\x1a\n":
		return ".png"
	case len(fileBytes) > 4 && string(fileBytes[0:4]) == "RIFF":
		return ".webp"
	default:
		return ".bin" // fallback for unknown formats
	}
}
