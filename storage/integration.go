// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage

import (
	"database/sql"
	"fmt"

	"github.com/miniflux/miniflux/model"
)

// UserByFeverToken returns a user by using the Fever API token.
func (s *Storage) UserByFeverToken(token string) (*model.User, error) {
	query := `
		SELECT
		users.id, users.is_admin, users.timezone
		FROM users
		LEFT JOIN integrations ON integrations.user_id=users.id
		WHERE integrations.fever_enabled='t' AND integrations.fever_token=$1
	`

	var user model.User
	err := s.db.QueryRow(query, token).Scan(&user.ID, &user.IsAdmin, &user.Timezone)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("unable to fetch user: %v", err)
	}

	return &user, nil
}

// Integration returns user integration settings.
func (s *Storage) Integration(userID int64) (*model.Integration, error) {
	query := `SELECT
			user_id,
			pinboard_enabled,
			pinboard_token,
			pinboard_tags,
			pinboard_mark_as_unread,
			instapaper_enabled,
			instapaper_username,
			instapaper_password,
			fever_enabled,
			fever_username,
			fever_password,
			fever_token
		FROM integrations
		WHERE user_id=$1
	`
	var integration model.Integration
	err := s.db.QueryRow(query, userID).Scan(
		&integration.UserID,
		&integration.PinboardEnabled,
		&integration.PinboardToken,
		&integration.PinboardTags,
		&integration.PinboardMarkAsUnread,
		&integration.InstapaperEnabled,
		&integration.InstapaperUsername,
		&integration.InstapaperPassword,
		&integration.FeverEnabled,
		&integration.FeverUsername,
		&integration.FeverPassword,
		&integration.FeverToken,
	)
	switch {
	case err == sql.ErrNoRows:
		return &integration, nil
	case err != nil:
		return &integration, fmt.Errorf("unable to fetch integration row: %v", err)
	}

	return &integration, nil
}

// UpdateIntegration saves user integration settings.
func (s *Storage) UpdateIntegration(integration *model.Integration) error {
	query := `
		UPDATE integrations SET
			pinboard_enabled=$1,
			pinboard_token=$2,
			pinboard_tags=$3,
			pinboard_mark_as_unread=$4,
			instapaper_enabled=$5,
			instapaper_username=$6,
			instapaper_password=$7,
			fever_enabled=$8,
			fever_username=$9,
			fever_password=$10,
			fever_token=$11
		WHERE user_id=$12
	`
	_, err := s.db.Exec(
		query,
		integration.PinboardEnabled,
		integration.PinboardToken,
		integration.PinboardTags,
		integration.PinboardMarkAsUnread,
		integration.InstapaperEnabled,
		integration.InstapaperUsername,
		integration.InstapaperPassword,
		integration.FeverEnabled,
		integration.FeverUsername,
		integration.FeverPassword,
		integration.FeverToken,
		integration.UserID,
	)

	if err != nil {
		return fmt.Errorf("unable to update integration row: %v", err)
	}

	return nil
}

// CreateIntegration creates initial user integration settings.
func (s *Storage) CreateIntegration(userID int64) error {
	query := `INSERT INTO integrations (user_id) VALUES ($1)`
	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("unable to create integration row: %v", err)
	}

	return nil
}
