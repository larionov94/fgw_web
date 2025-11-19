-- СОЗДАТЬ ТАБЛИЦУ РОЛЕЙ ДЛЯ СОТРУДНИКОВ
CREATE TABLE dbo.svRoles
(
    id          INT                     NOT NULL
        CONSTRAINT PK_svRoles_id PRIMARY KEY,     -- id - нумерация.
    name        VARCHAR(150) DEFAULT '' NOT NULL, -- name - наименование роли.
    description VARCHAR(300) DEFAULT '' NOT NULL, -- description - описание роли.
    created_at  DATETIME     DEFAULT GETDATE(),   -- created_at - дата создания записи.
    created_by  INT                     NOT NULL, -- created_by - табельный номер сотрудника.
    updated_at  DATETIME     DEFAULT GETDATE(),   -- updated_at - дата изменения записи.
    updated_by  INT                     NOT NULL, -- updated_by - табельный номер сотрудника изменивший запись.

    CONSTRAINT UQ_svRoles_name UNIQUE (name)
);

-- INSERT INTO dbo.svRoles
-- VALUES (0, 'user', N'Пользователь доступ к просмотру данных', getdate(), 6680, GETDATE(), 6680);
-- INSERT INTO dbo.svRoles
-- VALUES (1, 'storekeeper', N'Кладовщик доступ к CRUD операциям с продукцией', getdate(), 6680, GETDATE(), 6680);
-- INSERT INTO dbo.svRoles
-- VALUES (2, 'supervisor', N'Руководитель доступ к CRUD операциям с продукцией и редактированию данных', getdate(), 6680,
--         GETDATE(), 6680);
-- INSERT INTO dbo.svRoles
-- VALUES (3, 'administrator', N'Администратор доступ ко всем операциям', getdate(), 6680, GETDATE(), 6680);
-- INSERT INTO dbo.svRoles
-- VALUES (4, 'test', N'Тестовые операции для проверки работоспособности', getdate(), 6680, GETDATE(), 6680);
-- INSERT INTO dbo.svRoles
-- VALUES (5, 'operator', N'Диспетчер заводит сменно-суточное задание', getdate(), 6680, GETDATE(), 6680);

CREATE PROCEDURE dbo.svRolesAll -- ХП получение списка ролей.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, name, description, created_at, created_by, updated_at, updated_by FROM dbo.svRoles;
END
GO;

CREATE PROCEDURE dbo.svRolesAdd -- ХП добавляет роль
    @Id INT, -- ИД роли
    @Name VARCHAR(150), -- наименование роли
    @Description VARCHAR(300), -- описание роли
    @PerformerId INT -- ид сотрудника
AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO dbo.svRoles(id, name, description, created_at, created_by, updated_at, updated_by)
    VALUES (@Id, @Name, @Description, GETDATE(), @PerformerId, GETDATE(), @PerformerId);
END
GO;

CREATE PROCEDURE dbo.svRolesUpdById -- ХП обновляет роль по ид
    @Id INT, -- ИД роли
    @Name VARCHAR(150), -- наименование роли
    @Description VARCHAR(300), -- описание роли
    @PerformerId INT -- ид сотрудника
AS
BEGIN
    SET NOCOUNT ON;

    UPDATE dbo.svRoles
    SET name        = @Name,
        description = @Description,
        updated_by  = @PerformerId
    WHERE id = @Id;
END
GO;

CREATE PROCEDURE dbo.svRolesFindById -- ХП ищет роль
@Id INT -- ид роли
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, name, description, created_at, created_by, updated_at, updated_by
    FROM dbo.svRoles
    WHERE id = @Id;
END
GO;

CREATE PROCEDURE dbo.svRolesExistsById -- ХП проверяет, существует ли роль
@Id INT
AS
BEGIN
    SET NOCOUNT ON;

    DECLARE @Exists BIT = 0

    IF EXISTS(SELECT 1 FROM dbo.svRoles WHERE id = @Id)
        BEGIN
            SET @Exists = 1
        END

    SELECT @Exists AS exists_flag
END
GO;