package main

import (
	"agregat/checkdbg"
	"agregat/domain"
	"agregat/htmltmpl"
	"agregat/process"
	"agregat/reductor"
	"agregat/repo"
	"agregat/repo/configdb"
	"agregat/zaplog"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mechiko/dbscan"
	"github.com/mechiko/utility"
	"go.uber.org/zap"
)

var file = flag.String("file", "", "file to parse xlsx")

var Mode = "development"

func errMessageExit(title string, errDescription string) {
	utility.MessageBox(title, errDescription)
	os.Exit(-1)
}

func init() {
	flag.Parse()
}

func main() {
	var logsOutConfig = map[string][]string{
		"logger": {"stdout", "agregat.log"},
	}
	zl, err := zaplog.New(logsOutConfig, true)
	if err != nil {
		errMessageExit("ошибка", err.Error())
	}

	lg, err := zl.GetLogger("logger")
	if err != nil {
		errMessageExit("ошибка", err.Error())
	}
	loger := lg.Sugar()
	loger.Info("agregat started")
	listDbs := make(dbscan.ListDbInfoForScan)
	listDbs[dbscan.Config] = &dbscan.DbInfo{}
	listDbs[dbscan.TrueZnak] = &dbscan.DbInfo{}
	repo, err := repo.New(loger, listDbs, ".")
	if err != nil {
		errMessageExit("ошибка", err.Error())
	}
	// тесты
	if Mode == "development" {
		chk, err := checkdbg.NewChecks(loger, repo)
		if err != nil {
			loger.Errorf("create new check error %v", err)
			errMessageExit("ошибка", err.Error())
		}
		err = chk.Run()
		if err != nil {
			loger.Errorf("check error %v", err)
			errMessageExit("ошибка", err.Error())
		}
	}
	model := &reductor.Model{}
	if err := readModel(loger, repo, model); err != nil {
		errMessageExit("ошибка", err.Error())
	}
	// создаем редуктор с новой моделью
	reductor.New(model, loger)

	fileXLSX := *file
	if fileXLSX == "" {
		fileXLSX, err = utility.DialogOpenFile([]utility.FileType{utility.Excel, utility.All}, "", "in")
		if err != nil {
			errMessageExit("ошибка", err.Error())
		}
	}

	proc, err := process.New(fileXLSX, repo)
	if err != nil {
		errMessageExit("запуск обработки с ошибкой ", err.Error())
	}

	err = proc.ReadXlsx(fileXLSX)
	if err != nil {
		errMessageExit("ошибка чтения файла", err.Error())
	}

	err = proc.ScanPalet()
	if err != nil {
		// errMessageExit("ошибка ScanRecords", err.Error())
		loger.Errorf("ошибка ScanRecords %v", err)
	}

	// записываем короба и палеты раньше пусть будут
	err = proc.Save("out")
	if err != nil {
		// errMessageExit("ошибка ScanRecords", err.Error())
		loger.Errorf("ошибка ScanRecords %v", err)
	}

	// поиск заказов по маркам для отчетов нанесения
	err = proc.ScanRecords()
	if err != nil {
		// errMessageExit("ошибка ScanRecords", err.Error())
		loger.Errorf("ошибка ScanRecords %v", err)
	}
	// запись отчетов нанесения в БД
	err = proc.WriteUtilisation()
	if err != nil {
		errMessageExit("ошибка ScanRecords", err.Error())
		// loger.Errorf("ошибка ScanRecords %v", err)
	}

	htmString, err := htmltmpl.NewTemplate().StringHTML(proc)
	if err != nil {
		errMessageExit("ошибка ScanRecords", err.Error())
	}
	fileHtml := "report_" + utility.String(10)
	fileHtml = filepath.Join("out", fileHtml) + ".html"
	file, err := os.Create(fileHtml)
	if err != nil {
		errMessageExit("Error creating file", err.Error())
	}
	defer file.Close() // Ensure the file is closed when the function exits

	// Write the string content to the file
	_, err = file.WriteString(string(htmString))
	if err != nil {
		errMessageExit("Error writing to file", err.Error())
	}
	utility.OpenFileInShell(fileHtml)

}

func readModel(loger *zap.SugaredLogger, repo domain.Repo, model *reductor.Model) (err error) {
	info := repo.Info(dbscan.Config)
	if info == nil {
		return fmt.Errorf("базы конфиг не найдено")
	}
	db, err := configdb.New(loger, info, dbscan.Config)
	if err != nil {
		return fmt.Errorf("error open config.db")
	}
	defer func() {
		if errclose := db.Close(); errclose != nil {
			if err != nil {
				err = fmt.Errorf("%v %w", errclose, err)
			} else {
				err = fmt.Errorf("%v", errclose)
			}
		}
	}()
	model.Inn, err = db.Key("inn")
	if err != nil {
		return fmt.Errorf("error open config.db")
	}
	model.Kpp, err = db.Key("kpp")
	if err != nil {
		return fmt.Errorf("error open config.db")
	}
	return nil
}
