package models

import (
	"io"
	"time"
)

type RequestAttachment struct {
	ID        int64
	RequestID int64
	FileName  string
	FileURL   string
	FileSize  int64
	CreatedAt time.Time
}

type UploadAttachmentInput struct {
	UserID    int64
	UserRole  string
	RequestID int64
	FileName  string
	FileSize  int64
	Content   io.Reader
}
