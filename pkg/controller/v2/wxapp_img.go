package v2

import (
	"fmt"
	"os"
	"path/filepath"
	. "rentServer/core/controller"
	config "rentServer/pkg/config"
	"strings"
)

type WxappImgController interface {
	GetImages(c *Context)
}

type wxappImgController struct {
}

func (w wxappImgController) GetImages(c *Context) {

	wxappImg := config.GetConfig().WxappImg
	folderPath := wxappImg.Path
	tree, err := buildImageTree(folderPath, wxappImg.Domain)
	if err != nil {
		fmt.Println("Error building image tree:", err)
		return
	}

	// 这里只是简单打印树形结构，你可以根据需要序列化或处理它
	//fmt.Println(tree)
	c.JSONOK(tree)
}

func NewWxappImgController() WxappImgController {
	return &wxappImgController{}
}

// ImageInfo 包含图片名称和网络路径

// ImageInfo 包含图片名称和网络路径
type ImageInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// buildImageTree 递归构建图片树形结构
func buildImageTree(root string, host string) (map[string]interface{}, error) {
	tree := make(map[string]interface{})

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		// 如果是文件
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".png" || ext == ".svg" || ext == ".jpg" {
				// 获取相对于根目录的文件路径（不包含扩展名）作为键
				relPath, err := filepath.Rel(root, strings.TrimSuffix(path, ext))
				if err != nil {
					return err
				}
				// 替换文件夹分隔符为适合作为键的分隔符（例如：_）
				key := strings.ReplaceAll(relPath, string(filepath.Separator), "_")

				// 构造图片信息
				//imgInfo := ImageInfo{
				//	Name: info.Name(), // 保留原始文件名（包含扩展名）
				//	Path: path,        // 原始文件路径（绝对路径）
				//}

				// 逐级构建map
				parts := strings.Split(key, "_")
				current := tree
				for _, part := range parts[:len(parts)-1] {
					if _, ok := current[part].(map[string]interface{}); !ok {
						current[part] = make(map[string]interface{})
					}
					current = current[part].(map[string]interface{})
				}
				// 叶子节点是图片信息
				current[parts[len(parts)-1]] = host + "/" + path
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tree, nil
}
