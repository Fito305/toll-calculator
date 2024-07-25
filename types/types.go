package types


type OBUData struct {
	OBUID int `json:"obuID"`
	Lat float64 `json:"lat"` // Latitude
	Long float64 `json:"long"` // Longitude
}
