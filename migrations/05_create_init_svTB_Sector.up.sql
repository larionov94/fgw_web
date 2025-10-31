-- СОЗДАТЬ ТАБЛИЦУ ПЕЧЕК.
CREATE TABLE dbo.svTB_Sector
(
    idSector       INT IDENTITY (1,1)
        CONSTRAINT PK_svTB_PrSector PRIMARY KEY NONCLUSTERED,
    SectorName     VARCHAR(150)                   NOT NULL, -- SectorName - наименование печки.
    SectorEditDate DATETIME     DEFAULT GETDATE() NOT NULL, -- SectorEditDate - дата редактирования печки.
    SectorEditUser INT          DEFAULT 0         NOT NULL, -- SectorEditUser - право на редактирование печки.
    SecVPML        VARCHAR(150) DEFAULT ''        NOT NULL, -- SecVPML - список линий печки.
    performerid    INT          DEFAULT 0         NOT NULL, -- performerid - ид сотрудника (возможно роли не играет).
    dtact          DATETIME     DEFAULT GETDATE() NOT NULL, -- dtact - дата создания печки.
    TicketSize     VARCHAR(10)                              -- TicketSize - размер этикетки.
)

CREATE INDEX idx_svTB_Sector_SectorName ON dbo.svTB_Sector (SectorName);

-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'ВП №3 (Лин:31,32)', '20130330 19:32:10.000', 1, N'31,32', 0, '20210608 14:52:39.000', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'Без сектора (Тестовый)', '20220506 10:36:01', 1, N'0', 0, '20210608 14:52:39', N'10x20');
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'ВП №7 (ВП-№7)', '20130330 19:32:22.790', 1, N'71,72,73', 0, '20210608 14:52:39.323', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'ВП №5 (ВП-№5)', '20130718 16:13:11.267', 1, N'51,52,53,54', 0, N'20210608 14:52:39.323', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'ВП №61 (Лин:61)', '20150408 16:17:08.270', 1, N'61', 0, '20210608 14:52:39.323', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'ВП №1 (ВП-№1)', N'20130725 09:26:52.417', 1, N'11,12,13,14', 0, N'20210608 14:52:39.323', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'ВП №62 (Лин:62)', N'20150408 16:17:12.177', 1, N'62', 0, N'20210608 14:52:39.323', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'ВП №33 (Лин:33)', N'20230313 10:04:52.537', 1, N'33', 0, N'20210608 14:52:39.323', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'Переупаковка (УПУ)', N'20220825 09:06:25.350', 1, N'', 0, N'20210813 12:57:27.257', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'Участок декорирования стекла', N'20250122 12:10:56.900', 0, N'', 0, N'20250120 16:09:02.830', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'Без сектора (Тестовый 2)', N'20230616 09:35:10.987', 1, N'', 0, N'20230616 08:35:59.233', null);
-- INSERT INTO dbo.svTB_Sector (SectorName, SectorEditDate, SectorEditUser, SecVPML, performerid, dtact, TicketSize) VALUES (N'ВП №4 (ВП-№4)', N'20230616 14:45:23.080', 1, N'41', 0, N'20230616 14:45:23.080', null);
