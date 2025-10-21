-- СОЗДАТЬ ТАБЛИЦУ СПРОВОЧНИК ОПЕРАЦИЙ ДЛЯ AForms.
CREATE TABLE dbo.svTB_Oper
(
    idOper     INT IDENTITY
        CONSTRAINT PK_svTB_Oper PRIMARY KEY NONCLUSTERED,
    OperName   VARCHAR(100) NOT NULL,
    StateName  VARCHAR(100),
    ActionName VARCHAR(100)
);
CREATE INDEX idx_svTB_Oper_idOper ON dbo.svTB_Oper (idOper);
CREATE INDEX idx_svTB_Oper_OperName ON dbo.svTB_Oper (OperName);

INSERT INTO dbo.svTB_Oper (OperName, StateName, ActionName)
VALUES (N'Печать', N'На упаковке', N'Напечатать');
INSERT INTO dbo.svTB_Oper (OperName, StateName, ActionName)
VALUES (N'Упаковка', N'Упакован', N'Упаковать');
INSERT INTO dbo.svTB_Oper (OperName, StateName, ActionName)
VALUES (N'Разупаковка', N'Разупакован', N'Разупаковать');
INSERT INTO dbo.svTB_Oper (OperName, StateName, ActionName)
VALUES (N'Отгрузка', N'Отгружен', N'Отгрузить');
INSERT INTO dbo.svTB_Oper (OperName, StateName, ActionName)
VALUES (N'Оприходование', null, null);

CREATE PROCEDURE dbo.svTB_AllOper -- Получить список операций.
AS
BEGIN
    SET NOCOUNT ON;

    SELECT idOper, OperName, StateName, ActionName FROM dbo.svTB_Oper;
END
GO;

CREATE PROCEDURE dbo.svTB_GetOperByName -- Получить список операций по имени операции.
    @OperName VARCHAR(100)
AS
BEGIN
    SET NOCOUNT ON;

    SELECT idOper, OperName, StateName, ActionName FROM dbo.svTB_Oper WHERE OperName = @OperName;
END
GO;

CREATE PROCEDURE dbo.svTB_GetOperById -- Получить список операций по ИД операции.
    @idOper INT
AS
BEGIN
    SET NOCOUNT ON;

    SELECT idOper, OperName, StateName, ActionName FROM dbo.svTB_Oper WHERE idOper = @idOper;
END
GO;

CREATE PROCEDURE dbo.svTB_AddOper -- Добавить операцию.
    @OperName VARCHAR(100),
    @StateName VARCHAR(100),
    @ActionName VARCHAR(100)
AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO dbo.svTB_Oper (OperName, StateName, ActionName) VALUES (@OperName, @StateName, @ActionName);
END
GO;

CREATE PROCEDURE dbo.svTB_DelOperById -- Удалить операцию.
    @idOper INT
AS
BEGIN
    SET NOCOUNT ON;

    DELETE FROM dbo.svTB_Oper WHERE idOper = @idOper;
END
GO;

CREATE PROCEDURE dbo.svTB_UpdOper -- Обновить операцию.
    @idOper INT,
    @OperName VARCHAR(100),
    @StateName VARCHAR(100),
    @ActionName VARCHAR(100)
AS
BEGIN
    SET NOCOUNT ON;

    UPDATE dbo.svTB_Oper
    SET OperName   = @OperName,
        StateName  = @StateName,
        ActionName = @ActionName
    WHERE idOper = @idOper;
END
GO;