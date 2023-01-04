package utils

const FmtDate = "2006-01-02"
const FmtDateTime = "2006-01-02 15:04:05"
const FmtDateTimeStr = "2006-01-02T15:04:05.000000Z"

const Postgres = "postgres"
const Mysql = "mysql"
const Mongo = "mongo"

type str string

// SpanStr must have a type, if not, it will collide when it is used in middleware/http.go line43. It is for context.WithValue()
const SpanStr str = "ParentSpan"
