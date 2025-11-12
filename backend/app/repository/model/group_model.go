package model

import (
	"database/sql"
	"fmt"
	"time"

	"social-network/app/utils"
	"social-network/pkg/db/sqlite"
)

func FetchUserGroups(userID string) ([]utils.Group, error) {
	query := `SELECT 
		g.id,
		g.title,
		g.description,
		(
			SELECT COUNT(*)
			FROM group_members gm2
			WHERE gm2.group_id = g.id
		) AS member_count
	FROM groups g
	WHERE g.id IN (
		SELECT gm.group_id
		FROM group_members gm
		WHERE gm.user_id = ?
	);`

	rows, err := sqlite.Db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []utils.Group
	for rows.Next() {
		var g utils.Group
		if err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.MemberCount); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func FetchAllAvailableGroups(userID string) ([]utils.Group, error) {
	query := `SELECT 
		g.id, 
		g.title, 
		g.description,
		(
			SELECT COUNT(*)
			FROM group_members gm2
			WHERE gm2.group_id = g.id
		) AS member_count
	FROM groups g
	WHERE g.id NOT IN (
		SELECT gm.group_id 
		FROM group_members gm 
		WHERE gm.user_id = ?
	);`

	rows, err := sqlite.Db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []utils.Group
	for rows.Next() {
		var g utils.Group
		if err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.MemberCount); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func InsertNewGroup(tx *sql.Tx, groupID, title, description, adminID string) error {
	query := `INSERT INTO groups (id, title, description, admin_id) VALUES (?, ?, ?, ?)`
	_, err := tx.Exec(query, groupID, title, description, adminID)
	return err
}

func InsertAdminAsMember(tx *sql.Tx, adminID, groupID string) error {
	query := `INSERT INTO group_members (user_id, group_id) VALUES (?, ?)`
	_, err := tx.Exec(query, adminID, groupID)
	return err
}

func InsertGroupInvitation(tx *sql.Tx, rowID, groupID, userID, adminID string, createdAt time.Time) error {
	query := `INSERT INTO group_invitations (id, group_id, user_id, invited_by_user_id, request_type, created_at)
			VALUES (?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(query, rowID, groupID, userID, adminID, "invitation", createdAt)
	return err
}

func FetchCreatedGroup(groupID string) (utils.Group, error) {
	var g utils.Group
	query := `
		SELECT g.id, g.title, g.description 
		FROM groups g 
		ORDER BY g.created_at DESC
		LIMIT 1`
	err := sqlite.Db.QueryRow(query, groupID).Scan(&g.ID, &g.Title, &g.Description)
	return g, err
}

func CheckExistingInvitation(tx *sql.Tx, userID, groupID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
				SELECT 1 
				FROM group_invitations
				WHERE user_id = ? AND group_id = ?
			)`
	err := tx.QueryRow(query, userID, groupID).Scan(&exists)
	return exists, err
}

func CheckGroupMembership(tx *sql.Tx, userID, groupID string) (bool, error) {
	fmt.Println("USeeer is howa L", userID)
	var isMember bool
	query := `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	err := tx.QueryRow(query, userID, groupID).Scan(&isMember)
	return isMember, err
}

func InsertJoinRequest(tx *sql.Tx, invitationID, groupID, userID string) error {
	query := `INSERT INTO group_invitations (id, group_id, user_id, invited_by_user_id, request_type, created_at)
	        VALUES (?, ?, ?, ?, ?, ?)`
	createdAt := time.Now().UTC()
	_, err := tx.Exec(query, invitationID, groupID, userID, nil, "join", createdAt)
	return err
}

func CheckUserMembershipOrInvitation(tx *sql.Tx, userID, groupID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
				SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?
				UNION ALL
				SELECT 1 FROM group_invitations WHERE user_id = ? AND group_id = ?
			)`
	err := tx.QueryRow(query, userID, groupID, userID, groupID).Scan(&exists)
	return exists, err
}

func CheckUserExists(tx *sql.Tx, userID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = ?)`
	err := tx.QueryRow(query, userID).Scan(&exists)
	return exists, err
}

func InsertInvitation(tx *sql.Tx, invitationID, groupID, invitedUser, invitedBy string) error {
	query := `INSERT INTO group_invitations (id, group_id, user_id, invited_by_user_id, request_type, created_at)
	        VALUES (?, ?, ?, ?, ?, ?)`
	createdAt := time.Now().UTC()
	_, err := tx.Exec(query, invitationID, groupID, invitedUser, invitedBy, "invitation", createdAt)
	return err
}

func GetUserIDByInvitation(invitationID string) (string, error) {
	var userID string
	err := sqlite.Db.QueryRow(`SELECT user_id FROM group_invitations WHERE id = ?`, invitationID).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func GetGroupIDByInvitation(invitationID, userID string) (string, error) {
	var groupID string
	err := sqlite.Db.QueryRow(`
        SELECT group_id FROM group_invitations 
        WHERE id = ? AND user_id = ?`, invitationID, userID).Scan(&groupID)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("invalid invitation")
	}
	if err != nil {
		return "", err
	}
	return groupID, nil
}

func AddUserToGroup(tx *sql.Tx, userID, groupID string) error {
	query := `INSERT INTO group_members (user_id, group_id) VALUES (?, ?)`
	_, err := tx.Exec(query, userID, groupID)
	return err
}

func DeleteInvitation(tx *sql.Tx, invitationID string) error {
	query := `DELETE FROM group_invitations WHERE id = ?`
	_, err := tx.Exec(query, invitationID)
	return err
}

func FetchJoinRequests(userID string, groupID string) ([]utils.JoinRequest, error) {
	query := `
	SELECT 
		gi.id AS invitation_id,
		gi.user_id,
		u.first_name,
		u.last_name
	FROM group_invitations gi
	JOIN groups g ON gi.group_id = g.id
	JOIN users u ON gi.user_id = u.id
	WHERE g.admin_id = ? 
	AND gi.request_type = 'join'
	AND g.id = ?;
	`

	rows, err := sqlite.Db.Query(query, userID, groupID)
	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var joinRequests []utils.JoinRequest
	for rows.Next() {
		var req utils.JoinRequest
		if err := rows.Scan(&req.InvitationID, &req.UserID, &req.FirstName, &req.LastName); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		joinRequests = append(joinRequests, req)
	}

	return joinRequests, nil
}

func FetchGroupInvitations(userID string) ([]utils.FetchGroupInvitation, error) {
	query := `
	SELECT 
    g.id AS group_id,
    g.title,
    i.id AS invitation_id
	FROM groups g
	JOIN group_invitations i ON i.group_id = g.id
	WHERE i.user_id = ? AND i.request_type = 'invitation'
	`

	rows, err := sqlite.Db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var invitations []utils.FetchGroupInvitation
	for rows.Next() {
		var inv utils.FetchGroupInvitation
		if err := rows.Scan(&inv.GroupID, &inv.Title, &inv.InvitationID); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		invitations = append(invitations, inv)
	}
	fmt.Println("invitations in model :", invitations)

	return invitations, nil
}

func FetchFriendsNotInGroup(userID, groupID string) ([]utils.Friend, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.image
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		WHERE f.user_id = ?
		AND u.id NOT IN (
			SELECT user_id FROM group_members WHERE group_id = ?
		)
	`

	rows, err := sqlite.Db.Query(query, userID, groupID)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}
	defer rows.Close()

	var friends []utils.Friend
	for rows.Next() {
		var f utils.Friend
		if err := rows.Scan(&f.ID, &f.FirstName, &f.LastName, &f.Image); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		friends = append(friends, f)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return friends, nil
}
