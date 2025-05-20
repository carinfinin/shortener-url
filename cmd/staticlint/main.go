// Staticlint делает провекрку кода.
// проверяет код на ошибки проверка безопастности, оптимизация и стиль, статический онализ.
// так же имеет дополнительную проверку в пакете main от прямых вызовов os.Exit()
//
// Пример запуска
// go run main.go ./... # для проверки всех пакетов
package main

import (
	"github.com/carinfinin/shortener-url/internal/app/errchekers"
	"github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"strings"
)

func main() {

	checks := []*analysis.Analyzer{
		printf.Analyzer,             // проверяет соответствие форматных строк и аргументов в функциях типа Printf
		structtag.Analyzer,          // Пакет structtag определяет анализатор, который проверяет правильность формирования тегов полей структуры.
		assign.Analyzer,             //  обнаруживает бесполезные присваивания
		shadow.Analyzer,             //  находит затененные переменные
		analyzer.Analyzer,           // go critic оптимизация
		errchekers.ErrCheckAnalyzer, // проверка на вызов os.Exit() в пакете main
		bodyclose.Analyzer,          // bodyclose проверка закрыт ли body
	}

	checksStatic := map[string]bool{
		"ST1013": true, // проверка использования констант для кодов ошибок HTTP, а не магические числа
		"ST1012": true, // проверка емени для ошибок
	}
	// стандартные проверки staticcheck
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") || checksStatic[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}

	multichecker.Main(checks...)
}
