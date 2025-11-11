package repository

const (
	FGWsvPerformerAllQuery        = "exec dbo.svPerformerAll;"                                                        // ХП получение всех сотрудников (только не архивных).
	FGWsvPerformerAuthQuery       = "exec dbo.svPerformerAuth @id, @pass;"                                            // ХП проверяет сотрудника по табельному номеру и паролю для авторизации.
	FGWsvPerformerFindByIdQuery   = "exec dbo.svPerformerFindById @id;"                                               // ХП ищет информацию о сотруднике по ИД.
	FGWsvPerformerUpdByIdQuery    = "exec dbo.svPerformerUpdById @id, @id_role_a_forms, @id_role_a_fgw, @updated_by;" // ХП обновляет сотрудника по ИД.
	FGWsvPerformerExistsByIdQuery = "exec dbo.svPerformerExistsById @id;"                                             // ХП проверяет, существует ли сотрудник.
)
