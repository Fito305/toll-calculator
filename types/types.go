package types

// This should be transport indepedent. Which means that
// json should not play a role in this thing.
type Invoice struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
}

// Types is BUSINESS DOMAIN. This Distance is what your business logic needs
// to operate.
type Distance struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"` // This is a timestamp - unix is used instead of timestamp.
}

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`  // Latitude
	Long  float64 `json:"long"` // Longitude
	RequestID int `json:"requestID"` 
}

// OBUData.RequestID allows you to `trace` the requests from the clients 
// as that request travels through all the microservices. In the case that there
// is a bug, the RequestID will help you track that bug and find out where it occured
// in its flight, and in which microservice it might have occured in. The request starts
// at the client (maybe a mobile phone), then it goes to the GATEWAY, and finally it travels 
// through the microservices. The response is sent back to the client if no errors occur.
// The RequestID allows you to leave a trail of bread crums to find bugs that occur while
// the request from the client is in flight.
// You use `Elastic Search` to find the RequestID, that is basically `tracing`.

// Another option instead of the type OBUData.RequestID, you can use HTTP / GRPC headers. 
// These headers have request ids on them. 

