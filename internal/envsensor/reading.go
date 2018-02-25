package envsensor

type Reading struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	SensorType  string
}
