package handler

import (
	"net/http"
)

// ProjectHandler handles GET /projects.
// Expected payload: {"project": "web1", "testsuite_id": "login", "email": ""}
func (h *Handler) ProjectHandler(w http.ResponseWriter, r *http.Request) {
	runResp, err := h.projectUsecase.ShowProject()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	respondJSON(w, http.StatusOK, StandardResponse{
		Status:  "success",
		Message: "Projects data retrieved",
		Data:    runResp,
	})
}
