package main

import "ChildTuningGoBackend/backend/http"

func main() {
	service := http.NewRouter()
	service.StartServer()
}
