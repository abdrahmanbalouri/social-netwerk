package handlers

import (
	"net/http"
	"strings"

	service "social-network/app/api/service"
	"social-network/app/helper"
)

func Getlastcommnet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	_, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Comment not found")
		return
	}
	commentID := parts[3]

	comment, err := service.GetComment(commentID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			helper.RespondWithError(w, http.StatusNotFound, "Comment not found")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, comment)
}
