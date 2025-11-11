package repository

import (
	"FGW_WEB/internal/config/db"
	"FGW_WEB/internal/model"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

// createMock создание мок.
func createMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PerformerRepo) {
	mssqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	realLogger := &common.Logger{}

	// создаем репозиторий с моком
	repo := NewPerformerRepo(mssqlDB, realLogger)

	return mssqlDB, mock, repo
}

func getRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "fio", "bc", "pass", "archive",
		"id_role_a_forms", "id_role_a_fgw",
		"created_by", "created_at", "updated_by", "updated_at",
	})
}

// getPerformerNotArchive возвращает тестовые данные не архивного сотрудника
func getPerformerNotArchiveInvalid() (string, string, string, string, bool, int, int, int, string, int, string) {
	now := time.Now().String()
	return "invalid_id", "Тестов Тест Тестович", "000012345", "12345", false, 0, 0, 12345, now, 12345, now
}

// getPerformerNotArchive возвращает тестовые данные не архивного сотрудника
func getPerformerNotArchive() (int, string, string, string, bool, int, int, int, string, int, string) {
	now := time.Now().String()
	return 12345, "Тестов Тест Тестович", "000012345", "12345", false, 0, 0, 12345, now, 12345, now
}

// getPerformerArchive возвращает тестовые данные архивного сотрудника
func getPerformerArchive() (int, string, string, string, bool, int, int, int, string, int, string) {
	now := time.Now().String()
	return 12345, "Тест Тестов Тестович", "000012345", "12345", true, 0, 0, 12345, now, 12345, now
}

// verifyPerformer проверяет основные поля сотрудника.
func verifyPerformer(t *testing.T, performer *model.Performer, expectedFIO string) {
	t.Helper()

	require.NotNil(t, performer)
	assert.Equal(t, 12345, performer.Id)
	assert.Equal(t, "000012345", performer.BC)
	assert.Equal(t, "12345", performer.Pass)
	assert.False(t, performer.Archive)
	assert.Equal(t, 0, performer.IdRoleAForms)
	assert.Equal(t, 0, performer.IdRoleAFGW)
	assert.Equal(t, 12345, performer.AuditRec.CreatedBy)
	assert.Equal(t, 12345, performer.AuditRec.UpdatedBy)

	expectedFIOConverted, _ := convert.Win1251ToUTF8(expectedFIO)
	assert.Equal(t, expectedFIOConverted, performer.FIO)
}

// verifyArchivePerformer проверяет архивного сотрудника.
func verifyArchivePerformer(t *testing.T, performer *model.Performer, expectedFIO string) {
	t.Helper()

	require.NotNil(t, performer)
	assert.True(t, performer.Archive)

	expectedFIOConverted, _ := convert.Win1251ToUTF8(expectedFIO)
	assert.Equal(t, expectedFIOConverted, performer.FIO)
}

func TestPerformerRepo_All(t *testing.T) {
	mssqlDB, mock, repo := createMock(t)

	defer db.Close(mssqlDB)

	t.Run("Успех - возвращаем сотрудников", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAllQuery
		expectedRows := getRows().AddRow(getPerformerNotArchive())

		// ожидаем вызов запроса
		mock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		performers, err := repo.All(context.Background())

		require.NoError(t, err)
		require.Len(t, performers, 1)

		verifyPerformer(t, &performers[0], "Тестов Тест Тестович")

		// проверяем что ожидания все выполнены
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Успех - пустой результат", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAllQuery

		emptyRows := getRows()
		mock.ExpectQuery(expectedQuery).WillReturnRows(emptyRows)

		performers, err := repo.All(context.Background())

		require.NoError(t, err)
		assert.Empty(t, performers)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - сбой выполнения запроса", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAllQuery
		expectedError := errors.New("database connection failed")

		mock.ExpectQuery(expectedQuery).WillReturnError(expectedError)

		performers, err := repo.All(context.Background())

		assert.Error(t, err)
		assert.Nil(t, performers)
		assert.Contains(t, err.Error(), msg.E3202)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - сбой сканирования строки", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAllQuery

		// создаем невалидную строку
		invalidRows := getRows().AddRow(getPerformerNotArchiveInvalid())

		mock.ExpectQuery(expectedQuery).WillReturnRows(invalidRows)

		performers, err := repo.All(context.Background())

		assert.Error(t, err)
		assert.Nil(t, performers)
		assert.Contains(t, err.Error(), msg.E3204)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - сбой итерации строк", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAllQuery
		expectedError := errors.New("rows iteration error")

		expectedRows := getRows().AddRow(getPerformerNotArchive()).RowError(0, expectedError)

		mock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		performers, err := repo.All(context.Background())

		assert.Error(t, err)
		assert.Nil(t, performers)
		assert.Contains(t, err.Error(), msg.E3205)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPerformerRepo_AuthByIdAndPass(t *testing.T) {
	mssqlDB, mock, repo := createMock(t)
	defer db.Close(mssqlDB)

	t.Run("Успех - аутентификация завершена", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAuthQuery
		expectedRows := sqlmock.NewRows([]string{"auth_success"}).AddRow(true)

		mock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		result, err := repo.AuthByIdAndPass(context.Background(), 12345, "12345")

		require.NoError(t, err)
		assert.True(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Неудача - неверный пароль", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAuthQuery
		expectedRows := sqlmock.NewRows([]string{"auth_success"}).AddRow(false)

		mock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		result, err := repo.AuthByIdAndPass(context.Background(), 12345, "invalid_password")

		require.NoError(t, err)
		assert.False(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - запись не найдена", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAuthQuery

		mock.ExpectQuery(expectedQuery).WillReturnError(sql.ErrNoRows)

		result, err := repo.AuthByIdAndPass(context.Background(), 12345, "password")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Contains(t, err.Error(), msg.E3202)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - не удалось подключиться к БД", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAuthQuery
		expectedError := errors.New("database connection failed")

		mock.ExpectQuery(expectedQuery).WillReturnError(expectedError)

		result, err := repo.AuthByIdAndPass(context.Background(), 12345, "12345")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Contains(t, err.Error(), msg.E3202)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - неверный тип данных в результате", func(t *testing.T) {
		expectedQuery := FGWsvPerformerAuthQuery

		rows := sqlmock.NewRows([]string{"auth_success"}).AddRow("not_a_boolean")
		mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

		result, err := repo.AuthByIdAndPass(context.Background(), 12345, "12345")

		assert.Error(t, err)
		assert.False(t, result)
		assert.Contains(t, err.Error(), msg.E3202)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPerformerRepo_FindById(t *testing.T) {
	mssqlDB, mock, repo := createMock(t)
	defer db.Close(mssqlDB)

	t.Run("Успех - сотрудник найден", func(t *testing.T) {
		expectedId := 12345
		expectedFIO := "Тест Тестов Тестович"
		createdAt := time.Now()
		updatedAt := time.Now()

		rows := getRows().
			AddRow(expectedId, expectedFIO, "000012345", "12345", false, 0, 0, 12345, createdAt, 12345, updatedAt)

		mock.ExpectQuery(FGWsvPerformerFindByIdQuery).
			WithArgs(sql.Named("id", expectedId)).
			WillReturnRows(rows)

		performer, err := repo.FindById(context.Background(), expectedId)

		require.NoError(t, err)
		require.NotNil(t, performer)

		verifyPerformer(t, performer, expectedFIO)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - сотрудник не найден", func(t *testing.T) {
		expectedQuery := FGWsvPerformerFindByIdQuery

		mock.ExpectQuery(expectedQuery).
			WithArgs(sql.Named("id", 0)).
			WillReturnError(sql.ErrNoRows)

		performer, err := repo.FindById(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, performer)
		assert.Contains(t, err.Error(), msg.E3206)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - не удалось подключиться к БД", func(t *testing.T) {
		expectedQuery := FGWsvPerformerFindByIdQuery
		expectedError := errors.New("database connection failed")

		mock.ExpectQuery(expectedQuery).
			WithArgs(sql.Named("id", 12345)).
			WillReturnError(expectedError)

		performer, err := repo.FindById(context.Background(), 12345)

		assert.Error(t, err)
		assert.Nil(t, performer)
		assert.Contains(t, err.Error(), msg.E3202)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Успех - архивный сотрудник найден", func(t *testing.T) {
		expectedQuery := FGWsvPerformerFindByIdQuery

		expectedRows := getRows().AddRow(getPerformerArchive())

		mock.ExpectQuery(expectedQuery).
			WithArgs(sql.Named("id", 12345)).
			WillReturnRows(expectedRows)

		performer, err := repo.FindById(context.Background(), 12345)

		require.NoError(t, err)
		require.NotNil(t, performer)
		verifyArchivePerformer(t, performer, "Тест Тестов Тестович")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPerformerRepo_UpdById(t *testing.T) {
	mssqlDB, mock, repo := createMock(t)
	defer db.Close(mssqlDB)

	t.Run("Успех - обновлен сотрудник", func(t *testing.T) {
		expectedQuery := FGWsvPerformerUpdByIdQuery
		expectedId := 12345

		performer := &model.Performer{
			IdRoleAForms: 1,
			IdRoleAFGW:   2,
			AuditRec: model.Audit{
				UpdatedBy: 12354,
			},
		}

		mock.ExpectExec(expectedQuery).
			WithArgs(sql.Named("id", expectedId),
				sql.Named("id_role_a_forms", performer.IdRoleAForms),
				sql.Named("id_role_a_fgw", performer.IdRoleAFGW),
				sql.Named("updated_by", performer.AuditRec.UpdatedBy),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdById(context.Background(), expectedId, performer)

		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - не удалось выполнить запрос к БД", func(t *testing.T) {
		expectedId := 12345
		expectedQuery := FGWsvPerformerUpdByIdQuery
		expectedError := errors.New("database connection failed")

		performer := &model.Performer{
			IdRoleAForms: 1,
			IdRoleAFGW:   2,
			AuditRec: model.Audit{
				UpdatedBy: 12354,
			},
		}
		mock.ExpectExec(expectedQuery).
			WithArgs(sql.Named("id", expectedId),
				sql.Named("id_role_a_forms", performer.IdRoleAForms),
				sql.Named("id_role_a_fgw", performer.IdRoleAFGW),
				sql.Named("updated_by", performer.AuditRec.UpdatedBy),
			).
			WillReturnError(expectedError)

		err := repo.UpdById(context.Background(), expectedId, performer)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), msg.E3202)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - сбой сканирования строки", func(t *testing.T) {
		expectedId := 12345
		expectedQuery := FGWsvPerformerUpdByIdQuery

		performer := &model.Performer{
			IdRoleAForms: 1,
			IdRoleAFGW:   2,
			AuditRec: model.Audit{
				UpdatedBy: 12354,
			},
		}

		mockResult := sqlmock.NewErrorResult(errors.New("rows affected error"))
		mock.ExpectExec(expectedQuery).
			WithArgs(sql.Named("id", expectedId),
				sql.Named("id_role_a_forms", performer.IdRoleAForms),
				sql.Named("id_role_a_fgw", performer.IdRoleAFGW),
				sql.Named("updated_by", performer.AuditRec.UpdatedBy),
			).
			WillReturnResult(mockResult)

		err := repo.UpdById(context.Background(), expectedId, performer)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), msg.E3207)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - сотрудник не найден", func(t *testing.T) {
		expectedQuery := FGWsvPerformerUpdByIdQuery
		expectedId := 0

		performer := &model.Performer{
			IdRoleAForms: 1,
			IdRoleAFGW:   2,
			AuditRec: model.Audit{
				UpdatedBy: 12354,
			},
		}

		mock.ExpectExec(expectedQuery).
			WithArgs(sql.Named("id", expectedId),
				sql.Named("id_role_a_forms", performer.IdRoleAForms),
				sql.Named("id_role_a_fgw", performer.IdRoleAFGW),
				sql.Named("updated_by", performer.AuditRec.UpdatedBy),
			).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdById(context.Background(), expectedId, performer)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), msg.E3208)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPerformerRepo_ExistById(t *testing.T) {
	mssqlDB, mock, repo := createMock(t)
	defer db.Close(mssqlDB)

	t.Run("Успех - сотрудник существует", func(t *testing.T) {
		expectedId := 12345
		expectedQuery := FGWsvPerformerExistsByIdQuery

		expectedRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery(expectedQuery).
			WithArgs(sql.Named("id", expectedId)).
			WillReturnRows(expectedRows)

		exists, err := repo.ExistById(context.Background(), expectedId)

		require.NoError(t, err)
		assert.True(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Успех - сотрудник не существует", func(t *testing.T) {
		expectedId := 12345
		expectedQuery := FGWsvPerformerExistsByIdQuery

		expectedRows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery(expectedQuery).
			WithArgs(sql.Named("id", expectedId)).
			WillReturnRows(expectedRows)

		exists, err := repo.ExistById(context.Background(), expectedId)

		require.NoError(t, err)
		assert.False(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - запись не найдена", func(t *testing.T) {
		expectedQuery := FGWsvPerformerExistsByIdQuery
		expectedId := 12345

		mock.ExpectQuery(expectedQuery).
			WithArgs(sql.Named("id", expectedId)).
			WillReturnError(sql.ErrNoRows)

		exists, err := repo.ExistById(context.Background(), expectedId)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), msg.E3206)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - проблемы с подключением БД", func(t *testing.T) {
		expectedQuery := FGWsvPerformerExistsByIdQuery
		expectedId := 12345
		expectedError := errors.New("database connection failed")

		mock.ExpectQuery(expectedQuery).
			WithArgs(sql.Named("id", expectedId)).
			WillReturnError(expectedError)

		exists, err := repo.ExistById(context.Background(), expectedId)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), msg.E3206)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Ошибка - неверный тип данных", func(t *testing.T) {
		expectedQuery := FGWsvPerformerExistsByIdQuery
		expectedId := 12345

		expectedRows := sqlmock.NewRows([]string{"exists"}).AddRow("not_a_boolean")
		mock.ExpectQuery(expectedQuery).
			WithArgs(sql.Named("id", expectedId)).
			WillReturnRows(expectedRows)

		exists, err := repo.ExistById(context.Background(), expectedId)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), msg.E3206)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
