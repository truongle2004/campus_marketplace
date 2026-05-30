package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewWithDefaultConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := DefaultConfig()
	cfg.LogDir = dir

	l, err := New(cfg)
	require.NoError(t, err)

	l.Info("boot")
	_ = l.Sync()

	data, err := os.ReadFile(filepath.Join(dir, cfg.LogFile))
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestNewWithoutFile(t *testing.T) {
	cfg := DefaultConfig()
	cfg.WithFile = false

	l, err := New(cfg)
	require.NoError(t, err)
	assert.NotNil(t, l)

	l.Info("no file")
	_ = l.Sync()
}

func TestNewWithInvalidLogDir(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LogDir = "/dev/null/invalid"

	l, err := New(cfg)
	assert.Error(t, err)
	assert.Nil(t, l)
}

func TestNewWithInvalidLogFile(t *testing.T) {
	dir := t.TempDir()
	cfg := DefaultConfig()
	cfg.LogDir = dir
	cfg.LogFile = "sub/dir"
	require.NoError(t, os.MkdirAll(filepath.Join(dir, cfg.LogFile), 0755))

	l, err := New(cfg)
	assert.Error(t, err)
	assert.Nil(t, l)
}

func TestLogInfo(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	logger := NewWithCore(core)

	logger.Info("Test", zap.String("id", "123"))

	require.Equal(t, 1, logs.Len())

	entry := logs.All()[0]
	assert.Equal(t, "Test", entry.Message)
	assert.Equal(t, "id", entry.Context[0].Key)
	assert.Equal(t, "123", entry.Context[0].String)
}

func TestDebugFilteredAtInfoLevel(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	logger := NewWithCore(core)

	logger.Debug("should be filtered")

	assert.Equal(t, 0, logs.Len())
}

func TestErrorLoggedAtInfoLevel(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	logger := NewWithCore(core)

	logger.Error("error happened")

	require.Equal(t, 1, logs.Len())
	assert.Equal(t, "error happened", logs.All()[0].Message)
}

func TestWarnLoggedAtInfoLevel(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	logger := NewWithCore(core)

	logger.Warn("warning")

	require.Equal(t, 1, logs.Len())
	assert.Equal(t, "warning", logs.All()[0].Message)
}

func TestNewNop(t *testing.T) {
	l := NewNop()
	assert.NotNil(t, l)
}

func TestSugar(t *testing.T) {
	l := NewNop()
	assert.NotNil(t, l.Sugar())
}
