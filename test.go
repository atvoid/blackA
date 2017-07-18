package main
/*
import (
	"fmt"
	"blackA/room"
	"net"
)
*/
func main() {
	/*
	room := room.CreateRoom()
	listener, _ := net.Listen("tcp", "192.168.199.189:789")
	id := 0
	for {
		conn, _ := listener.Accept()
		id++
		if room.AddUser(id, string(id), &conn) {
			fmt.Println("someone In:", id)
			//conn.Write([]byte("Success"))
			//c <- id
		}
		if room.IsFull() {
			fmt.Println("Started");
			go room.Start()
		}
	}
	*/
}
