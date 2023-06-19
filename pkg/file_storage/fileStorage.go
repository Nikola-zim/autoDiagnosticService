package file_storage

import "time"

type FileStorage struct {
	RecognitionInterval time.Duration
}

const (
	_defaultRecognitionInterval = time.Duration(10 * time.Second)
)

func New(url string, opts ...Option) (*FileStorage, error) {
	fs := &FileStorage{
		RecognitionInterval: _defaultRecognitionInterval,
	}

	// Custom options
	for _, opt := range opts {
		opt(fs)
	}
	return fs, nil
}
