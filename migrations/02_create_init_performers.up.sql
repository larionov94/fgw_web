-- СОЗДАТЬ ТАБЛИЦУ СОТРУДНИКОВ ДЛЯ ТЛК
CREATE TABLE dbo.svPerformers
(
    id              INT          DEFAULT 0  NOT NULL PRIMARY KEY, -- id - табельный номер.
    fio             VARCHAR(150) DEFAULT '' NOT NULL,             -- fio - ФИО сотрудника.
    bc              VARCHAR(13)  DEFAULT '' NOT NULL,             -- bc - код доступа сотрудника.
    pass            VARCHAR(30)  DEFAULT '' NOT NULL,             -- pass - пароль сотрудника.
    archive         BIT          DEFAULT 0  NOT NULL,             -- archive - флаг архивного сотрудника.
    id_role_a_forms INT          DEFAULT 0  NOT NULL,             -- id_role_a_forms - id роли.
    id_role_a_fgw   INT          DEFAULT 0  NOT NULL,             -- id_role_fgw - id роли.
);
CREATE INDEX idx_svPerformers_id ON dbo.svPerformers (id);
CREATE INDEX idx_svPerformers_id_role_a_forms ON dbo.svPerformers (id_role_a_forms);
CREATE INDEX idx_svPerformers_id_role_a_fgw ON dbo.svPerformers (id_role_a_fgw);

CREATE PROCEDURE dbo.svGetPerformers -- Процедура получения всех сотрудников.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, fio, bc, pass, archive, id_role_a_forms, id_role_a_fgw FROM dbo.svPerformers;
END
GO;

CREATE PROCEDURE dbo.svGetPerformerById -- Процедура получения сотрудника по id.
@Id INT -- id сотрудника.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, fio, bc, pass, archive, id_role_a_forms, id_role_a_fgw FROM dbo.svPerformers WHERE id = @Id;
END
GO;

CREATE PROCEDURE dbo.svGetPerformerByIdRoleAForms -- Процедура получения сотрудника по id роли.
@IdRoleAForms INT -- id роли.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, fio, bc, pass, archive, id_role_a_forms, id_role_a_fgw
    FROM dbo.svPerformers
    WHERE id_role_a_forms = @IdRoleAForms;
END
GO;

CREATE PROCEDURE dbo.svGetPerformerByIdRoleFGW -- Процедура получения сотрудника по id роли.
@IdRoleFGW INT -- id роли.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, fio, bc, pass, archive, id_role_a_forms, id_role_a_fgw
    FROM dbo.svPerformers
    WHERE id_role_a_fgw = @IdRoleFGW;
END
GO;

-- exec dbo.svGetPerformers;
-- exec dbo.svGetPerformerById 123456;
-- exec dbo.svGetPerformerByIdRole 1;

