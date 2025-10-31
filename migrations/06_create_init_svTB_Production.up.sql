-- СОЗДАТЬ ТАБЛИЦУ ДЛЯ РЕДАКТИРОВАНИЯ ПРОДУКЦИИ.
CREATE TABLE dbo.svTB_Production
(
    idProduction   INT IDENTITY (1,1)
        CONSTRAINT PK_skTB_Productuion
            PRIMARY KEY NONCLUSTERED,                         -- idProduction - ид продукции.
    PrFolder       INT            DEFAULT 0         NOT NULL, -- PrFolder - пока поле неизвестно.
    PrName         VARCHAR(300)   DEFAULT ''        NOT NULL, -- PrName - наименование варианта упаковки продукции для упаковщика.
    PrShortName    VARCHAR(100)   DEFAULT ''        NOT NULL, -- PrShortName - короткое наименование продукции для этикетки.
    PrPackName     VARCHAR(300)   DEFAULT ''        NOT NULL, -- PrPackName - вариант упаковки.
    PrType         VARCHAR(100)   DEFAULT ''        NOT NULL, -- PrType - декларированная или нет.
    PrArticle      VARCHAR(5)     DEFAULT ''        NOT NULL  -- PrArticle - артикул варианта упаковки.
        CONSTRAINT IX_svTB_Production
            UNIQUE,
    PrColor        VARCHAR(20)    DEFAULT ''        NOT NULL, -- PrColor - цвет продукции.
    PrBarCode      VARCHAR(13)    DEFAULT ''        NOT NULL, -- PrBarCode - бар-код.
    PrCount        INT            DEFAULT 0         NOT NULL, -- PrCount - количество продукции в ряду.
    PrRows         INT            DEFAULT 0         NOT NULL, -- PrRows - количество рядов.
    PrWeight       DECIMAL(19, 3) DEFAULT 0         NOT NULL, -- PrWeight - вес п\п (кг).
    PrHWD          VARCHAR(100)   DEFAULT ''        NOT NULL, -- PrHWD - габариты (мм) 1000(высота)х1200(ширина)х1000(глубина).
    PrInfo         VARCHAR(1024),                             -- PrInfo - информация о продукции\комментарий.
    PrStatus       BIT            DEFAULT 1         NOT NULL, -- PrStatus - статус продукции.
    PrEditDate     DATETIME       DEFAULT GETDATE(),          -- PrEditDate - дата и время изменения записи.
    PrEditUser     INT            DEFAULT 1,                  -- PrEditUser - роль сотрудника. По умолчанию 1 - администратор.
    PrPart         INT            DEFAULT 0         NOT NULL, -- PrPart - номер текущей партии, номер партии и дата указываются вручную и не будут изменяться автоматически с течением времени.
    PrPartLastDate DATETIME       DEFAULT GETDATE() NOT NULL, -- PrPartLastDate - дата выпуска партии.
    PrPartAutoInc  SMALLINT       DEFAULT 1         NOT NULL, -- PrPartAutoInc - нумерация партии и даты! Ручная(0), Автоматическая(1), С указанной даты(2).
    PrPartRealDate DATETIME,                                  -- PrPartRealDate - дата продукции пока неизвестное поле.
    PrArchive      BIT            DEFAULT 0         NOT NULL, -- PrArchive - архивная запись или нет.
    PrPerGodn      SMALLINT       DEFAULT 0,                  -- PrPerGodn - срок годности в месяцах.
    PrSAP          VARCHAR(15),                               -- PrSAP - сап-код.
    PrProdType     BIT            DEFAULT 1         NOT NULL, -- PrProdType - тип продукции пищевая\не пищевая.
    PrUmbrella     BIT            DEFAULT 1         NOT NULL, -- PrUmbrella - беречь от влаги.
    PrSun          BIT            DEFAULT 1         NOT NULL, -- PrSun - беречь от солнца.
    PrDecl         BIT            DEFAULT 0         NOT NULL, -- PrDecl - декларирования или нет.
    PrParty        BIT            DEFAULT 0         NOT NULL, -- PrParty - партионная или нет.
    PrGL           SMALLINT       DEFAULT 0         NOT NULL, -- PrGL - петля Мёбиуса.
    PrVP           SMALLINT       DEFAULT 0         NOT NULL, -- PrVP - ванная печь.
    PrML           SMALLINT       DEFAULT 0         NOT NULL, -- PrML - машинная линия на печи.
    performerid    INT            DEFAULT 0         NOT NULL, -- performerid - изначально 0 это робот делает, теперь робота не будет.
    dtact          DATETIME       DEFAULT GETDATE() NOT NULL  -- dtact - дата актуальная для создания записи.
);