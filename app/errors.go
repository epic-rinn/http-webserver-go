package main

func (app *application) NotFound(w ResponseWriter, r *Request) {
	w.Write(StatusNotFound, []byte("Not Found"))
}
