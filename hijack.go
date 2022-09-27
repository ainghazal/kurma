package kurma

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

// HijackAndProxyHandler returns a http handler function that will hijack the
// underlying connection and stablish a transparent TCP proxy to the passed
// remoteAddr address.
func HijackAndProxyHandler(remoteAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			panic(err)
		}
		conn.Write([]byte{})
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n")

		go func() {
			defer conn.Close()
			conn2, err := net.Dial("tcp", remoteAddr)
			log.Println("Dialing remote:", remoteAddr)

			if err != nil {
				log.Println("error dialing remote addr", err)
				return
			}
			defer conn2.Close()
			closer := make(chan struct{}, 2)
			go copyBytes(closer, conn2, conn)
			go copyBytes(closer, conn, conn2)
			<-closer
			log.Println("Connection complete", conn.RemoteAddr())
		}()
	}
}

func copyBytes(closer chan struct{}, dst io.Writer, src io.Reader) {
	_, _ = io.Copy(dst, src)
	closer <- struct{}{} // connection is closed, send signal to stop proxy
}
