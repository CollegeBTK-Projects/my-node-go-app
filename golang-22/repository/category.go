package repository

import "strings"

func Category(err error) string {
	msg := err.Error()

	if strings.Contains(msg, "password authentication failed") {
		return "неверный пароль"
	}
	if strings.Contains(msg, "syntax error") {
		return "ошибка синтаксиса"
	}
	if strings.Contains(msg, "duplicate key") {
		return "дубликат"
	}
	if strings.Contains(msg, "null value") {
		return "NULL в NOT NULL"
	}
	if strings.Contains(msg, "foreign key") {
		return "нарушение внешнего ключа"
	}
	if strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "deadline exceeded") {
		return "сервер недоступен"
	}
	return "другая ошибка"
}
