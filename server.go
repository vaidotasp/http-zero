/*
HTTP Server implementation without using http std lib
this is for purely learning purposes, hoping to touch on tcp/ip, maybe udp and some networking basics and use Go for this thing

Because HTTP is built upon TCP, we should start with that. We need some info on what TCP protocol is. Essentially it is a communications standard (https://en.wikipedia.org/wiki/Transmission_Control_Protocol).

TCP is connection-oriented, and a connection between client and server is established before data can be sent. The server must be listening (passive open) for connection requests from clients before a connection is established.

TCP is a reliable byte stream delivery service which guarantees that all bytes received will be identical and in the same order as those sent. Since packet transfer by many networks is not reliable, TCP achieves this using a technique known as positive acknowledgement with re-transmission.

While IP handles actual delivery of the data, TCP keeps track of segments - the individual units of data transmission that a message is divided into for efficient routing through the network. For example, when an HTML file is sent from a web server, the TCP software layer of that server divides the file into segments and forwards them individually to the internet layer in the network stack.


HTTP is an application layer protocol designed within the framework of the Internet protocol suite. Its definition presumes an underlying and reliable transport layer protocol.[19] In the latest version HTTP/3, the Transmission Control Protocol (TCP) is no longer used, but the older versions are still more used and they most commonly use TCP.


a request line, consisting of the case-sensitive request method, a space, the requested URL, another space, the protocol version, a carriage return, and a line feed
GET /images/logo.png HTTP/1.1

request header fields:
Host: www.example.com
Accept-Language: en


Response syntax
a status line, consisting of the protocol version, a space, the response status code, another space, a possibly empty reason phrase, a carriage return and a line feed, e.g.:
HTTP/1.1 200 OK
Content-Type: text/html


*/

/*
https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview

Client Request:
GET / HTTP/1.1
Host: www.example.com
User-Agent: Mozilla/5.0
Accept: text/html,application/xhtml+xml
Accept-Language: en-GB,en;q=0.5
Accept-Encoding: gzip, deflate, br
Connection: keep-alive

Server response
HTTP/1.1 200 OK
Date: Mon, 23 May 2005 22:38:34 GMT
Content-Type: text/html; charset=UTF-8
Content-Length: 155
Last-Modified: Wed, 08 Jan 2003 23:11:55 GMT
Server: Apache/1.3.3.7 (Unix) (Red-Hat/Linux)
ETag: "3f80f-1b6-3e1cb03b"
Accept-Ranges: bytes
Connection: close

<html>
  <head>
    <title>An Example Page</title>
  </head>
  <body>
    <p>Hello World, this is a very simple HTML document.</p>
  </body>
</html>
*/

package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	fmt.Println("connection handling")
	clientAddr := conn.RemoteAddr()
	fmt.Printf("Client address: %s\n", clientAddr)

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("buffer: %s", buffer)

	// example of response here
	r := []byte("HTTP/1.1 200 OK\r\nConnection: close\r\nContent-Type: text/html\r\nContent-Length: 19\r\n\r\n<h1>Hola Mundo</h1>")
	conn.Write(r)

	defer conn.Close()
}

func main() {
	fmt.Println("TCP server started")
	// get TCP server going
	l, err := net.Listen("tcp", "localhost:1337")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	// this basically is while (true)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}

}
