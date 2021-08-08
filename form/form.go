package form

type HumidityForm struct {
	Device string  `json:"device"`
	Value  float32 `json:"value"`
}

type TemperatureForm struct {
	Device string  `json:"device"`
	Value  float32 `json:"value"`
}
