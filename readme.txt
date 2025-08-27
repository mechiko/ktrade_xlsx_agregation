программа генерирует csv файлы по xlsx файлу с 5 колонками 
и отчет нанесения когда есть база данных

CREATE INDEX "idx_omcsn_code" ON "order_mark_codes_serial_numbers" (
	"code"
);
