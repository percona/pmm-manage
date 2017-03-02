package user

import (
	"database/sql"
	"github.com/grafana/grafana/pkg/util"
	_ "github.com/mattn/go-sqlite3" // sqlite driver requires such import
)

func createGrafanaUser(newUser PMMUser) error {
	email := newUser.Username + "@localhost"
	salt := util.GetRandomString(10)
	rands := util.GetRandomString(10)
	password := util.EncodePassword(newUser.Password, salt)

	db, err := sql.Open("sqlite3", PMMConfig.GrafanaDBPath)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	if _, err = db.Exec("PRAGMA busy_timeout = 60000"); err != nil {
		return err
	}

	affect, err := updateUser(db, newUser.Username, email, password, salt, rands)
	if err != nil {
		return err
	}

	if affect == 0 {
		userID, err := insertUser(db, newUser.Username, email, password, salt, rands)
		if err != nil {
			return err
		}
		return addUserToOrg(db, userID)
	}

	return nil
}

func updateUser(db *sql.DB, username, email, password, salt, rands string) (int64, error) {
	stmt, err := db.Prepare(`
        UPDATE user
        SET    version  = 1,
               email    = ?,
               password = ?,
               salt     = ?,
               rands    = ?,
               updated  = date('now')
        WHERE  login    = ?
    `)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(email, password, salt, rands, username)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func insertUser(db *sql.DB, username, email, password, salt, rands string) (int64, error) {
	stmt, err := db.Prepare(`
        INSERT INTO user (version, login, email, password, salt, rands, org_id, is_admin,     created,     updated)
        VALUES           (      1,     ?,     ?,        ?,    ?,     ?,      1,        1, date('now'), date('now'))
    `)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(username, email, password, salt, rands)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func addUserToOrg(db *sql.DB, userID int64) error {
	stmt, err := db.Prepare(`
        INSERT INTO org_user (org_id, user_id,    role,     created,     updated)
        VALUES               (     1,       ?, 'Admin', date('now'), date('now'));
    `)
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(userID); err != nil {
		return err
	}
	return nil
}

func deleteGrafanaUser(username string) error {
	db, err := sql.Open("sqlite3", PMMConfig.GrafanaDBPath)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	stmt, err := db.Prepare("DELETE FROM user WHERE login = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username)
	return err
}
