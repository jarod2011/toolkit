package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestWithWriter(t *testing.T) {
	options := Options{}
	assert.Nil(t, options.Writer)
	WithWriter(os.Stdout)(&options)
	assert.NotNil(t, options.Writer)
	WithWriter(nil)(&options)
	assert.Nil(t, options.Writer)
}

func TestWithLevel(t *testing.T) {
	options := Options{}
	assert.True(t, options.Level == Info)
	WithLevel(Warn)(&options)
	assert.True(t, options.Level == Warn)
}

func TestWithJsonFormat(t *testing.T) {
	options := Options{}
	assert.False(t, options.JsonFormat)
	WithJsonFormat(true)(&options)
	assert.True(t, options.JsonFormat)
}

func TestNewLogger(t *testing.T) {
	lg := NewLogger(WithLevel(Debug))
	assert.True(t, lg.l.Logger.Level == logrus.DebugLevel)
	assert.IsType(t, &logrus.JSONFormatter{}, lg.l.Logger.Formatter)
	lg = NewLogger(WithJsonFormat(false))
	assert.IsType(t, &logrus.TextFormatter{}, lg.l.Logger.Formatter)
}

func TestMyLogger_DebugF(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	assert.Equal(t, buf.Len(), 0)
	lg := NewLogger(WithLevel(Debug), WithWriter(buf))
	lg.DebugF("1111")
	assert.Greater(t, buf.Len(), 0)
	buf.Reset()
	lg.SetLevel(Info)
	lg.DebugF("1111")
	assert.Equal(t, buf.Len(), 0)
}

func TestMyLogger_InfoF(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	assert.Equal(t, buf.Len(), 0)
	lg := NewLogger(WithLevel(Info), WithWriter(buf))
	lg.InfoF("1111")
	assert.Greater(t, buf.Len(), 0)
	buf.Reset()
	lg.SetLevel(Warn)
	lg.InfoF("1111")
	assert.Equal(t, buf.Len(), 0)
}
func TestMyLogger_WarnF(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	assert.Equal(t, buf.Len(), 0)
	lg := NewLogger(WithLevel(Warn), WithWriter(buf))
	lg.WarnF("1111")
	assert.Greater(t, buf.Len(), 0)
	buf.Reset()
	lg.SetLevel(Error)
	lg.WarnF("1111")
	assert.Equal(t, buf.Len(), 0)
}
func TestMyLogger_ErrorF(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	assert.Equal(t, buf.Len(), 0)
	lg := NewLogger(WithLevel(Error), WithWriter(buf))
	lg.ErrorF("1111")
	assert.Greater(t, buf.Len(), 0)
	buf.Reset()
	lg.SetLevel(Fatal)
	lg.ErrorF("1111")
	assert.Equal(t, buf.Len(), 0)
}

func TestMyLogger_SetLevel(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	assert.Equal(t, buf.Len(), 0)
	lg := NewLogger(WithLevel(Info), WithWriter(buf))
	lg.InfoF("123")
	assert.Greater(t, buf.Len(), 0)
	buf.Reset()
	lg.SetLevel(Error)
	lg.InfoF("333")
	assert.Equal(t, buf.Len(), 0)
}

func TestMyLogger_WithField(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	assert.Equal(t, buf.Len(), 0)
	lg := NewLogger(WithWriter(buf))
	lg.WithField("name", "value").InfoF("13456")
	assert.True(t, bytes.Contains(buf.Bytes(), []byte("value")))
}
