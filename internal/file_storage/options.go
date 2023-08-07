package file_storage

import "time"

// Option -.
type Option func(storage *FileStorage)

// MaxPoolSize -.
func RecognitionInterval(interval time.Duration) Option {
	return func(c *FileStorage) {
		c.RecognitionInterval = interval
	}
}
