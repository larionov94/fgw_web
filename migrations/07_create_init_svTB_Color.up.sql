-- СОЗДАТЬ ТАБЛИЦУ ЦВЕТОВ ДЛЯ ПРОДУКЦИИ.
CREATE TABLE dbo.svTB_Color
(
    idColor INT IDENTITY (1, 1)
        CONSTRAINT PK_svTB_Color PRIMARY KEY NONCLUSTERED,
    Color   VARCHAR(100) NOT NULL, -- Color - цвет продукции.
    GL      INT          NOT NULL  -- GL - петля Мёбиуса.
);

-- INSERT INTO dbo.svTB_Color (Color, GL) VALUES(N'бесцветное', 70);
-- INSERT INTO dbo.svTB_Color (Color, GL) VALUES(N'голубое', 73);
-- INSERT INTO dbo.svTB_Color (Color, GL) VALUES(N'зеленое', 71);
-- INSERT INTO dbo.svTB_Color (Color, GL) VALUES(N'коричневое', 72);
-- INSERT INTO dbo.svTB_Color (Color, GL) VALUES(N'красное', 78);
-- INSERT INTO dbo.svTB_Color (Color, GL) VALUES(N'оранжевое тарное', 0);
-- INSERT INTO dbo.svTB_Color (Color, GL) VALUES(N'серое', 74);
-- INSERT INTO dbo.svTB_Color (Color, GL) VALUES(N'серое', 75);
-- INSERT INTO dbo.svTB_Color (Color, GL) VALUES(N'черное', 76);
