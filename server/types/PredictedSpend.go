package types
import "time"

type PredictedSpend struct {
	Category    string    `json:"category"`
	Likelihood  float64   `json:"likelihood"`
	PredictedDate time.Time `json:"predictedDate"`
	Warning     string    `json:"warning"`
}
