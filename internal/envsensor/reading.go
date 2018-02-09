package envsensor

type Reading struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}
