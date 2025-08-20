package main

import "net/http"

func handleSuccess(w http.ResponseWriter,r *http.Request){
	responseWithJson(w,http.StatusOK,struct{}{})
}