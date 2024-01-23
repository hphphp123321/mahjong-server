// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/hphphp123321/mahjong-server/app/entity"
)

func newLog(db *gorm.DB, opts ...gen.DOOption) log {
	_log := log{}

	_log.logDo.UseDB(db, opts...)
	_log.logDo.UseModel(&entity.Log{})

	tableName := _log.logDo.TableName()
	_log.ALL = field.NewAsterisk(tableName)
	_log.ID = field.NewUint(tableName, "id")
	_log.CreatedAt = field.NewTime(tableName, "created_at")
	_log.UpdatedAt = field.NewTime(tableName, "updated_at")
	_log.DeletedAt = field.NewField(tableName, "deleted_at")
	_log.Content = field.NewString(tableName, "content")
	_log.Users = logManyToManyUsers{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Users", "entity.User"),
		Logs: struct {
			field.RelationField
			Users struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Users.Logs", "entity.Log"),
			Users: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Users.Logs.Users", "entity.User"),
			},
		},
	}

	_log.fillFieldMap()

	return _log
}

type log struct {
	logDo logDo

	ALL       field.Asterisk
	ID        field.Uint
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	Content   field.String
	Users     logManyToManyUsers

	fieldMap map[string]field.Expr
}

func (l log) Table(newTableName string) *log {
	l.logDo.UseTable(newTableName)
	return l.updateTableName(newTableName)
}

func (l log) As(alias string) *log {
	l.logDo.DO = *(l.logDo.As(alias).(*gen.DO))
	return l.updateTableName(alias)
}

func (l *log) updateTableName(table string) *log {
	l.ALL = field.NewAsterisk(table)
	l.ID = field.NewUint(table, "id")
	l.CreatedAt = field.NewTime(table, "created_at")
	l.UpdatedAt = field.NewTime(table, "updated_at")
	l.DeletedAt = field.NewField(table, "deleted_at")
	l.Content = field.NewString(table, "content")

	l.fillFieldMap()

	return l
}

func (l *log) WithContext(ctx context.Context) ILogDo { return l.logDo.WithContext(ctx) }

func (l log) TableName() string { return l.logDo.TableName() }

func (l log) Alias() string { return l.logDo.Alias() }

func (l log) Columns(cols ...field.Expr) gen.Columns { return l.logDo.Columns(cols...) }

func (l *log) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := l.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (l *log) fillFieldMap() {
	l.fieldMap = make(map[string]field.Expr, 6)
	l.fieldMap["id"] = l.ID
	l.fieldMap["created_at"] = l.CreatedAt
	l.fieldMap["updated_at"] = l.UpdatedAt
	l.fieldMap["deleted_at"] = l.DeletedAt
	l.fieldMap["content"] = l.Content

}

func (l log) clone(db *gorm.DB) log {
	l.logDo.ReplaceConnPool(db.Statement.ConnPool)
	return l
}

func (l log) replaceDB(db *gorm.DB) log {
	l.logDo.ReplaceDB(db)
	return l
}

type logManyToManyUsers struct {
	db *gorm.DB

	field.RelationField

	Logs struct {
		field.RelationField
		Users struct {
			field.RelationField
		}
	}
}

func (a logManyToManyUsers) Where(conds ...field.Expr) *logManyToManyUsers {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a logManyToManyUsers) WithContext(ctx context.Context) *logManyToManyUsers {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a logManyToManyUsers) Session(session *gorm.Session) *logManyToManyUsers {
	a.db = a.db.Session(session)
	return &a
}

func (a logManyToManyUsers) Model(m *entity.Log) *logManyToManyUsersTx {
	return &logManyToManyUsersTx{a.db.Model(m).Association(a.Name())}
}

type logManyToManyUsersTx struct{ tx *gorm.Association }

func (a logManyToManyUsersTx) Find() (result []*entity.User, err error) {
	return result, a.tx.Find(&result)
}

func (a logManyToManyUsersTx) Append(values ...*entity.User) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a logManyToManyUsersTx) Replace(values ...*entity.User) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a logManyToManyUsersTx) Delete(values ...*entity.User) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a logManyToManyUsersTx) Clear() error {
	return a.tx.Clear()
}

func (a logManyToManyUsersTx) Count() int64 {
	return a.tx.Count()
}

type logDo struct{ gen.DO }

type ILogDo interface {
	gen.SubQuery
	Debug() ILogDo
	WithContext(ctx context.Context) ILogDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ILogDo
	WriteDB() ILogDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ILogDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ILogDo
	Not(conds ...gen.Condition) ILogDo
	Or(conds ...gen.Condition) ILogDo
	Select(conds ...field.Expr) ILogDo
	Where(conds ...gen.Condition) ILogDo
	Order(conds ...field.Expr) ILogDo
	Distinct(cols ...field.Expr) ILogDo
	Omit(cols ...field.Expr) ILogDo
	Join(table schema.Tabler, on ...field.Expr) ILogDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ILogDo
	RightJoin(table schema.Tabler, on ...field.Expr) ILogDo
	Group(cols ...field.Expr) ILogDo
	Having(conds ...gen.Condition) ILogDo
	Limit(limit int) ILogDo
	Offset(offset int) ILogDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ILogDo
	Unscoped() ILogDo
	Create(values ...*entity.Log) error
	CreateInBatches(values []*entity.Log, batchSize int) error
	Save(values ...*entity.Log) error
	First() (*entity.Log, error)
	Take() (*entity.Log, error)
	Last() (*entity.Log, error)
	Find() ([]*entity.Log, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.Log, err error)
	FindInBatches(result *[]*entity.Log, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*entity.Log) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ILogDo
	Assign(attrs ...field.AssignExpr) ILogDo
	Joins(fields ...field.RelationField) ILogDo
	Preload(fields ...field.RelationField) ILogDo
	FirstOrInit() (*entity.Log, error)
	FirstOrCreate() (*entity.Log, error)
	FindByPage(offset int, limit int) (result []*entity.Log, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ILogDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (l logDo) Debug() ILogDo {
	return l.withDO(l.DO.Debug())
}

func (l logDo) WithContext(ctx context.Context) ILogDo {
	return l.withDO(l.DO.WithContext(ctx))
}

func (l logDo) ReadDB() ILogDo {
	return l.Clauses(dbresolver.Read)
}

func (l logDo) WriteDB() ILogDo {
	return l.Clauses(dbresolver.Write)
}

func (l logDo) Session(config *gorm.Session) ILogDo {
	return l.withDO(l.DO.Session(config))
}

func (l logDo) Clauses(conds ...clause.Expression) ILogDo {
	return l.withDO(l.DO.Clauses(conds...))
}

func (l logDo) Returning(value interface{}, columns ...string) ILogDo {
	return l.withDO(l.DO.Returning(value, columns...))
}

func (l logDo) Not(conds ...gen.Condition) ILogDo {
	return l.withDO(l.DO.Not(conds...))
}

func (l logDo) Or(conds ...gen.Condition) ILogDo {
	return l.withDO(l.DO.Or(conds...))
}

func (l logDo) Select(conds ...field.Expr) ILogDo {
	return l.withDO(l.DO.Select(conds...))
}

func (l logDo) Where(conds ...gen.Condition) ILogDo {
	return l.withDO(l.DO.Where(conds...))
}

func (l logDo) Order(conds ...field.Expr) ILogDo {
	return l.withDO(l.DO.Order(conds...))
}

func (l logDo) Distinct(cols ...field.Expr) ILogDo {
	return l.withDO(l.DO.Distinct(cols...))
}

func (l logDo) Omit(cols ...field.Expr) ILogDo {
	return l.withDO(l.DO.Omit(cols...))
}

func (l logDo) Join(table schema.Tabler, on ...field.Expr) ILogDo {
	return l.withDO(l.DO.Join(table, on...))
}

func (l logDo) LeftJoin(table schema.Tabler, on ...field.Expr) ILogDo {
	return l.withDO(l.DO.LeftJoin(table, on...))
}

func (l logDo) RightJoin(table schema.Tabler, on ...field.Expr) ILogDo {
	return l.withDO(l.DO.RightJoin(table, on...))
}

func (l logDo) Group(cols ...field.Expr) ILogDo {
	return l.withDO(l.DO.Group(cols...))
}

func (l logDo) Having(conds ...gen.Condition) ILogDo {
	return l.withDO(l.DO.Having(conds...))
}

func (l logDo) Limit(limit int) ILogDo {
	return l.withDO(l.DO.Limit(limit))
}

func (l logDo) Offset(offset int) ILogDo {
	return l.withDO(l.DO.Offset(offset))
}

func (l logDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ILogDo {
	return l.withDO(l.DO.Scopes(funcs...))
}

func (l logDo) Unscoped() ILogDo {
	return l.withDO(l.DO.Unscoped())
}

func (l logDo) Create(values ...*entity.Log) error {
	if len(values) == 0 {
		return nil
	}
	return l.DO.Create(values)
}

func (l logDo) CreateInBatches(values []*entity.Log, batchSize int) error {
	return l.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (l logDo) Save(values ...*entity.Log) error {
	if len(values) == 0 {
		return nil
	}
	return l.DO.Save(values)
}

func (l logDo) First() (*entity.Log, error) {
	if result, err := l.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Log), nil
	}
}

func (l logDo) Take() (*entity.Log, error) {
	if result, err := l.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Log), nil
	}
}

func (l logDo) Last() (*entity.Log, error) {
	if result, err := l.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Log), nil
	}
}

func (l logDo) Find() ([]*entity.Log, error) {
	result, err := l.DO.Find()
	return result.([]*entity.Log), err
}

func (l logDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.Log, err error) {
	buf := make([]*entity.Log, 0, batchSize)
	err = l.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (l logDo) FindInBatches(result *[]*entity.Log, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return l.DO.FindInBatches(result, batchSize, fc)
}

func (l logDo) Attrs(attrs ...field.AssignExpr) ILogDo {
	return l.withDO(l.DO.Attrs(attrs...))
}

func (l logDo) Assign(attrs ...field.AssignExpr) ILogDo {
	return l.withDO(l.DO.Assign(attrs...))
}

func (l logDo) Joins(fields ...field.RelationField) ILogDo {
	for _, _f := range fields {
		l = *l.withDO(l.DO.Joins(_f))
	}
	return &l
}

func (l logDo) Preload(fields ...field.RelationField) ILogDo {
	for _, _f := range fields {
		l = *l.withDO(l.DO.Preload(_f))
	}
	return &l
}

func (l logDo) FirstOrInit() (*entity.Log, error) {
	if result, err := l.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Log), nil
	}
}

func (l logDo) FirstOrCreate() (*entity.Log, error) {
	if result, err := l.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Log), nil
	}
}

func (l logDo) FindByPage(offset int, limit int) (result []*entity.Log, count int64, err error) {
	result, err = l.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = l.Offset(-1).Limit(-1).Count()
	return
}

func (l logDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = l.Count()
	if err != nil {
		return
	}

	err = l.Offset(offset).Limit(limit).Scan(result)
	return
}

func (l logDo) Scan(result interface{}) (err error) {
	return l.DO.Scan(result)
}

func (l logDo) Delete(models ...*entity.Log) (result gen.ResultInfo, err error) {
	return l.DO.Delete(models)
}

func (l *logDo) withDO(do gen.Dao) *logDo {
	l.DO = *do.(*gen.DO)
	return l
}
