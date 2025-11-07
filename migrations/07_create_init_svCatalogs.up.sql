-- СОЗДАТЬ ТАБЛИЦУ КАТАЛОГОВ СО СПРАВОЧНИКАМИ.
-- kodcat:
--  0 - конструкторское наименование продукции
--  1 - действие над объектами учета
--  2 - действие над этикеткой
--  3 - цвет продукции
--  4 - принтеры
--  5 - действия для заявок
--  6 - приоритеты
--  7 - статусы заявок
--  8 - ТСД, компьютеры
--  9 - участки упаковки
--  10 - участки хранения
--  11 - объекты учета
--  12 - назначение при списании
--  13 - комментарии к п\п
--  14 - размеры этикеток
--  15 - типы документов
CREATE TABLE dbo.svCatalogs
(
    id          INT IDENTITY (1,1)
        CONSTRAINT PK_svCatalogs
            PRIMARY KEY CLUSTERED ,              -- id - ид
    parid       INT           DEFAULT 0  NOT NULL, -- parid - ид родительской записи. (id).
    kodcat      SMALLINT      DEFAULT 0  NOT NULL, -- kodcat - код справочника.
    kod         SMALLINT      DEFAULT 0  NOT NULL, -- kod - пользовательский код записи (инкремент внутри кода справочника).
    name        VARCHAR(254)  DEFAULT '' NOT NULL, -- name - наименование.
    comm        VARCHAR(1500) DEFAULT '' NOT NULL, -- comm - комментарий.
    dop_int_1   INT           DEFAULT 0  NOT NULL, -- dop_int_1 - дополнительное поле (для kodcat 0, 1, 10, 3, 4, 9).
    dop_int_2   INT           DEFAULT 0  NOT NULL, -- dop_int_2 - дополнительное поле.
    dop_float_1 FLOAT         DEFAULT 0  NOT NULL, -- dop_float_1 - для kodcat=10, возможный процент использования.
    dop_float_2 FLOAT         DEFAULT 0  NOT NULL, -- dop_float_2 - для kodcat=10, вместимость, площадь, объём.
    dop_bit_1   BIT           DEFAULT 0  NOT NULL, -- dop_bit_1 - для kodcat=9, переупаковка 0\1, kodcat=10, площадка (0 -открытая, 1 - закрытая).
    dop_bit_2   BIT           DEFAULT 0  NOT NULL, -- dop_bit_2 - для kodcat=10, наличие ЖД путей. 0\1.
    archive     BIT           DEFAULT 0  NOT NULL, -- archive - архивная запись.
    png         VARBINARY(max) NULL , -- png - хранит картинку складка со строчками.
);
CREATE INDEX idx_svCatalogs_parid ON dbo.svCatalogs(parid);
CREATE INDEX idx_svCatalogs_kodcat ON dbo.svCatalogs(kodcat);
CREATE INDEX idx_svCatalogs_kodcat_kod ON dbo.svCatalogs(kodcat, kod);
CREATE INDEX idx_svCatalogs_kodcat_name ON dbo.svCatalogs(kodcat, name);
CREATE INDEX idx_svCatalogs_archive ON dbo.svCatalogs(archive);