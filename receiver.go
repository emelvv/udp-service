package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func inMap(element string, Map map[string]chan []byte) bool {
	for key, _ := range Map {
		if key == element {
			return true
		}
	}
	return false
}

func correctFileName(fileName string) string {

	fileFormat := func(fileName string) (string, string) {
		for i, num := range fileName {
			if string(num) == "." {
				return fileName[i:], fileName[:i]
			}
		}
		return "", fileName
	}

	exists := func(path string) bool {
		_, err := os.Stat(path)
		if err == nil {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	if !exists(fileName) {
		return fileName
	} else {

		format, fileTag := fileFormat(fileName)

		for i := 1; ; i++ {
			var (
				name string
				c    int
			)

			_, _ = fmt.Sscanf(fileName, "%v (%v)", &name, &c)

			fileName = fmt.Sprintf("%v (%v)%v", fileTag, i, format)

			if !exists(fileName) {
				return fileName
			}
		}
	}
}

func cr_data(data_map map[int][]byte) (res []byte) {
	for i := 1; i < len(data_map)+1; i++ {
		res = append(res, data_map[i]...)
	}
	return
}

func receive(bytes_chan chan []byte, proc *(map[string]chan []byte), remoteAddr net.Addr, uuid string, conn net.PacketConn, timers *(map[string]*time.Timer)) any {
	// var name, l, hash []byte
	// c := 0
	// for _, num := range title{
	//     if num == byte(0){c+=1
	//     }else {
	//         switch c {
	//         case 0: name=append(name, num)
	//         case 1: l = append(l, num)
	//         case 2: hash = append(hash, num)
	//         }
	//     }
	// }

	data_map := make(map[int][]byte)
	var data []byte
	var name, l, hash []byte
	for {
		select {
		case d := <-bytes_chan:
			data = d
		}
		if data == nil {
			return nil
		}

		if data[0] == byte(0) {
			c := 0
			for _, num := range data[37:] {
				if num == byte(0) {
					c += 1
				} else {
					switch c {
					case 0:
						name = append(name, num)
					case 1:
						l = append(l, num)
					case 2:
						hash = append(hash, num)
					}
				}
			}
			break
		}

		var p []byte
		data_map[int(data[37])] = append(p, data...)[38:]
	}

	l_int, _ := strconv.Atoi(string(l))

	for l_int != len(data_map) {
		select {
		case d := <-bytes_chan:
			data = d
		}
		if data == nil {
			return nil
		}

		var j []byte
		data_map[int(data[37])] = append(j, data...)[38:]
	}

	// for key, val:=range data_map{
	//     log.Print(key, val[:100])
	// }

	all_data := cr_data(data_map)

	x := sha1.Sum(all_data)
	if string(hash) == string(x[:]) { //проверка контрольных сумм
		_, err := conn.WriteTo([]byte("File received successfully."), remoteAddr)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(correctFileName(string(name)), all_data, 0644)
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
		log.Printf("%v successfully sended from %v", string(name), uuid)

	} else {
		_, err := conn.WriteTo([]byte("File received unsuccessfully."), remoteAddr)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v: Error in check sum when receiving file %v", uuid, string(name))
	}

	delete(*proc, uuid)
	delete(*timers, uuid)

	return nil
}

// 32
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Port needed.")
	}

	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%s", os.Args[1]))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	buffer := make([]byte, 65000)

	inProcess := make(map[string]chan []byte)
	timers := make(map[string]*time.Timer)

	for {
		bytesRead, remoteAddr, _ := conn.ReadFrom(buffer)

		// log.Print(buffer[:100])

		if buffer[0] == byte(0) || buffer[0] == byte(1) { // если тип 1 или 0

			var uuid string
			uuid = string(buffer[1:37])

			timers[uuid].Stop()

			timers[uuid] = time.AfterFunc(5*time.Second, func() {
				inProcess[uuid] <- nil
				log.Print(uuid + " ERROR: File not received (Timeout expired).")
				timers[uuid].Stop()
				delete(timers, uuid)
			})

			var b []byte
			b = append(b, buffer[:bytesRead]...)

			inProcess[uuid] <- b

		} else if buffer[0] == byte(2) { // если тип 2
			n := uuid.New().String()
			_, err := conn.WriteTo([]byte(n), remoteAddr)
			if err != nil {
				log.Fatal(err)
			}
			inProcess[n] = make(chan []byte, 1)

			timers[n] = time.AfterFunc(5*time.Second, func() {
				inProcess[n] <- nil
				log.Print(n + " ERROR: File not received (Timeout expired).")
				timers[n].Stop()
				delete(timers, n)
			})
			defer timers[n].Stop()

			go receive(inProcess[n], &inProcess, remoteAddr, n, conn, &timers)

		}
	}

}
