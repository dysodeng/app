package db

import (
	"fmt"
	"log"
	"time"

	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	otelTrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// TracerPlugin 数据库tracer插件
type TracerPlugin struct{}

func (op *TracerPlugin) Name() string {
	return "tracer"
}

func (op *TracerPlugin) Initialize(db *gorm.DB) error {
	_ = db.Callback().Create().Before("gorm:before_create").Register("tracer:before_create", op.beforeCreate)
	_ = db.Callback().Create().After("gorm:after_create").Register("tracer:after_create", op.afterCreate)
	_ = db.Callback().Query().Before("gorm:before_query").Register("tracer:before_query", op.beforeQuery)
	_ = db.Callback().Query().After("gorm:after_query").Register("tracer:after_query", op.afterQuery)
	_ = db.Callback().Update().Before("gorm:before_update").Register("tracer:before_update", op.beforeUpdate)
	_ = db.Callback().Update().After("gorm:after_update").Register("tracer:after_update", op.afterUpdate)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register("tracer:before_delete", op.beforeDelete)
	_ = db.Callback().Delete().After("gorm:after_delete").Register("tracer:after_delete", op.afterDelete)
	_ = db.Callback().Row().Before("gorm:before_row").Register("tracer:before_row", op.beforeRow)
	_ = db.Callback().Row().After("gorm:after_row").Register("tracer:after_row", op.afterRow)
	_ = db.Callback().Raw().Before("gorm:before_raw").Register("tracer:before_raw", op.beforeRaw)
	_ = db.Callback().Raw().After("gorm:after_raw").Register("tracer:after_raw", op.afterRaw)
	return nil
}

func (op *TracerPlugin) beforeCreate(db *gorm.DB) {
	op.startSpan(db, "create")
}

func (op *TracerPlugin) afterCreate(db *gorm.DB) {
	op.endSpan(db, "create")
}

func (op *TracerPlugin) beforeQuery(db *gorm.DB) {
	op.startSpan(db, "query")
}

func (op *TracerPlugin) afterQuery(db *gorm.DB) {
	op.endSpan(db, "query")
}

func (op *TracerPlugin) beforeUpdate(db *gorm.DB) {
	op.startSpan(db, "update")
}

func (op *TracerPlugin) afterUpdate(db *gorm.DB) {
	op.endSpan(db, "update")
}

func (op *TracerPlugin) beforeDelete(db *gorm.DB) {
	op.startSpan(db, "delete")
}

func (op *TracerPlugin) afterDelete(db *gorm.DB) {
	op.endSpan(db, "delete")
}

func (op *TracerPlugin) beforeRow(db *gorm.DB) {
	op.startSpan(db, "row")
}

func (op *TracerPlugin) afterRow(db *gorm.DB) {
	op.endSpan(db, "row")
}

func (op *TracerPlugin) beforeRaw(db *gorm.DB) {
	op.startSpan(db, "raw")
}

func (op *TracerPlugin) afterRaw(db *gorm.DB) {
	op.endSpan(db, "raw")
}

func (op *TracerPlugin) startSpan(db *gorm.DB, operation string) {
	ctx, _ := trace.Tracer().Start(db.Statement.Context, "gorm:"+operation)
	db.Statement.Context = ctx
	db.Statement.Set("startTime", time.Now())
}

func (op *TracerPlugin) endSpan(db *gorm.DB, operation string) {
	span := otelTrace.SpanFromContext(db.Statement.Context)
	if span == nil {
		log.Println("span is nil")
		return
	}

	startTime, ok := db.Statement.Get("startTime")
	if !ok {
		log.Println("startTime not found")
		span.End()
		return
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime.(time.Time))

	// 记录操作是否成功
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		span.SetStatus(codes.Error, db.Error.Error())
		span.RecordError(db.Error)
	} else {
		span.SetStatus(codes.Ok, "ok")
	}

	span.SetAttributes(
		attribute.String("sql.query", db.Statement.SQL.String()),
		attribute.String("sql.args", fmt.Sprintf("%+v", db.Statement.Vars)),
		attribute.Int64("sql.rows_affected", db.RowsAffected),         // 记录操作记录数
		attribute.Float64("sql.duration_ms", duration.Seconds()*1000), // 记录执行时间，单位为毫秒
	)

	span.End()
}
