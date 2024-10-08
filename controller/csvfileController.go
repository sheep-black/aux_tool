package controller

import (
	"backEnd/common"
	"backEnd/common/response"
	"backEnd/model"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
	"github.com/gin-gonic/gin"
)

type Header struct {
	XMLName xml.Name `xml:"header"`
	Name    string `xml:"name"`
	Type    string `xml:"type"`
	Precise int    `xml:"precise"`
	Remark  string `xml:"remark"`
}

type Table struct {
	XMLName   xml.Name `xml:"table"`
	Tablename string   `xml:"tablename,attr"`
	Heds      []Header `xml:"header"`
}

type Database struct {
	XMLName      xml.Name `xml:"database"`
	Databasename string   `xml:"databasename,attr"`
	Tabs         []Table  `xml:"table"`
}

// 获取csv文件属性
func GetCSVInfo(csv_path string) (int, int, []string) {
	// 打开CSV文件
	file, error := os.Open(csv_path)
	if error != nil {
		fmt.Println("无法打开CSV文件:", error)
	}
	defer file.Close()

	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 读取CSV行
	lines, error := reader.ReadAll()
	if error != nil {
		fmt.Println("读取CSV文件失败:", error)
	}

	// 获取列数和行数
	numColumns := len(lines[0])
	numRows := len(lines)

	// 创建一个切片用于保存每列的数据类型
	columnTypes := make([]string, numColumns)

	// 遍历第二行以识别每列的数据类型
	for columnIndex := 0; columnIndex < numColumns; columnIndex++ {
		// 默认为字符串类型
		columnType := "字符串"

		// 尝试将第二行数据识别为rune
		if utf8.RuneCountInString(lines[1][columnIndex]) == 1 {
			columnType = "char"
		}

		// 尝试将第二行数据转换为int
		_, error := strconv.Atoi(lines[1][columnIndex])
		if error == nil {
			columnType = "int"
		}

		// 尝试将第二行数据转换为float64
		_, error = strconv.ParseFloat(lines[1][columnIndex], 64)
		if error == nil {
			columnType = "float"
		}

		// 将数据类型保存到切片中
		columnTypes[columnIndex] = columnType
	}

	// 创建一个映射来保存不重复的字符串元素
	uniqueStrings := make(map[string]bool)

	// 将元素添加到映射中
	for _, element := range columnTypes {
		uniqueStrings[element] = true
	}

	var Types []string
	for uniqueElement := range uniqueStrings {
		Types = append(Types, uniqueElement)
	}

	return numColumns, numRows, Types
}

// 获取某数据集的csv表
func GetTable(ctx *gin.Context) {
	// 参数
	Task := ctx.Query("Task")
	Dataset_Id := ctx.Query("Dataset_Id")
	Type := ctx.Query("Type")

	db := common.InitDB()

	datatables := []model.Datatable{}
	db.Where("Task = ? and Type = ? and Dataset_Id = ?", Task, Type, Dataset_Id).Find(&datatables)
	if len(datatables) == 0 {

		response.Response(ctx, http.StatusOK, 404, nil, "No corresponding card found")
	} else {

		response.Success(ctx, gin.H{"datatables": datatables}, "success")
	}
}

// 获取某数据集的csv表
func DeleteTable(ctx *gin.Context) {
	// 参数
	db := common.InitDB()
	Id := ctx.Query("id")
	CsvName := ctx.Query("csv")
	Task := ctx.Query("task")
	Type := ctx.Query("type")
	DatasetName := ctx.Query("datasetName")
	IdList := strings.Split(Id, "/")
	CsvList := strings.Split(CsvName, "/")

	//此循环为删除文件的操作
	for _, value := range CsvList {
		if value == "" {
			fmt.Println("文件删除结束")
		} else {
			erro := os.Remove("./auxTool-frontEnd-main/" + Type + "/" + Task + "/" + DatasetName + "/" + value)
			if erro != nil {
				fmt.Println("delete fail")
			}
		}
	}
	//此循环为删除数据库记录的操作
	for _, value := range IdList {
		if value == "" {
			fmt.Println("记录删除结束")
		} else {

			datatables := []model.Datatable{}
			result := db.Where("Id = ?", value).Delete(&datatables)

			if result.RowsAffected == 0 {
				fmt.Println("删除失败")
				response.Response(ctx, http.StatusOK, 404, nil, "fail")
			} else {
				fmt.Println("删除成功")
				response.Success(ctx, nil, "success")
			}
		}

	}
}

// 添加csv表
func CreateTable(Task string, Type string, Dataset_Id int, numColumns int, numRows int, Types []string, csv_name string, csv_path string, time string) string {
	db := common.InitDB()
	types := ""
	for _, columnType := range Types {
		types = types + columnType
	}
	datatable := model.Datatable{
		Type:       Type,
		Task:       Task,
		Dataset_Id: uint(Dataset_Id),
		Table_name: csv_name,
		Header_num: uint(numColumns),
		Data_len:   uint(numRows),
		Data_type:  types,
		Csv_path:   csv_path,
		Time:       time,
	}
	// 判重处理
	db.Where("Task = ? and Type = ? and Dataset_Id = ? and Table_name = ?", Task, Type, Dataset_Id, csv_name).First(&datatable)
	if datatable.Id != 0 {
		fmt.Println("该卡片已存在")
		// response.Response(ctx, http.StatusOK, 404, nil, "The card already exists")
		return "fail"
	} else {
		// 新增卡片
		db.Create(&datatable)
		fmt.Println("创建表")
		// response.Success(ctx, nil, "success")
		return "success"
	}
}

// 获取表信息

func GetCsvData(ctx *gin.Context) {
	Dataset_id := ctx.Query("id")

	db := common.InitDB()

	datatables := []model.Datatable{}
	db.Where("Dataset_Id = ?", Dataset_id).Find(&datatables)
	if len(datatables) == 0 {

		response.Response(ctx, http.StatusOK, 404, nil, "No corresponding card found")
	} else {
		fmt.Println(datatables)
		response.Success(ctx, gin.H{"datatables": datatables}, "success")
	}

}
func OutPutXml(ctx *gin.Context) {
	CsvPath := ctx.Query("path")
	DataName := ctx.Query("data_name")
	CsvPathList := strings.Split(CsvPath, ",")
	XmlData := Database{Databasename: DataName}
	for _, value := range CsvPathList {
		dst := "./auxTool-frontEnd-main/" + value
		csvPathSingle := strings.Split(value, "/")
		Tab := Table{Tablename: csvPathSingle[len(csvPathSingle)-1]}
		file, error := os.Open(dst)
		if error != nil {
			fmt.Println("无法打开CSV文件:", error)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		lines, error := reader.ReadAll()
		if error != nil {
			fmt.Println("读取CSV文件失败:", error)
		}
		columns := lines[0]
		fmt.Println("columns", columns)

		for index, column := range columns {
			// 默认为字符串类型
			var columnType string = "string"
			var precise int = 0
			// 尝试将第二行数据识别为rune
			if utf8.RuneCountInString(lines[1][index]) == 1 {
				columnType = "string"
			}

			// 尝试将第二行数据转换为int
			_, error := strconv.Atoi(lines[1][index])
			if error == nil {
				columnType = "int"
			}

			// 尝试将第二行数据转换为float64
			_, error = strconv.ParseFloat(lines[1][index], 64)
			if error == nil {
				columnType = "float"
				precise = len(strings.Split(lines[1][index], ".")[1])
			}
			Hed := Header{Name: column, Type: columnType, Precise: precise, Remark: "损失函数"}
			Tab.Heds = append(Tab.Heds, Hed)
		}
		XmlData.Tabs = append(XmlData.Tabs, Tab)
	}

	b, _ := xml.MarshalIndent(XmlData, "", "	")
	b = append([]byte(xml.Header), b...)

	xmlPath := "./auxTool-frontEnd-main/xml/" + DataName + ".xml"

	err := ioutil.WriteFile(xmlPath, b, 0666)
	if err != nil {
		fmt.Println("后端xml文件写入失败:", err)
	}
	ctx.File(xmlPath)
}
