package service

import (
	"Week02/dao"
	"fmt"
	"net/http"
	"strconv"
)

// Authentication is func
func Authentication(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	v := r.URL.Query()
	id, err := strconv.Atoi(v.Get("UserID"))
	// id, err := strconv.Atoi(r.PostForm.Get("userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)

		return
	}
	if err := getUser(id); err != nil {
		fmt.Printf("%+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(""))
}

func getUser(ID int) error {
	_, err := dao.QueryUserByID(ID)
	return err
}
