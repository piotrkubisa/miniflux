// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package api // import "miniflux.app/api"

import (
	"errors"
	"net/http"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
)

// CreateCategory is the API handler to create a new category.
func (c *Controller) CreateCategory(w http.ResponseWriter, r *http.Request) {
	category, err := decodeCategoryPayload(r.Body)
	if err != nil {
		json.BadRequest(w, r, err)
		return
	}

	userID := request.UserID(r)
	category.UserID = userID
	if err := category.ValidateCategoryCreation(); err != nil {
		json.BadRequest(w, r, err)
		return
	}

	if c, err := c.store.CategoryByTitle(userID, category.Title); err != nil || c != nil {
		json.BadRequest(w, r, errors.New("This category already exists"))
		return
	}

	if err := c.store.CreateCategory(category); err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, category)
}

// UpdateCategory is the API handler to update a category.
func (c *Controller) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := request.RouteInt64Param(r, "categoryID")

	category, err := decodeCategoryPayload(r.Body)
	if err != nil {
		json.BadRequest(w, r, err)
		return
	}

	category.UserID = request.UserID(r)
	category.ID = categoryID
	if err := category.ValidateCategoryModification(); err != nil {
		json.BadRequest(w, r, err)
		return
	}

	err = c.store.UpdateCategory(category)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, category)
}

// GetCategories is the API handler to get a list of categories for a given user.
func (c *Controller) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.store.Categories(request.UserID(r))
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.OK(w, r, categories)
}

// RemoveCategory is the API handler to remove a category.
func (c *Controller) RemoveCategory(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)
	categoryID := request.RouteInt64Param(r, "categoryID")

	if !c.store.CategoryExists(userID, categoryID) {
		json.NotFound(w, r)
		return
	}

	if err := c.store.RemoveCategory(userID, categoryID); err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.NoContent(w, r)
}
