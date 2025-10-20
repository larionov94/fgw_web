-- СОЗДАТЬ ТАБЛИЦУ РОЛЕЙ
CREATE TABLE dbo.svRoles
(
    id          INT          DEFAULT 0  NOT NULL PRIMARY KEY, -- id - id роли.
    name        VARCHAR(150) DEFAULT '' NOT NULL,             -- name - название роли.
    description VARCHAR(250) DEFAULT '' NOT NULL,             -- description - описание роли.
);

CREATE INDEX idx_svRoles_id ON dbo.svRoles (id);

CREATE PROCEDURE dbo.svGetRoles -- Процедура получения всех ролей.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, name, description FROM dbo.svRoles;
END
GO;

CREATE PROCEDURE dbo.svGetRoleById -- Процедура получения роли по id.
    @id INT -- id роли.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, name, description FROM dbo.svRoles WHERE id = @id;
END
GO;

CREATE PROCEDURE dbo.svGetRoleByName -- Процедура получения роли по имени.
    @name VARCHAR(150) -- имя роли.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id, name, description FROM dbo.svRoles WHERE name = @name;
END
GO;

CREATE PROCEDURE dbo.svGetRoleDeleteById -- Процедура удаления роли по id.
    @id INT -- id роли.
AS
BEGIN
    SET NOCOUNT ON;

    DELETE FROM dbo.svRoles WHERE id = @id;
    exec svUpdRolePerformer;
END
GO;

CREATE PROCEDURE dbo.svUpdRolePerformer -- Процедура обновления роли сотрудника.
AS
BEGIN
    SET NOCOUNT ON;

    UPDATE dbo.svPerformers
    SET id_role = 0
    WHERE id_role NOT IN (SELECT id FROM dbo.svRoles); -- id_role = 0, если роли нет в таблице svRoles.
END
GO;

CREATE PROCEDURE dbo.svAddRole -- Процедура добавления роли.
    @id INT, -- id роли.
    @name VARCHAR(150), -- имя роли.
    @description VARCHAR(250) -- описание роли.
AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO dbo.svRoles (id, name, description) VALUES (@id, @name, @description);
END
GO;

INSERT INTO dbo.svRoles
VALUES (0, 'user', N'Пользователь доступ к просмотру данных');
INSERT INTO dbo.svRoles
VALUES (1, 'storekeeper', N'Кладовщик доступ к CRUD операциям с продукцией');
INSERT INTO dbo.svRoles
VALUES (2, 'supervisor', N'Руководитель доступ к CRUD операциям с продукцией и редактированию данных');
INSERT INTO dbo.svRoles
VALUES (3, 'administrator', N'Администратор доступ ко всем операциям');

-- exec dbo.svGetRoles;
-- exec dbo.svGetRoleById 1;
-- exec dbo.svGetRoleByName 'administrator';
-- exec dbo.svGetRoleDeleteById 1;
-- exec dbo.svUpdRolePerformer;