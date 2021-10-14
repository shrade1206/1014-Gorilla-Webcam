package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	fmt.Println("Go WebSocket")

	setRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func reader(conn *websocket.Conn) {
	reply := []byte("已收到訊息")
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("使用者訊息: " + string(p))

		if err := conn.WriteMessage(messageType, reply); err != nil {
			log.Println(err)
			return
		}

		// 拍照存擋，前端輸入訊息 = 檔名
		webcam, err := gocv.VideoCaptureDevice(0)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second)

		img := gocv.NewMat()

		webcam.Read(&img)

		gocv.IMWrite(string(p)+".jpg", img)

		webcam.Close()
		img.Close()

		break

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("使用者已連線")

	reader(ws)
}

// func client(conn *websocket.Conn) {
// 	err := conn.WriteMessage(websocket.TextMessage, []byte("Welcome Tese"))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }

func setRoutes() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/ws", wsEndpoint)

}
