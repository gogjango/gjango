package mockgopg

import (
	"strings"
	"sync"
)

// SQLMock handles query mocks
type SQLMock struct {
	lock          *sync.RWMutex
	currentQuery  string
	currentParams []interface{}
	queries       map[string]buildQuery
}

// ExpectExec is a builder method that accepts a query in string and returns an SQLMock pointer
func (sqlMock *SQLMock) ExpectExec(query string) *SQLMock {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	sqlMock.currentQuery = strings.TrimSpace(query)
	return sqlMock
}

// ExpectQuery accepts a query in string and returns an SQLMock pointer
func (sqlMock *SQLMock) ExpectQuery(query string) *SQLMock {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	sqlMock.currentQuery = strings.TrimSpace(query)
	return sqlMock
}

// WithArgs is a builder method that accepts a query in string and returns an SQLMock pointer
func (sqlMock *SQLMock) WithArgs(params ...interface{}) *SQLMock {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	sqlMock.currentParams = make([]interface{}, 0)
	for _, p := range params {
		sqlMock.currentParams = append(sqlMock.currentParams, p)
	}

	return sqlMock
}

// Returns accepts expected result and error, and completes the build of our sqlMock object
func (sqlMock *SQLMock) Returns(result *OrmResult, err error) {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	q := buildQuery{
		query:  sqlMock.currentQuery,
		params: sqlMock.currentParams,
		result: result,
		err:    err,
	}

	sqlMock.queries[sqlMock.currentQuery] = q
	sqlMock.currentQuery = ""
	sqlMock.currentParams = nil
}

// FlushAll resets our sqlMock object
func (sqlMock *SQLMock) FlushAll() {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	sqlMock.currentQuery = ""
	sqlMock.currentParams = nil
	sqlMock.queries = make(map[string]buildQuery)
}
