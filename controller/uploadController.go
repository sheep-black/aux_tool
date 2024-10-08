package controller

import (
	"backEnd/common"
	"backEnd/common/response"
	"backEnd/model"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Record struct {
	Name   string `json:"name"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

// 上传文件
func UploadCsvFile(ctx *gin.Context) {
	// 路径参数

	Task := ctx.PostForm("task")
	Type := ctx.PostForm("type")
	Id := ctx.PostForm("id")
	num, _ := strconv.Atoi(Id)
	time := ctx.PostForm("time")
	file, _ := ctx.FormFile("file")


	// 要创建的文件夹的路径
	folderPath := "./auxTool-frontEnd-main/" + Type + "/" + Task + "/" + Id
	// 使用os.Mkdir创建文件夹
	err := os.Mkdir(folderPath, 0755) // 0755是文件夹的权限设置
	if err != nil {
		fmt.Println("创建文件夹失败:", err)
	}

	dst := folderPath + "/" + file.Filename
	dst_sql := Type + "/" + Task + "/" + Id + "/" + file.Filename
	// 上传文件至指定的完整文件路径
	ctx.SaveUploadedFile(file, dst)

	// 读取csv文件信息，获取行数、列数、数据类型
	numColumns, numRows, Types := GetCSVInfo(dst)
	result := CreateTable(Task, Type, num, numColumns, numRows, Types, file.Filename, dst_sql, time)
	if result == "success" {
		response.Success(ctx, nil, "success")
	} else {
		response.Success(ctx, nil, "fail")
	}
}

// 将参数解析为map（因为可多选）
func DeleteCsvFile(ctx *gin.Context) {
	var records []Record
	db := common.InitDB()
	// 解析JSON数据到结构体切片中
	if err := ctx.BindJSON(&records); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		fmt.Println("解析json失败")
	}

	var flag int
	flag = 1
	for _, record := range records {
		Task := reflect.ValueOf(record).FieldByName("Task").String()
		Dataset_Id := reflect.ValueOf(record).FieldByName("Dataset_Id").String()
		Type := reflect.ValueOf(record).FieldByName("Type").String()
		Csv_name := reflect.ValueOf(record).FieldByName("Csv_name").String()

		// 在这里执行删除记录的操作
		datatable := []model.Datatable{}
		result := db.Where("Task = ? and Dataset_name = ? and Type = ? and Csv_name = ?", Task, Dataset_Id, Type, Csv_name).Delete(&datatable)
		if result.RowsAffected == 0 {
			fmt.Println("删除失败")
			flag = 0
		} else {
			fmt.Println("删除成功")
		}
	}
	// 有一条或者多条记录删除失败
	if flag == 0 {
		response.Response(ctx, http.StatusOK, 404, nil, "fail")
	}
	response.Success(ctx, nil, "success")

}

// 下载文件
func DownloadCsvFile(ctx *gin.Context) {
	db := common.InitDB()
	// 路径参数
	Task := ctx.PostForm("task")
	Dataset_Id := ctx.PostForm("dataset_name")
	Type := ctx.PostForm("type")
	Table_name := ctx.PostForm("table_name")
	datatables := []model.Datatable{}

	db.Where("Task = ? and Type = ? and Dataset_Id = ? and Table_name = ?", Task, Type, Dataset_Id, Table_name).Find(&datatables)
	if len(datatables) == 0 {
		fmt.Println("没找到")
		response.Response(ctx, http.StatusOK, 404, nil, "fail")
	} else {
		filePath := datatables[0].Csv_path
		fmt.Println(filePath)
		response.Response(ctx, http.StatusOK, 404, gin.H{"url": filePath}, filePath)
	}

}
