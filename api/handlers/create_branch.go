package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/treeverse/lakefs/api/gen/models"
	"github.com/treeverse/lakefs/api/gen/restapi/operations"
	"github.com/treeverse/lakefs/permissions"
)

func NewCreateBranchHandler(serverContext ServerContext) operations.CreateBranchHandler {
	return &createBranchHandler{serverContext}
}

type createBranchHandler struct {
	serverContext ServerContext
}

func (h *createBranchHandler) Handle(params operations.CreateBranchParams, user *models.User) middleware.Responder {
	err := authorize(h.serverContext, user, permissions.ManageRepos, repoArn(params.RepositoryID))
	if err != nil {
		return operations.NewCreateBranchUnauthorized().WithPayload(responseErrorFrom(err))
	}

	err = h.serverContext.GetIndex().CreateBranch(params.RepositoryID, params.BranchID, *params.Branch.CommitID)
	if err != nil {
		return operations.NewCreateBranchDefault(http.StatusInternalServerError).WithPayload(responseErrorFrom(err))
	}

	return operations.NewCreateBranchCreated().WithPayload(&models.Refspec{
		CommitID: params.Branch.CommitID,
		ID:       &params.RepositoryID,
	})
}