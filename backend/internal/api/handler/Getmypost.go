package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	service "social-network/internal/api/service"
	"social-network/internal/helper"
	"social-network/internal/repository"
)

func Getmypost(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	authUserID := parts[3]



	offset, _ := strconv.Atoi(parts[4])
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if authUserID == "0" {
		authUserID = userID
	}
		
	

	posts, err := service.FetchPostsByUser(repository.Db , userID, authUserID, offset,10)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(posts)
	helper.RespondWithJSON(w, http.StatusOK, posts)
}
