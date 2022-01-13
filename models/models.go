package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/leffss/go-admin/pkg/common"
	"github.com/leffss/go-admin/pkg/logging"
	"github.com/leffss/go-admin/pkg/setting"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

var DB *gorm.DB
var databaseSetting = setting.GetDatabaseSetting()
var logger = logging.GetLogger()

type GormLogger struct {
	gormLogger.Interface
	logger *zap.Logger
	logLevel gormLogger.LogLevel
	SlowThreshold time.Duration
	IgnoreRecordNotFoundError bool
}

func NewGormLogger(logger *zap.Logger, logLevel gormLogger.LogLevel, slowThreshold time.Duration, ignoreRecordNotFoundError bool) *GormLogger {
	return &GormLogger{
		logger: logger,
		logLevel: logLevel,
		SlowThreshold: slowThreshold,
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
	}
}

func (gl *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface  {
	return gl
}

func (gl *GormLogger) Info(ctx context.Context, msg string, values ...interface{}) {
	if gl.logLevel >= gormLogger.Info {
		gl.logger.Info(fmt.Sprintf("%s %s %v", utils.FileWithLineNum(), msg, values))
	}
}

func (gl *GormLogger) Warn(ctx context.Context, msg string, values ...interface{}) {
	if gl.logLevel >= gormLogger.Warn {
		gl.logger.Warn(fmt.Sprintf("%s %s %v", utils.FileWithLineNum(), msg, values))
	}
}

func (gl *GormLogger) Error(ctx context.Context, msg string, values ...interface{}) {
	if gl.logLevel >= gormLogger.Error {
		gl.logger.Error(fmt.Sprintf("%s %s %v", utils.FileWithLineNum(), msg, values))
	}
}

func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if gl.logLevel <= gormLogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && gl.logLevel >= gormLogger.Error && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !gl.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			gl.logger.Error(fmt.Sprint(utils.FileWithLineNum(), "##", err, "##", float64(elapsed.Nanoseconds())/1e6, "##", "-", "##", sql))
		} else {
			gl.logger.Error(fmt.Sprint(utils.FileWithLineNum(), "##", err, "##", float64(elapsed.Nanoseconds())/1e6, "##", rows, "##", sql))
		}
	case elapsed > gl.SlowThreshold && gl.SlowThreshold != 0 && gl.logLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", gl.SlowThreshold)
		if rows == -1 {
			gl.logger.Warn(fmt.Sprint(utils.FileWithLineNum(), "##", slowLog, "##", float64(elapsed.Nanoseconds())/1e6, "##", "-", "##", sql))
		} else {
			gl.logger.Warn(fmt.Sprint(utils.FileWithLineNum(), "##", slowLog, "##", float64(elapsed.Nanoseconds())/1e6, "##", rows, "##", sql))
		}
	case gl.logLevel == gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			gl.logger.Debug(fmt.Sprint(utils.FileWithLineNum(), "##", float64(elapsed.Nanoseconds())/1e6, "##", "-", "##", sql))
		} else {
			gl.logger.Debug(fmt.Sprint(utils.FileWithLineNum(), "##", float64(elapsed.Nanoseconds())/1e6, "##", rows, "##", sql))
		}
	}
}

func Setup() (err error) {
	config := &gorm.Config{
		Logger: NewGormLogger(logger, gormLogger.Error, 200 * time.Millisecond, true),
	}
	if databaseSetting.Type == "mysql" {
		//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", databaseSetting.User,
		//	databaseSetting.Password, databaseSetting.Host, databaseSetting.Port, databaseSetting.Db)
		timeZone := strings.Replace(databaseSetting.TimeZone, "/", "%2F", -1)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", databaseSetting.User,
			databaseSetting.Password, databaseSetting.Host, databaseSetting.Port, databaseSetting.Db, timeZone)
		// method 1
		//db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		// method 2
		DB, err = gorm.Open(mysql.New(mysql.Config{
			DSN: dsn, // data source name
			DefaultStringSize: 256, // default size for string fields
			DisableDatetimePrecision: true, // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex: true, // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn: true, // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
		}), config)
		if err != nil {
			return err
		}
	} else if databaseSetting.Type == "postgresql" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", databaseSetting.Host,
			databaseSetting.User, databaseSetting.Password, databaseSetting.Db, databaseSetting.Port, databaseSetting.TimeZone)
		DB, err = gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), config)
		if err != nil {
			return err
		}
	} else if databaseSetting.Type == "sqlserver" {
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", databaseSetting.User,
			databaseSetting.Password, databaseSetting.Host, databaseSetting.Port, databaseSetting.Db)
		DB, err = gorm.Open(sqlserver.Open(dsn), config)
		if err != nil {
			return err
		}
	} else if databaseSetting.Type == "sqlite3" {
		DB, err = gorm.Open(sqlite.Open(databaseSetting.DdFile), config)
		if err != nil {
			return err
		}
	} else {
		return errors.New("error database type setting")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(databaseSetting.ConnMaxLifetime)

	return autoMigrate()
}

func autoMigrate() error {
	if err := DB.AutoMigrate(&Department{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&Role{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&User{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&Permission{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&LoginLog{}); err != nil {
		return err
	}

	return nil
}

func InitDatabase(password string)  {
	// 创建权限
	p1 := &Permission{
		Name: "测试权限1",
		Sign: "test1",
		Menu: true,
		Method: "get",
		Path: "/test1",
	}
	p2 := &Permission{
		Name: "测试权限2",
		Sign: "test2",
		Menu: true,
		Method: "get",
		Path: "/test2",
	}
	p3 := &Permission{
		Name: "测试权限3",
		Sign: "test3",
		Menu: true,
		Method: "get",
		Path: "/test3",
	}
	p31 := &Permission{
		Name: "测试权限31",
		Sign: "test31",
		Menu: false,
		Method: "post",
		Path: "/test3/1",
	}
	DB.Create(p1)
	DB.Create(p2)
	DB.Create(p3)
	p31.ParentId = p3.Id
	DB.Create(p31)

	// 添加超级管理员
	u1 := &User{
		Username:      "admin",
		Password:      common.EncodeSHA256(password),
		IsSuper:       true,
		Status:        true,
	}
	DB.Create(u1)

	// 添加默认部门
	d1 := &Department{
		Name:          "默认",
		ControlBy:     ControlBy{CreateBy: u1.Id, UpdateBy: u1.Id},
	}
	DB.Create(d1)

	// 添加角色
	r1 := &Role{
		Name:          "管理员",
		DeptId:        d1.Id,
		ControlBy:     ControlBy{CreateBy: u1.Id, UpdateBy: u1.Id},
	}
	r1.Permissions = append(r1.Permissions, *p1)
	r1.Permissions = append(r1.Permissions, *p2)
	r1.Permissions = append(r1.Permissions, *p3)
	r1.Permissions = append(r1.Permissions, *p31)

	DB.Create(r1)
	// 清除角色权限
	//db.Model(r1).Association("Permissions").Delete(p1)
	//db.Model(r1).Association("Permissions").Clear()

	r2 := &Role{
		Name:          "用户",
		DeptId:        d1.Id,
		ControlBy:     ControlBy{CreateBy: u1.Id, UpdateBy: u1.Id},
	}
	DB.Create(r2)
	err := DB.Model(r2).Association("Permissions").Append(p1, p2)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 添加普通用户
	u2 := &User{
		Username:      "leffss",
		Password:      common.EncodeSHA256(password),
		Status:        true,
		RoleId:        r2.Id,
	}
	DB.Create(u2)
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
