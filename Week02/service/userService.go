package service

import (
	"Week02/dao"
	"net/http"
	"strconv"
	"fmt"
)

func Authentication(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	id, err := strconv.Atoi(r.Form.Get("userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if err := getUser(id); err != nil {
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
