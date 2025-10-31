-- СОЗАТЬ ТАБЛИЦУ СОТРУДНИКОВ ДЛЯ AForms. Данные получает от Галактики через скрипт. Смотреть svet-02.
CREATE TABLE dbo.svTB_Performer
(
    idPerformer INT IDENTITY (10, 1) PRIMARY KEY NONCLUSTERED, -- idPerformer - ИД.
    extSector   INT,                                           -- extSector - внешний ключ на таблицу svTB_Sector.
    PfFolder    INT         DEFAULT (-1)      NOT NULL,        -- PfFolder - Папка. (под вопросом) sysTB_Folder.
    PfName      VARCHAR(100)                  not null,        -- PfName - Имя (ФИО).
    PfBarcode   VARCHAR(20) DEFAULT ''        NOT NULL,        -- PfBarcode - Штрих-код.
    PfEditDate  DATETIME    DEFAULT GETDATE() NOT NULL,        -- PfEditDate - Дата редактирования.
    PfEditUser  INT         DEFAULT (-1)      NOT NULL,        -- PfEditUser - ИД пользователя редактирования. (1 администратор)
    PfTabnum    INT         DEFAULT (0)
        CONSTRAINT IX_svTB_Performer UNIQUE   NOT NULL,        -- PfTabnum - Табельный номер.
    PfPassword  VARCHAR(255)                                   -- PfPassword - пароль сотрудника.
);

CREATE PROCEDURE dbo.svTB_UpdPerformerBySector -- Процедура для обновления привязки печки к сотруднику.
    @IdPerformer INT, -- ид сотрудника
    @IdSector INT -- ид печки
AS
BEGIN
    SET NOCOUNT ON;

    UPDATE dbo.svTB_Performer

    SET extSector = @IdSector
    WHERE idPerformer = @IdPerformer
END
GO;