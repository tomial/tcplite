package net

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("Failed to listen: ", err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Failed to accept from socket: ", err)
				return
			}

			packer := NewDataPacker()
			// 处理客户端请求
			go func(conn net.Conn) {

				for {
					headerBytes := make([]byte, packer.GetHeadLen())

					_, err := io.ReadFull(conn, headerBytes)
					if err != nil {
						fmt.Println("Failed to read header from conn: ", err)
						return
					}

					header, err := packer.Unpack(headerBytes)
					if err != nil {
						fmt.Println("Failed to unpack header bytes: ", err)
					}

					if header.GetMsgLen() > 0 {

						msg := header.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("Failed to read data from conn: ", err)
						}

						fmt.Println("Successfully received data!")
						fmt.Printf("DataLen: %d, DataID: %d, Data: %s\n", msg.DataLen, msg.Id, msg.Data)
					}
				}
			}(conn)
		}
	}()

	// 客户端
	go func() {

		conn, err := net.Dial("tcp", "127.0.0.1:9000")
		if err != nil {
			fmt.Println("Failed to connect to the server: ", err)
			return
		}

		packer := NewDataPacker()

		msg1 := &Message{
			Id:      1,
			DataLen: 5,
			Data:    []byte("hello"),
		}

		pack1, err := packer.Pack(msg1)
		if err != nil {
			fmt.Println("Failed to pack data: ", err)
			return
		}

		msg2 := &Message{
			Id:      2,
			DataLen: 7,
			Data:    []byte("tcplite"),
		}
		pack2, err := packer.Pack(msg2)
		if err != nil {
			fmt.Println("Failed to pack data: ", err)
			return
		}

		msg := append(pack1, pack2...)

		conn.Write(msg)
	}()

	select {
	case <-time.After(time.Second):
		return
	}
}
