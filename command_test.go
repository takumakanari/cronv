package cronv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCronvCommandDefaultParameters(t *testing.T) {
	cmd := NewCronvCommand()
	assert.NotNil(t, cmd)
	assert.NotNil(t, cmd.FromDate)
	assert.NotNil(t, cmd.FromTime)
	assert.Equal(t, cmd.OutputFilePath, optDefaultOutputPath, "")
	assert.Equal(t, cmd.Duration, optDefaultDuration, "")
}

func TestToDurationMinutesValid(t *testing.T) {
	cmd := NewCronvCommand()
	cmd.Duration = "1m"
	ret, _ := cmd.toDurationMinutes()
	assert.Equal(t, ret, float64(1), "")

	cmd = NewCronvCommand()
	cmd.Duration = "2h"
	ret2, _ := cmd.toDurationMinutes()
	assert.Equal(t, ret2, float64(1*60*2), "")

	cmd = NewCronvCommand()
	cmd.Duration = "3d"
	ret3, _ := cmd.toDurationMinutes()
	assert.Equal(t, ret3, float64(1*60*24*3), "")
}

func TestToDurationMinutesInValid(t *testing.T) {
	cmd := NewCronvCommand()
	cmd.Duration = "INVALID"
	_, err := cmd.toDurationMinutes()
	assert.NotNil(t, err)

	cmd = NewCronvCommand()
	cmd.Duration = "1F"
	_, err2 := cmd.toDurationMinutes()
	assert.NotNil(t, err2)
}

func TestToFromTime(t *testing.T) {
	cmd := NewCronvCommand()
	cmd.FromDate = "2016/11/08"
	cmd.FromTime = "01:30"
	ret, err := cmd.toFromTime()
	assert.NotNil(t, ret)
	assert.Nil(t, err)
}
