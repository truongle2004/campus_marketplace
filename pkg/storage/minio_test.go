package storage

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	tcminio "github.com/testcontainers/testcontainers-go/modules/minio"
)

const testBucket = "test-bucket"

func setupMinio(t *testing.T) Service {
	t.Helper()
	ctx := context.Background()

	container, err := tcminio.Run(ctx, "minio/minio:latest")
	if err != nil {
		t.Fatalf("start minio container: %v", err)
	}
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	endpoint, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	user := container.Username
	pass := container.Password

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(user, pass, ""),
		Secure: false,
	})
	if err != nil {
		t.Fatalf("create minio client: %v", err)
	}

	if err := client.MakeBucket(ctx, testBucket, minio.MakeBucketOptions{}); err != nil {
		t.Fatalf("create bucket: %v", err)
	}

	return &service{client: client, bucket: testBucket}
}

func TestUploadFile(t *testing.T) {
	svc := setupMinio(t)
	ctx := context.Background()

	body := "hello world"
	path, err := svc.UploadFile(ctx, "test.txt", strings.NewReader(body), int64(len(body)), "text/plain")
	if err != nil {
		t.Fatalf("upload: %v", err)
	}
	if path != "/test-bucket/test.txt" {
		t.Fatalf("unexpected path: %s", path)
	}
}

func TestGetFileURL(t *testing.T) {
	svc := setupMinio(t)
	ctx := context.Background()

	body := "hello"
	if _, err := svc.UploadFile(ctx, "url-test.txt", strings.NewReader(body), int64(len(body)), "text/plain"); err != nil {
		t.Fatalf("upload: %v", err)
	}

	u, err := svc.GetFileURL(ctx, "url-test.txt", 5*time.Minute)
	if err != nil {
		t.Fatalf("get url: %v", err)
	}
	if u == "" {
		t.Fatal("expected non-empty url")
	}
	if !strings.Contains(u, "url-test.txt") {
		t.Fatalf("url should contain object name: %s", u)
	}
}

func TestDeleteFile(t *testing.T) {
	svc := setupMinio(t)
	ctx := context.Background()

	body := "to-delete"
	if _, err := svc.UploadFile(ctx, "del.txt", strings.NewReader(body), int64(len(body)), "text/plain"); err != nil {
		t.Fatalf("upload: %v", err)
	}

	if err := svc.DeleteFile(ctx, "del.txt"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	s := svc.(*service)
	_, err := s.client.StatObject(ctx, testBucket, "del.txt", minio.StatObjectOptions{})
	if err == nil {
		t.Fatal("expected error after delete, object still exists")
	}
}

func TestUploadFile_EmptyContent(t *testing.T) {
	svc := setupMinio(t)
	ctx := context.Background()

	path, err := svc.UploadFile(ctx, "empty.txt", strings.NewReader(""), 0, "text/plain")
	if err != nil {
		t.Fatalf("upload empty: %v", err)
	}
	if path != "/test-bucket/empty.txt" {
		t.Fatalf("unexpected path: %s", path)
	}
}
