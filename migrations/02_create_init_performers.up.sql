-- СОЗДАТЬ ТАБЛИЦУ СОТРУДНИКОВ ДЛЯ ТЛК
CREATE TABLE dbo.svPerformers
(
    id      INT          DEFAULT 0        NOT NULL PRIMARY KEY, -- id - табельный номер.
    fio     VARCHAR(150) DEFAULT ''       NOT NULL,             -- fio - ФИО сотрудника.
    bc      VARCHAR(13)  DEFAULT ''       NOT NULL,             -- bc - код доступа сотрудника.
    pass    VARCHAR(30)  DEFAULT '123456' NOT NULL,             -- pass - пароль сотрудника.
    archive BIT          DEFAULT 0        NOT NULL,             -- archive - флаг архивного сотрудника.
    id_role INT          DEFAULT 0        NOT NULL,             -- id_role - id роли.
);
CREATE INDEX idx_svPerformers_id ON dbo.svPerformers (id);
CREATE INDEX idx_svPerformers_id_role ON dbo.svPerformers (id_role);

CREATE PROCEDURE dbo.svAllPerformers -- Процедура получения всех сотрудников.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, fio, bc, pass, archive, id_role FROM dbo.svPerformers;
END
GO;

CREATE PROCEDURE dbo.svGetPerformerById -- Процедура получения сотрудника по id.
    @Id INT -- id сотрудника.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, fio, bc, pass, archive, id_role FROM dbo.svPerformers WHERE id = @Id;
END
GO;

CREATE PROCEDURE dbo.svGetPerformerByIdRole -- Процедура получения сотрудника по id роли.
    @IdRole INT -- id роли.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, fio, bc, pass, archive, id_role FROM dbo.svPerformers WHERE id_role = @IdRole;
END
GO;

-- exec dbo.svGetPerformers;
-- exec dbo.svGetPerformerById 123456;
-- exec dbo.svGetPerformerByIdRole 1;

