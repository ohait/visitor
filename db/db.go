package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ohait/visitor/ctx"
	"github.com/ohait/visitor/data"
)

type Store interface {
	PersonById(c ctx.C, id string) (data.Person, bool, error)

	UpsertPerson(c ctx.C, v data.Person) error

	EventsByPerson(c ctx.C, id string) ([]data.Event, error)

	InsertEvent(c ctx.C, ev data.Event) error
}

type SQLStore struct {
	DB *sql.DB
}

//var _ Store = &SQLStore{}

func NewSQLStore(c ctx.C, driver, datasource string) (*SQLStore, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, fmt.Errorf("sql.Open(%q,%q): %v", driver, datasource, err)
	}
	s := &SQLStore{
		DB: db,
	}
	return s, s.init(c)
}

func (s *SQLStore) query(c ctx.C, sql string, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.DB.QueryContext(c, sql, args...)
	if err != nil {
		c = ctx.WithTag(c, `sql`, sql)
		return nil, ctx.Wrap(c, err)
	}
	ctx.Log(c).Debugf(`%s -- %+v`, sql, args)
	return rows, nil
}

func (s *SQLStore) exec(c ctx.C, sql string, args ...interface{}) (sql.Result, error) {
	c = ctx.WithTag(c, `sql`, sql)
	res, err := s.DB.ExecContext(c, sql, args...)
	if err != nil {
		return res, ctx.Wrap(c, err)
	}
	return res, nil
}

func (s *SQLStore) init(c ctx.C) error {
	var err error

	_, err = s.query(c, `SELECT count(1) FROM person`)
	if err != nil {
		_, err = s.exec(c, `CREATE TABLE person (
			id VARCHAR(32),
			email VARCHAR(64),
			details TEXT,
			PRIMARY KEY(id))`)
		if err != nil {
			return err
		}
		_, err = s.exec(c, `CREATE INDEX person_email ON person(email)`)
		if err != nil {
			return err
		}
	}

	_, err = s.query(c, `SELECT count(1) FROM event`)
	if err != nil {
		_, err = s.exec(c, `CREATE TABLE event (
			person VARCHAR(32),
			campaign VARCHAR(64),
			time DATETIME,
			type VARCHAR(16))`)
		if err != nil {
			return err
		}
		_, err = s.exec(c, `CREATE INDEX event_person ON event(person)`)
		if err != nil {
			return err
		}
		_, err = s.exec(c, `CREATE INDEX event_campaign ON event(campaign)`)
		if err != nil {
			return err
		}
		_, err = s.exec(c, `CREATE INDEX event_time ON event(time)`)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLStore) EventsByPerson(c ctx.C, id string) ([]data.Event, error) {
	rows, err := s.query(c, `SELECT person, campaign, type, time
		FROM event WHERE person = ? ORDER BY time DESC LIMIT 1000`, id)
	if err != nil {
		return nil, ctx.Wrapf(c, "can't read events: %v", err)
	}
	out := []data.Event{}
	for rows.Next() {
		var ev data.Event
		err = rows.Scan(&ev.Person, &ev.Campaign, &ev.Type, &ev.Time)
		if err != nil {
			return out, ctx.Wrapf(c, "can't scan row: %v", err)
		}
		out = append(out, ev)
	}
	return out, nil
}

func (s *SQLStore) InsertEvent(c ctx.C, ev data.Event) error {
	_, err := s.exec(c, `INSERT INTO event (person, campaign, type, time) VALUES (?,?,?,?)`,
		ev.Person, ev.Campaign, ev.Type, time.Now())
	return ctx.Wrap(c, err)
}

func (s *SQLStore) UpsertPerson(c ctx.C, p data.Person) error {
	j, _ := json.Marshal(p.Details)

	_, err := s.exec(c, `INSERT INTO person (id, email, details) VALUES (?, ?, ?)`,
		p.Id, p.Email, j)
	if err == nil {
		return nil // inserted
	}

	_, err = s.exec(c, `UPDATE person SET email = ?, details = ? WHERE id = ?`,
		p.Email, j, p.Id)
	return err // updated
}

func (s *SQLStore) PersonById(c ctx.C, id string) (data.Person, bool, error) {
	p := data.Person{}
	rows, err := s.query(c, `SELECT id, email, details FROM person WHERE id = ?`, id)
	if err != nil {
		return p, false, ctx.Wrapf(c, "Query(%q): %v", id, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return p, false, nil
	}
	var j string
	err = rows.Scan(&p.Id, &p.Email, &j)
	if err != nil {
		return p, true, ctx.Wrapf(c, "can't scan db row: %v", err)
	}
	err = json.Unmarshal([]byte(j), &p.Details)
	if err != nil {
		return p, true, ctx.Wrapf(c, "can't unmarshal details: %v", err)
	}
	return p, true, nil
}
