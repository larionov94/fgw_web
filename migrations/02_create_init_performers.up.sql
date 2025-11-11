-- СОЗДАТЬ ТАБЛИЦУ СОТРУДНИКОВ ДЛЯ ТЛК
CREATE TABLE dbo.svPerformers
(
    id              INT          DEFAULT 0  NOT NULL
        CONSTRAINT PK_svPerformers_id PRIMARY KEY,    -- id - табельный номер.
    fio             VARCHAR(150) DEFAULT '' NOT NULL, -- fio - ФИО сотрудника.
    bc              VARCHAR(13)  DEFAULT '' NOT NULL, -- bc - код доступа сотрудника.
    pass            VARCHAR(30)  DEFAULT '' NOT NULL, -- pass - пароль сотрудника.
    archive         BIT          DEFAULT 0  NOT NULL, -- archive - флаг архивного сотрудника.
    id_role_a_forms INT          DEFAULT 0  NOT NULL, -- id_role_a_forms - id роли.
    id_role_a_fgw   INT          DEFAULT 0  NOT NULL, -- id_role_fgw - id роли.
    created_at      DATETIME     DEFAULT GETDATE(),   -- created_at - дата создания записи.
    created_by      INT                     NOT NULL, -- created_by - табельный номер сотрудника.
    updated_at      DATETIME     DEFAULT GETDATE(),   -- updated_at - дата изменения записи.
    updated_by      INT                     NOT NULL, -- updated_by - табельный номер сотрудника изменивший запись.

    CONSTRAINT CHK_svPerformers_fio_not_empty CHECK (LEN(TRIM(fio)) > 0),
    CONSTRAINT FK_svPerformers_role_forms FOREIGN KEY (id_role_a_forms) REFERENCES dbo.svRoles (id),
    CONSTRAINT FK_svPerformers_role_fgw FOREIGN KEY (id_role_a_fgw) REFERENCES dbo.svRoles (id)
);

CREATE INDEX IX_svPerformers_bc ON dbo.svPerformers (bc);
CREATE INDEX IX_svPerformers_archive ON dbo.svPerformers (archive);
CREATE INDEX IX_svPerformers_fio ON dbo.svPerformers (fio);
CREATE INDEX IX_svPerformers_created_at ON dbo.svPerformers (created_at);

CREATE PROCEDURE dbo.svPerformerAll -- ХП получение всех сотрудников (только не архивных).
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id,
           fio,
           bc,
           pass,
           archive,
           id_role_a_forms,
           id_role_a_fgw,
           created_at,
           created_by,
           updated_at,
           updated_by
    FROM dbo.svPerformers
    WHERE archive = 0
END
GO;

CREATE PROCEDURE dbo.svPerformerAuth -- ХП проверяет сотрудника по табельному номеру и паролю для авторизации.
    @Id INT,
    @Pass VARCHAR(30)
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS(SELECT 1
              FROM dbo.svPerformers
              WHERE id = @Id
                AND pass = @Pass
                AND archive = 0)
        BEGIN
            SELECT 1 AS auth_success
        END
    ELSE
        BEGIN
            SELECT 0 AS auth_success
        END
END
GO;

CREATE PROCEDURE dbo.svPerformerFindById -- ХП ищет информацию о сотруднике по ИД.
@Id INT
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id,
           fio,
           bc,
           pass,
           archive,
           id_role_a_forms,
           id_role_a_fgw,
           created_at,
           created_by,
           updated_at,
           updated_by
    FROM dbo.svPerformers
    WHERE id = @Id
      AND archive = 0
END
GO;

CREATE PROCEDURE dbo.svPerformerUpdById -- ХП обновляет сотрудника по ИД.
    @Id INT,
    @Id_role_a_forms INT,
    @Id_role_a_fgw INT,
    @Updated_by INT
AS
BEGIN
    SET NOCOUNT ON;

    UPDATE dbo.svPerformers
    SET id_role_a_forms = @Id_role_a_forms,
        id_role_a_fgw   = @Id_role_a_fgw,
        updated_at      = GETDATE(),
        updated_by      = @Updated_by
    WHERE id = @Id
      AND archive = 0
END
GO;

CREATE PROCEDURE dbo.svPerformerExistsById -- ХП проверяет, существует ли сотрудник.
@Id INT
AS
BEGIN
    SET NOCOUNT ON;

    DECLARE @Exists BIT = 0;

    IF EXISTS(SELECT 1 FROM dbo.svPerformers WHERE id = @Id AND archive = 0)
        BEGIN
            SET @Exists = 1
        END

    SELECT @Exists AS exists_flag;
END
GO;