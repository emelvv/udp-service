package main

import (
    "log"
	"os"
	"net"
	"io/ioutil"
	"fmt"
	"time"
	"crypto/sha1"
)

func fileName(path string) (name string) {
	for i:=len(path)-1;i>=0;i-=1{
		if string(path[i]) == "/"{return}
		name=string(path[i])+name
	}
	return
}

func divide(buffer []byte, n int, uuid []byte) (res [][]byte){
	c:=0
	piece_counter :=1
	var a []byte
	s := false
	for _, num := range buffer{
		s = false
		if c >=n{
			s = true
			send := []byte{1}
			send = append(send, uuid...)   // тип (1 байт) + uuid (36 байт) + номер куска(1 байт)
			send = append(send, byte(piece_counter))
			piece_counter+=1
			
			send = append(send, a...)
			res = append(res, send)
			a =[]byte{}
			c=0
		}
		c+=1
		a = append(a, num)
	}
	if !s{
		send := []byte{1}
		send = append(send, uuid...)   // тип (1 байт) + uuid (36 байт) + номер куска(1 байт)
		send = append(send, byte(piece_counter))
		send = append(send, a...)
		res = append(res, send)
	}
	return
}

func comb(a []string) (res string){
	
	for i, num := range a {
		if i != 0{
			res=res + " " +num
		}else{
			res = num
		}
	}
	return
}



//go run transmitter.go 192.168.1.67:1221 file.txt  

func main() {

	if len(os.Args)==1{log.Fatal("Adress needed")}
	if len(os.Args)==2{log.Fatal("File path needed")}

	buffer := make([]byte, 65000)

	file, err := ioutil.ReadFile(comb(os.Args[2:]))
	if err!=nil{log.Fatal(err)}
	
	check_sum := sha1.Sum(file) // контрольная сумма

	remoteAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	conn, err := net.ListenPacket("udp", "")


	_, err=conn.WriteTo([]byte{2}, remoteAddr)
	
	var uuid []byte
	n, _, _ := conn.ReadFrom(buffer)
	uuid = buffer[:n]
	buffer = make([]byte, 65000)
	log.Printf("Your uuid is %v", string(uuid))
	defer conn.Close() // выполнится в конце
	
	var all_bytes int

	all:=divide(file, 60000, uuid)

	title := []byte{byte(0)}
	title = append(title, uuid...)
	title = append(title, []byte(fileName(comb(os.Args[2:])))...)
	title = append(title, byte(0))
	title = append(title, []byte(fmt.Sprintf("%v", len(all)))...)
	title = append(title, byte(0))
	title = append(title, check_sum[:]...)

	b, _ := conn.WriteTo(title, remoteAddr)

	all_bytes+=b

	log.Print(fmt.Sprintf("Sending %v, in %v piece(s)", fileName(comb(os.Args[2:])), len(all)))

	// all[0], all[1], all[2] = all[2], all[0], all[1] // отправляю в разном порядке

	for i, num := range all{
		b, err:=conn.WriteTo(num, remoteAddr)
		if err!=nil{
			log.Printf("ERROR: %v", err)
		}else{
			all_bytes+=b
			log.Printf("Sending file... (%v/%v)", i+1, len(all))
		}
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("Written %v bytes to %v", all_bytes, remoteAddr)


	_, _, err = conn.ReadFrom(buffer)

	log.Print(string(buffer))

}

