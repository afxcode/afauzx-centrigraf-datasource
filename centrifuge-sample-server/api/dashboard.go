package api

import (
	"encoding/json"
)

var vref = 0

func (c *CentrifugoAPI) BroadcastKV() (err error) {
	type KV struct {
		Key   string  `json:"keyx"`
		Value float64 `json:"valuex"`
	}

	payloadFanSpeed, err := json.Marshal(KV{
		Key:   "fan_speed",
		Value: float64(vref) * 11.1,
	})
	if err != nil {
		return err
	}

	payloadCpuTemp, err := json.Marshal(KV{
		Key:   "cpu_temp",
		Value: float64(vref) * 1.2,
	})
	if err != nil {
		return err
	}

	vref++
	if vref > 100 {
		vref = 1
	}

	if _, err = c.node.Publish("fan_speed", payloadFanSpeed); err != nil {
		return
	}
	if _, err = c.node.Publish("cpu_temp", payloadCpuTemp); err != nil {
		return
	}
	return
}
