package rds

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

type Rds struct {
	*sqlx.DB
	cfg *Config
}

func NewRds(dbType Type, opts ...ConfigOptionFunc) (*Rds, error) {
	config := func() *Config {
		cfg := newConfig()
		for _, opt := range opts {
			opt(cfg)
		}
		return cfg
	}()
	rdbInst := &Rds{cfg: func() *Config {
		cfg := newConfig()
		for _, opt := range opts {
			opt(cfg)
		}
		return cfg
	}()}
	{
		db, dbErr := func(dbType Type, dbConfig *Config) (*sqlx.DB, error) {
			switch dbType {
			case Mysql:
				{
					return func(m *Config) (*sqlx.DB, error) {
						path := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?collation=utf8mb4_general_ci", m.UserName, m.PassWord, m.IP, m.Port, m.DbName)
						return sqlx.Open("mysql", path)
					}(dbConfig)
				}
			case Sqlite:
				{
					return func(m *Config) (*sqlx.DB, error) {
						return sqlx.Open("sqlite3", m.DbName+".db")
					}(dbConfig)
				}
			default:
				{
					return nil, errWrapper(IllegalParams)
				}
			}
		}(dbType, config)
		if dbErr != nil {
			return nil, dbErr
		}
		rdbInst.DB = db
	}
	{
		rdbInst.DB.SetConnMaxLifetime(config.MaxLifeTime)
		rdbInst.DB.SetMaxIdleConns(config.MaxIdleConns)
	}
	{
		pingErr := rdbInst.DB.Ping()
		if pingErr != nil {
			return nil, pingErr
		}
	}
	return rdbInst, nil
}

func (r *Rds) transaction(ctx context.Context, instructs ...*Instruct) (map[string]int64, error) {
	sid := func() string {
		_sid := ctx.Value("sid")
		if _sidString, _sidOk := _sid.(string); _sidOk {
			return _sidString
		} else {
			return strconv.Itoa(time.Now().Nanosecond())
		}
	}()

	r.cfg.Logger.Debugf("sid:%v,fn:transaction=>instructs:%#v", sid, InstructsToString(instructs))
	tx, txErr := r.DB.BeginTx(ctx, nil)
	if txErr != nil {
		r.cfg.Logger.Errorf("sid:%v,fn:transaction=>txErr:%v", sid, txErr)
		return nil, txErr
	}
	var affectedRows = make(map[string]int64)
	execErrs := make([]error, len(instructs))
	for _, instruct := range instructs {
		r.cfg.Logger.Debugf("sid:%v,about to do ExecContext for instruct:%v", sid, instruct)
		execRst, execErr := r.DB.ExecContext(ctx, instruct.Sql, instruct.Args...)
		if execErr == nil {
			r.cfg.Logger.Debugf("sid:%v,ExecContext successfully,about to get RowsAffected", sid)
			affected, _ := execRst.RowsAffected()
			affectedRows[fmt.Sprintf("%v", instruct)] = affected
			r.cfg.Logger.Debugf("sid:%v,instruct exec successfully,affected:%v", affected, sid)
			continue
		}
		r.cfg.Logger.Errorf("sid:%v,fn:transaction=>execErr:%v", sid, execErr)
		execErrs = append(execErrs, execErr)
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			r.cfg.Logger.Errorf("sid:%v,fn:transaction=>rollbackErr:%v", sid, rollbackErr)
			return affectedRows, fmt.Errorf("execErr:%v,rollbackErr:%v", execErr, rollbackErr)
		}
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		r.cfg.Logger.Errorf("sid:%v,fn:transaction=>commitErr:%v", sid, commitErr)
		return affectedRows, fmt.Errorf("execErrs:%v,commitErr:%v", execErrs, commitErr)
	}
	return affectedRows, nil
}
