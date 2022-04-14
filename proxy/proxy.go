package main
import(
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"time"
	"fmt"
)

func transfer(dst io.WriteCloser, src io.ReadCloser){
	defer dst.Close()
	defer src.Close()
	io.Copy(dst,src)
}

func handleTunneling(w http.ResponseWriter , r *http.Request){
	fmt.Print("In tunnel")
	dest_conn,err := net.DialTimeout("tcp",r.Host,10*time.Second)
	
	if err!=nil {
		http.Error(w,err.Error(),http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker,ok := w.(http.Hijacker)
	if !ok {
		http.Error(w,"Hijacking not supported",http.StatusInternalServerError)
		return
	}
	client_conn,_,err := hijacker.Hijack()
	if err!=nil {
		http.Error(w,err.Error(),http.StatusServiceUnavailable)
		return
	}
	go transfer(dest_conn,client_conn)
	go transfer(client_conn,dest_conn)
	
}

func copyHeader(dst, src http.Header){
	for key,value := range src{
		for _,v := range value{
			dst.Add(key,v)
		}
	}
}

func handleHTTP( w http.ResponseWriter , r *http.Request){
	resp,err := http.DefaultTransport.RoundTrip(r)
	fmt.Print(resp)
	if err!=nil {
		http.Error(w,err.Error(),http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()
	copyHeader(w.Header(),resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main(){
	var pemPath string
	flag.StringVar(&pemPath,"pem","server.pem","path to pem file")
	var keyPath string
	flag.StringVar(&keyPath,"key","server.key","path to key file")
	var proto string
	flag.StringVar(&proto,"proto","https","Proxy Protocol(http or https)")
	flag.Parse()

	if proto!="http" && proto!="https" {
		log.Fatal("Protocol must be either http or https")
	}

	server := &http.Server{
		Addr : ":8888",
		Handler : http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
			fmt.Print(r.Method)
			if r.Method == http.MethodConnect {
				handleTunneling(w,r)
			} else {
				handleHTTP(w,r)
			}
		}),
		TLSNextProto: make(map[string]func(*http.Server,*tls.Conn,http.Handler)),
	}

	if proto=="http" {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(server.ListenAndServeTLS(pemPath,keyPath))
	}
}