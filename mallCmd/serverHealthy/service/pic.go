package service

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/server/pkg/database/mongo"
	"fmt"
	json "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"serverHealthy/modelHealthy"
	"strconv"
	"strings"
)

func DownLoadPic(name string) ([]string, error) {
	var ids []string
	dirPath := fmt.Sprintf("../clientHealthy/pic/%s", name)
	fileInfo, _ := os.Stat(dirPath)
	if fileInfo.IsDir() {
		list, err := listImageFiles(dirPath)
		if err != nil {
			return ids, app.Err(app.Fail, "读取图片失败")
		}
		for _, v := range list {
			var info modelHealthy.Pic
			pic, err := convertToJSON(name, v)
			if err != nil {
				return ids, err
			}
			if _, err := mongo.Collection(&pic).InsertOne(pic); err == nil {
				err = mongo.Collection(&pic).Fields("_id").Where(bson.M{"file_path": v}).FindOne(&info)
				if err == nil {
					info.FilePath = fmt.Sprintf("../clientHealthy/pic/%s/%s", name, strconv.Itoa(int(info.Id))+filepath.Ext(pic.FilePath))
					if err = os.Rename(pic.FilePath, info.FilePath); err != nil {
						return ids, err
					}
					info.Size = pic.Size
					info.Name = pic.Name
					_, _ = mongo.Collection(&pic).Where(bson.M{"_id": info.Id}).UpdateOne(&info)
				}
			}
			id, _ := json.Marshal(info.Id)
			ids = append(ids, string(id))
		}
	}
	return ids, nil
}

func listImageFiles(root string) ([]string, error) {
	var imageFiles []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查是否为图片文件（JPEG 或 PNG）
		if !info.IsDir() && isPic(path) {
			newPath := changeChar(path)
			if count, _ := mongo.Collection(&modelHealthy.Pic{}).Where(bson.M{"file_path": newPath}).Count(); count == 0 {
				imageFiles = append(imageFiles, newPath)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return imageFiles, nil
}

func changeChar(str string) (newStr string) {
	strBytes := []rune(str)
	for i := 0; i < len(strBytes); i++ {
		if strBytes[i] == '\\' {
			strBytes[i] = '/'
		}
	}
	newStr = string(strBytes)
	return newStr
}

func isPic(filePath string) bool {
	extension := strings.ToLower(filepath.Ext(filePath))
	return extension == ".jpg" || extension == ".jpeg" || extension == ".png"
}

func convertToJSON(name string, imagePath string) (modelHealthy.Pic, error) {
	var pic modelHealthy.Pic
	file, err := os.Open(imagePath)
	if err != nil {
		return pic, app.Err(app.Fail, "打开图片失败")
	}
	defer file.Close()
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("解码图像配置时出错:", err)
		return pic, app.Err(app.Fail, "解码图像配置时出错")
	}
	pic = modelHealthy.Pic{
		Name:     name,
		FilePath: imagePath,
		Size:     image.Point{X: img.Width, Y: img.Height},
	}
	return pic, nil
}
