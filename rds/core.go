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

type core struct {
	*sqlx.DB
	*config
}

func newCore(dbType Type, opts ...ConfigOptionFunc) (*core, error) {
	cfg := func() *config {
		cfg := newDefaultConfig()
		for _, opt := range opts {
			opt(cfg)
		}
		return cfg
	}()
	inst := &core{config: cfg}

	switch dbType {
	case Mysql:
		{
			inst.placeHolder = "?"
			path := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?collation=utf8mb4_general_ci", cfg.UserName, cfg.PassWord, cfg.IP, cfg.Port, cfg.DbName)
			db, dbErr := sqlx.Open("mysql", path)
			if dbErr != nil {
				return inst, dbErr
			}
			inst.DB = db
		}
	case Sqlite:
		{
			inst.placeHolder = "?"
			db, dbErr := sqlx.Open("sqlite3", cfg.DbName+".db")
			if dbErr != nil {
				return inst, dbErr
			}
			inst.DB = db
		}
	default:
		{
			return nil, errWrapper(IllegalParams)
		}
	}

	inst.DB.SetConnMaxLifetime(cfg.MaxLifeTime)
	inst.DB.SetMaxIdleConns(cfg.MaxIdleConns)

	return inst, inst.DB.Ping()
}

func (c *core) Stub() *Rds {
	return &Rds{c}
}

func (c *core) transaction(ctx context.Context, instructs ...*Instruct) (map[string]int64, error) {
	sid := func() string {
		_sid := ctx.Value("sid")
		if _sidString, _sidOk := _sid.(string); _sidOk {
			return _sidString
		} else {
			return strconv.Itoa(time.Now().Nanosecond())
		}
	}()

	c.config.Logger.Debugf("sid:%v,fn:transaction=>instructs:%#v", sid, instructsToString(instructs))
	tx, txErr := c.DB.BeginTx(ctx, nil)
	if txErr != nil {
		c.config.Logger.Errorf("sid:%v,fn:transaction=>txErr:%v", sid, txErr)
		return nil, txErr
	}
	var affectedRows = make(map[string]int64)
	execErrs := make([]error, len(instructs))
	for _, instruct := range instructs {
		c.config.Logger.Debugf("sid:%v,about to do ExecContext for instruct:%v", sid, instruct)
		execRst, execErr := c.DB.ExecContext(ctx, instruct.Sql, instruct.Args...)
		if execErr == nil {
			c.config.Logger.Debugf("sid:%v,ExecContext successfully,about to get RowsAffected", sid)
			affected, _ := execRst.RowsAffected()
			affectedRows[fmt.Sprintf("%v", instruct)] = affected
			c.config.Logger.Debugf("sid:%v,instruct exec successfully,affected:%v", affected, sid)
			continue
		}
		c.config.Logger.Errorf("sid:%v,fn:transaction=>execErr:%v", sid, execErr)
		execErrs = append(execErrs, execErr)
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			c.config.Logger.Errorf("sid:%v,fn:transaction=>rollbackErr:%v", sid, rollbackErr)
			return affectedRows, fmt.Errorf("execErr:%v,rollbackErr:%v", execErr, rollbackErr)
		}
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		c.config.Logger.Errorf("sid:%v,fn:transaction=>commitErr:%v", sid, commitErr)
		return affectedRows, fmt.Errorf("execErrs:%v,commitErr:%v", execErrs, commitErr)
	}
	return affectedRows, nil
}
