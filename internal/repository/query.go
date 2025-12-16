package repository

// СОТРУДНИКИ
const (
	FGWsvPerformerAllQuery         = "exec dbo.svPerformerAll;"                // ХП получение всех сотрудников (только не архивных).
	FGWsvPerformerAuthQuery        = "exec dbo.svPerformerAuth ?, ?;"          // ХП проверяет сотрудника по табельному номеру и паролю для авторизации.
	FGWsvPerformerFindByIdQuery    = "exec dbo.svPerformerFindById ?;"         // ХП ищет информацию о сотруднике по ИД.
	FGWsvPerformerUpdByIdQuery     = "exec dbo.svPerformerUpdById ?, ?, ?, ?;" // ХП обновляет сотрудника по ИД.
	FGWsvPerformerExistsByIdQuery  = "exec dbo.svPerformerExistsById ?;"       // ХП проверяет, существует ли сотрудник.
	FGWsvPerformersCountQuery      = "exec dbo.svPerformersCount;"             // ХП считает общее кол-во сотрудников.
	FGWsvPerformersPaginationQuery = "exec dbo.svPerformersPagination ?, ?;"   // ХП получает сотрудников с нумерации страниц.
	FGWsvPerformerFilterByIdQuery  = "exec dbo.svPerformerFilterById ?;"       // ХП ищет сотрудника по табельному номеру.
)

// РОЛИ
const (
	FGWsvRoleAllQuery        = "exec dbo.svRoleAll;"                // ХП получение списка ролей.
	FGWsvRoleAddQuery        = "exec dbo.svRoleAdd ?, ?, ?, ?;"     // ХП добавляет роль.
	FGWsvRoleFindByIdQuery   = "exec dbo.svRoleFindById ?;"         // ХП ищет роль.
	FGWsvRoleUpdByIdQuery    = "exec dbo.svRoleUpdById ?, ?, ?, ?;" // ХП обновляет роль по ид.
	FGWsvRoleExistsByIdQuery = "exec dbo.svRoleExistsById ?;"       // ХП проверяет, существует ли роль.
	FGWsvRoleDelByIdQuery    = "exec dbo.svRoleDelById ?;"          // ХП проверяет, существует ли роль.
)
