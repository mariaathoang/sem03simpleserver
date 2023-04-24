package main

import (
	"io"
	"log"
	"net"
	"sync"
	"github.com/mariaathoang/is105sem03/mycrypt"
	"fmt"
	"strconv"
	"strings"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:5000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekryptert melding: ", string(dekryptertMelding))
					kryptertMelding := mycrypt.Krypter([]rune(string(dekryptertMelding)), mycrypt.ALF_SEM03, 4)
					log.Println("Kryptert melding: ", string(kryptertMelding))
					/*switch msg := string(dekryptertMelding); msg {
  				        case "ping":
						_, err = c.Write([]byte("pong"))
					case strings.HasPrefix(msg, "Kjevik"):
						celsius, err := strconv.ParseFloat(msg[6:], 64)
						if err != nil {
							log.Println(err)
							return
						}
						fahrenheit := (celsius * 1.8) + 32
						response := fmt.Sprintf("Temperature in Fahrenheit is: %.2f", fahrenheit)
						_, err = c.Write([]byte(response))
					default:
                                                _, err = c.Write([]byte(string(kryptertMelding)))
					}*/
				msg := string(dekryptertMelding)
				if msg == "ping" {
                                                _, err = c.Write([]byte("pong"))
                                        } else if strings.HasPrefix(msg, "Kjevik") {
                                                celsius, err := strconv.ParseFloat(msg[len(msg)-1:], 64)
                                                if err != nil {
                                                        log.Println(err)
							fmt.Println("Test")
                                                        return
                                                }
                                                fahrenheit := (celsius * 1.8) + 32
                                                response := fmt.Sprintf("Temperature in Fahrenheit is: %.2f", fahrenheit)
                                                _, err = c.Write([]byte(response))
                                        } else {
                                                _, err = c.Write([]byte(string(kryptertMelding)))
                                        }
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
