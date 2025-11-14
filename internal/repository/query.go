package repository

const (
	FGWsvPerformerAllQuery        = "exec dbo.svPerformerAll;"                // ХП получение всех сотрудников (только не архивных).
	FGWsvPerformerAuthQuery       = "exec dbo.svPerformerAuth ?, ?;"          // ХП проверяет сотрудника по табельному номеру и паролю для авторизации.
	FGWsvPerformerFindByIdQuery   = "exec dbo.svPerformerFindById ?;"         // ХП ищет информацию о сотруднике по ИД.
	FGWsvPerformerUpdByIdQuery    = "exec dbo.svPerformerUpdById ?, ?, ?, ?;" // ХП обновляет сотрудника по ИД.
	FGWsvPerformerExistsByIdQuery = "exec dbo.svPerformerExistsById ?;"       // ХП проверяет, существует ли сотрудник.
)
