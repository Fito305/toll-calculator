package main

import (
	"fmt"
    "log"
	"time"
	"math/rand"

	// "github.com/alexsasharegan/dotenv"
    "github.com/Fito305/tolling/types"
    "github.com/gorilla/websocket"
)
const wsEndpoint = "ws://127.0.0.1:30000/ws"
var sendInterval = time.Second * 5


func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func main() {
	obuIDS := generateOBUIDS(20)
    conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint,  nil)
    if err != nil {
        log.Fatal(err)
    }
	for {
		// Each time we send a coordinate we are going to pick a random OBU.
		for i := 0; i <len(obuIDS); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat: lat,
				Long: long,
			}
			fmt.Printf("%v+\n", data)
            if err := conn.WriteJSON(data); err != nil {
                log.Fatal(err)
            }
		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(999999)
	}
	return ids
}

func init() {
	// add this line to use the `github.com/alexsasharegen/dotenv` .env file Enviroment variables.
	// if err := dotenv.Load(); err != nil {
	// 	log.Fatal(err)
	// }
	rand.Seed(time.Now().UnixNano()) // Why do we need to do that? To make sure we have random data. 
}
