package controller

import (
	"backEnd/common"
	"backEnd/common/response"
	"backEnd/model"
	"encoding/json"
	"io/ioutil"
	// "backEnd/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	// "strconv"
	// "strings"
)


/* 
前端保存按钮
*/
func SaveCanvas(ctx *gin.Context) {
	fmt.Println("SaveCanvas")
	// Type := ctx.Query("type")
	// Task := ctx.Query("task")
	// pageKind := ctx.Query("pageKind")
	id := ctx.Query("id")
	newCellData := ctx.Query("cell")

	// 读取JSON文件
	jsonData, err := ioutil.ReadFile("./data/data.json")
	if err != nil {
		response.Response(ctx, http.StatusOK, 404, nil, "fail")
		return
	}

    // 解析JSON数据
    var data model.Data
    if err := json.Unmarshal(jsonData, &data); err != nil {
        fmt.Println("解析JSON数据失败:", err)
		response.Response(ctx, http.StatusOK, 404, nil, "fail")
        return
    }

    // 找到要覆盖的design
    var targetDesign *model.Design
    for i := range data.Designs {
        if fmt.Sprintf("%d", data.Designs[i].Id) == id {
            targetDesign = &data.Designs[i]
            break
        }
    }

    if targetDesign != nil {
        // 解析新的cell数据为map[string]interface{}
        var newCellMap map[string]interface{}
        if err := json.Unmarshal([]byte(newCellData), &newCellMap); err != nil {
            fmt.Println("解析新的cell数据失败:", err)
			response.Response(ctx, http.StatusOK, 404, nil, "fail")
            return
        }

        // 将新的cell数据覆盖到目标design的Cells字段
        targetDesign.Cell = newCellMap

        // 将更新后的数据重新写入JSON文件
        updatedData, err := json.MarshalIndent(data, "", "  ")
        if err != nil {
            fmt.Println("序列化更新后的数据失败:", err)
			response.Response(ctx, http.StatusOK, 404, nil, "fail")
            return
        }

        if err := ioutil.WriteFile("data.json", updatedData, os.ModePerm); err != nil {
            fmt.Println("写入JSON文件失败:", err)
			response.Response(ctx, http.StatusOK, 404, nil, "fail")
            return
        }
        fmt.Println("数据已更新")
		response.Success(ctx, nil, "success")
    } else {
        fmt.Println("未找到指定的design")
		response.Response(ctx, http.StatusOK, 404, nil, "fail")
    }
}

/*
前端运行按钮
*/
func RunCanvas(ctx *gin.Context) {
	fmt.Println("SaveCanvas")
	design_name := ctx.Query("design_name")
	id := ctx.Query("id")
	newCellData := ctx.Query("cell")
	Rank := ctx.Query("rank")

	// 读取JSON文件
	jsonData, err := ioutil.ReadFile("./data/data.json")
	if err != nil {
		response.Response(ctx, http.StatusOK, 404, nil, "fail")
		return
	}

    // 解析JSON数据
    var data model.Data
    if err := json.Unmarshal(jsonData, &data); err != nil {
        fmt.Println("解析JSON数据失败:", err)
		response.Response(ctx, http.StatusOK, 404, nil, "fail")
        return
    }

    // 找到要覆盖的design
    var targetDesign *model.Design
    for i := range data.Designs {
        if fmt.Sprintf("%d", data.Designs[i].Id) == id {
            targetDesign = &data.Designs[i]
            break
        }
    }

    if targetDesign != nil {
        // 解析新的cell数据为map[string]interface{}
        var newCellMap map[string]interface{}
        if err := json.Unmarshal([]byte(newCellData), &newCellMap); err != nil {
            fmt.Println("解析新的cell数据失败:", err)
			response.Response(ctx, http.StatusOK, 404, nil, "fail")
            return
        }

        // 将新的cell数据覆盖到目标design的Cells字段
        targetDesign.Cell = newCellMap

        // 将更新后的数据重新写入JSON文件
        updatedData, err := json.MarshalIndent(data, "", "  ")
        if err != nil {
            fmt.Println("序列化更新后的数据失败:", err)
			response.Response(ctx, http.StatusOK, 404, nil, "fail")
            return
        }

        if err := ioutil.WriteFile("data.json", updatedData, os.ModePerm); err != nil {
            fmt.Println("写入JSON文件失败:", err)
			response.Response(ctx, http.StatusOK, 404, nil, "fail")
            return
        }
        fmt.Println("数据已更新")

		// example添加
		db := common.InitDB()

		example := model.Example{
			Example_name:   	design_name,
			Rank:   			Rank,
			State:           	"运行",
			Cpu_num:           	4,
			Gpu_num: 			1,
			Post_data:         	0,
			Dataset_url:      	"",
			Model_name:    		"",
			Model_type:         "",
			Epoch_num:          "200e",
			Loss:           	"loss",
			Optimizer:          "optimizer",
			Decay:           	"decay",
			Evaluation:         "evaluation",
			Model_url:          "model_url",
			Memory:           	"2000M",
			Start_time:         0,
			End_time:           0,
		}
		// 判重处理
		//pageKind、task、type、dataset_name
		db.Where("Example_name = ?", example).First(&example)
		if example.Id != 0 {
			fmt.Println("该实例已存在")
			response.Response(ctx, http.StatusOK, 404, nil, "The example already exists")
		} else {
			// 新增
			db.Create(&example)
			response.Success(ctx, nil, "success")
		}
    } else {
        fmt.Println("未找到指定的design")
		response.Response(ctx, http.StatusOK, 404, nil, "fail")
    }
}

/*
前端保存png
*/
func SaveCanvasPNG(ctx *gin.Context) {
	fmt.Println("SaveCanvasPNG")
	Id := ctx.PostForm("id")
	Type := ctx.Query("type")
	Task := ctx.Query("task")
	image, _ := ctx.FormFile("image")
	// 要创建的文件夹的路径
	folderPath := "./auxTool-frontEnd-main/" + Type + "/" + Task + "/" + Id
	// 使用os.Mkdir创建文件夹
	err := os.Mkdir(folderPath, 0755) // 0755是文件夹的权限设置
	if err != nil {
		fmt.Println("创建文件夹失败:", err)
	}

	dst := folderPath + "/image.png"
	// 上传文件至指定的完整文件路径
	ctx.SaveUploadedFile(image, dst)

	response.Success(ctx, nil, "success")
	
}

/*
前端上传reward.txt
*/
func UploadReward(ctx *gin.Context) {
	fmt.Println("SaveCanvasPNG")
	Id := ctx.PostForm("id")
	Type := ctx.Query("type")
	Task := ctx.Query("task")
	reward, _ := ctx.FormFile("image")
	// 要创建的文件夹的路径
	folderPath := "./auxTool-frontEnd-main/" + Type + "/" + Task + "/" + Id
	// 使用os.Mkdir创建文件夹
	err := os.Mkdir(folderPath, 0755) // 0755是文件夹的权限设置
	if err != nil {
		fmt.Println("创建文件夹失败:", err)
	}

	dst := folderPath + "/reward.txt"
	// 上传文件至指定的完整文件路径
	ctx.SaveUploadedFile(reward, dst)

	response.Success(ctx, nil, "success")
}

/*
前端上传actions.json
*/
func UploadActions(ctx *gin.Context) {
	fmt.Println("SaveCanvasPNG")
	Id := ctx.PostForm("id")
	Type := ctx.Query("type")
	Task := ctx.Query("task")
	actions, _ := ctx.FormFile("image")
	// 要创建的文件夹的路径
	folderPath := "./auxTool-frontEnd-main/" + Type + "/" + Task + "/" + Id
	// 使用os.Mkdir创建文件夹
	err := os.Mkdir(folderPath, 0755) // 0755是文件夹的权限设置
	if err != nil {
		fmt.Println("创建文件夹失败:", err)
	}

	dst := folderPath + "/action.json"
	// 上传文件至指定的完整文件路径
	ctx.SaveUploadedFile(actions, dst)

	response.Success(ctx, nil, "success")
}

/*
前端上传loss.csv
*/
func UploadLoss(ctx *gin.Context) {
	fmt.Println("SaveCanvasPNG")
	Id := ctx.PostForm("id")
	Type := ctx.Query("type")
	Task := ctx.Query("task")
	loss, _ := ctx.FormFile("image")
	// 要创建的文件夹的路径
	folderPath := "./auxTool-frontEnd-main/" + Type + "/" + Task + "/" + Id
	// 使用os.Mkdir创建文件夹
	err := os.Mkdir(folderPath, 0755) // 0755是文件夹的权限设置
	if err != nil {
		fmt.Println("创建文件夹失败:", err)
	}

	dst := folderPath + "/loss.csv"
	// 上传文件至指定的完整文件路径
	ctx.SaveUploadedFile(loss, dst)

	response.Success(ctx, nil, "success")
}