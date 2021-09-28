package album

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	userFilters "github.com/tupyy/gophoto/internal/domain/filters/user"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

// parseAlbumIDHandler decrypt the album id passes as parameters and set the id in the context.
func parseAlbumIDHandler(c *gin.Context) {
	logger := logutil.GetLogger(c)

	// decrypt album id
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	decryptedID, err := gen.DecryptData(c.Param("id"))
	if err != nil {
		logger.WithError(err).Error("cannot decrypt album id")
		c.AbortWithError(http.StatusNotFound, err) // explicit return not found here

		return
	}

	id, err := strconv.Atoi(decryptedID)
	if err != nil {
		logger.WithError(err).WithField("id", decryptedID).Error("cannot parse album id")
		c.AbortWithError(http.StatusNotFound, err)

		return
	}

	c.Set("id", id)
}

// generateFilters generates 3 filters: notUserNameFilter, FilterByRole and FilterByCanShare.
func generateFilters(currentUser entity.User) (userFilters.Filters, error) {
	filters := make(userFilters.Filters)

	// get other users with can_share true except the current user
	usernameFilter, err := userFilters.GenerateFilterFuncs(userFilters.NotFilterByUsername, []string{currentUser.Username})
	if err != nil {
		return nil, err
	}

	filters[currentUser.ID] = usernameFilter

	// only can share users
	canShareFilter, err := userFilters.GenerateFilterFuncs(userFilters.FilterByCanShare, []string{})
	if err != nil {
		return nil, err
	}

	filters["canshare"] = canShareFilter

	// remove admins
	notAdminFilter, err := userFilters.GenerateFilterFuncs(userFilters.FilterByRole, []entity.Role{entity.RoleUser, entity.RoleEditor})
	if err != nil {
		return nil, err
	}

	filters["admin"] = notAdminFilter
	return filters, nil
}
