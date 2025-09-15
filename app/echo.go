package main

func (app *application) Echo(w ResponseWriter, r *Request) {
	w.Write(StatusOK, []byte(r.Params["str"]))
}
