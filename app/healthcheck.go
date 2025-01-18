package main

func (app *application) Healthcheck(w ResponseWriter, r *Request) {
	w.Write(StatusOK, []byte("OK!"))
}
