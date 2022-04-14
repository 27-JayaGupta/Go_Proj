package main
import(
	"fmt"
	"net/http"
	"io"
	"time"
)
var verbose = false

var passthruRequestHeaderKeys = [...]string{
    "Accept",
    "Accept-Encoding",
    "Accept-Language",
    "Cache-Control",
    "Cookie",
    "Referer",
    "User-Agent",
}

var passthruResponseHeaderKeys = [...]string{
    "Content-Encoding",
    "Content-Language",
    "Content-Type",
    "Cache-Control",
    "Etag",
    "Expires",
    "Last-Modified",
    "Location",
    "Server",
    "Vary",
}

func handleFunc(w http.ResponseWriter, r *http.Request){
	fmt.Printf("--> %v %v\n",r.Method,r.URL)

	hh := http.Header{}

	for _,hk := range passthruRequestHeaderKeys{
		if hv,ok := r.Header[hk]; ok{
			hh[hk] = hv
		}
	}

	rr := &http.Request{
		Method : r.Method,
		URL : r.URL,
		Header: hh,
		Body: r.Body,
		ContentLength: r.ContentLength,
		Close: r.Close,
	}

	resp,err := http.DefaultTransport.RoundTrip(rr)
	if err!=nil {
		http.Error(w,"Could not reach the server",http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()

	if verbose{
		fmt.Printf("--> %v %+v\n",resp.Status,resp.Header)
	} else {
		fmt.Printf("--> %v\n",resp.Status)
	}

	respH :=  w.Header()

	for _,hk := range passthruResponseHeaderKeys{
		if hv,ok := resp.Header[hk];ok{
			respH[hk] = hv;
		}
	}

	w.WriteHeader(resp.StatusCode)

	if resp.ContentLength > 0{
		io.Copy(w, resp.Body)
	} else if (resp.Close){
		for {
			if _,err := io.Copy(w,resp.Body); err!=nil {
				break
			}
		}
	}
}

func main(){

	handler := http.DefaultServeMux
	handler.HandleFunc("/",handleFunc)

	server := &http.Server{
		Addr: ":8888",
		Handler : handler,
		ReadTimeout: 10*time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1<<20,
	}

	server.ListenAndServe()
}