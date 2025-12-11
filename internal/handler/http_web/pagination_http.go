package http_web

import (
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
	"fmt"
	"math"
)

const numerationDefault = 1

// GeneratePageRange функция для генерации диапазона страниц.
func GeneratePageRange(current, total, maxPages int) []int {
	var pages []int

	if total <= maxPages {
		// 1. Если страниц меньше или равно maxPages, показываем все.
		for i := 1; i <= total; i++ {
			pages = append(pages, i)
		}
	} else {
		// 2. Определяем начальную и конечную страницу.
		start := current - maxPages/2
		end := current + maxPages/2

		if start < 1 {
			start = 1
			end = maxPages
		}

		if end > total {
			end = total
			start = total - maxPages + 1
		}

		for i := start; i <= end; i++ {
			pages = append(pages, i)
		}
	}

	return pages
}

// GetParametersPagination получить параметры пагинации.
func GetParametersPagination(pageStr string, numberPageDefault int) (int, error) {
	if pageStr == "" {
		return numberPageDefault, nil
	}

	page := convert.ConvStrToInt(pageStr)

	if page <= 0 {
		return numberPageDefault, fmt.Errorf("%s: %s", msg.E3300, pageStr)
	}

	return page, nil
}

// CalculateRangeOfElements рассчитать отображаемый диапазон элементов.
func CalculateRangeOfElements(offset int, totalCount int, countOnPage int) (int, int, error) {
	if offset < 0 || totalCount < 0 || countOnPage < 0 {
		return 0, 0, fmt.Errorf("%s", msg.E3301)
	}

	if totalCount == 0 {
		return 0, 0, nil
	}

	if offset >= totalCount {
		return 0, 0, fmt.Errorf("%s: %d > %d", msg.E3303, offset, totalCount)
	}

	startItem := offset + numerationDefault
	endItem := offset + countOnPage

	if startItem > totalCount {
		startItem = 0
	}

	if endItem > totalCount {
		endItem = totalCount
	}

	return startItem, endItem, nil
}

// CalculatePage рассчитать пагинацию.
func CalculatePage(totalCount int, pageSize int, page int) (int, error) {
	if totalCount < 0 {
		return 0, fmt.Errorf("%s: %d", msg.E3301, totalCount)
	}

	if pageSize <= 0 {
		return 0, fmt.Errorf("%s: %d", msg.E3302, pageSize)
	}

	if totalCount == 0 {
		return 1, nil
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	if page > totalPages {
		page = totalPages
	}

	return totalPages, nil
}
