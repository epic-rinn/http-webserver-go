package main

func HandleRoute(rw ResponseWriter, req *Request) {
	switch req.Path {
	case "/":
		rw.Write(StatusOK, []byte("OK"))
	case "/echo/:str":
		rw.Write(StatusOK, []byte(req.Path[len("/echo/"):]))
	default:
		rw.Write(StatusNotFound, []byte("Not Found"))
	}
}
