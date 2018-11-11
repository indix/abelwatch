package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewWatch(t *testing.T) {
	rawJSON := []byte(`{
    "duration": 0,
    "name": "completed",
    "tags": [
        "www.finishline.com",
        "rmn_p1_variant_20180918"
    ],
    "condition": {
        "op": "<",
        "value": 10250
    },
    "slackChannel": "#z_development"
}`)
	watch := NewWatch("1", rawJSON, nil)

	assert.Equal(t, "1", watch.ID)
	assert.Equal(t, int64(0), watch.Duration)
	assert.Equal(t, "<", watch.Condition.Op)
	assert.Equal(t, int64(10250), watch.Condition.Value)
	assert.Equal(t, "completed [[www.finishline.com rmn_p1_variant_20180918]] at 0s", watch.String())
	assert.Equal(t, int64(-1), watch.NextCheck)
	assert.Equal(t, true, watch.LastChecked > 1000)
}

func TestNewCondition(t *testing.T) {
	condition := NewCondition(">", int64(10240))
	assert.Equal(t, ">", condition.Op)
	assert.Equal(t, int64(10240), condition.Value)
}

func TestConditionHasBreached(t *testing.T) {
	assert.Equal(t, true, NewCondition(">", int64(10240)).HasBreached(int64(10241)))
	assert.Equal(t, false, NewCondition(">", int64(10240)).HasBreached(int64(10240)))

	assert.Equal(t, true, NewCondition("<", int64(10240)).HasBreached(int64(10239)))
	assert.Equal(t, false, NewCondition("<", int64(10240)).HasBreached(int64(10241)))

	assert.Equal(t, true, NewCondition("=", int64(10240)).HasBreached(int64(10240)))
	assert.Equal(t, false, NewCondition("=", int64(10240)).HasBreached(int64(10241)))

	assert.Equal(t, true, NewCondition("<=", int64(10240)).HasBreached(int64(10240)))
	assert.Equal(t, true, NewCondition("<=", int64(10240)).HasBreached(int64(10240)))
	assert.Equal(t, false, NewCondition("<=", int64(10240)).HasBreached(int64(10241)))

	assert.Equal(t, true, NewCondition(">=", int64(10240)).HasBreached(int64(10240)))
	assert.Equal(t, true, NewCondition(">=", int64(10240)).HasBreached(int64(10241)))
	assert.Equal(t, false, NewCondition(">=", int64(10240)).HasBreached(int64(10239)))

	assert.Equal(t, false, NewCondition("!", int64(10240)).HasBreached(int64(10241)))
}
