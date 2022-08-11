package file

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/itnotebooks/zip"
)

// From:https://github.com/itnotebooks
// ZipDo 递归压缩，默认采用AES256加密方式加密
// 支持以下加密方式
// Standard         ZIP标准，安全性最低
// AES128           AES128位，安全性高
// AES192           AES192位，安全性高
// AES256           AES256位，安全性最高，本程序默认采用此加密方式
func ZipDo(dst, src string, encrypt bool, password, enc string) (err error) {
	var dstFileBaseName = ""
	// 创建压缩文件对象
	zfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zfile.Close()

	// 通过文件对象生成写入对象
	zFileWriter := zip.NewWriter(zfile)
	defer func() {
		// 检测一下是否成功关闭
		if err := zFileWriter.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(src)

		if !strings.HasSuffix(src, "/") {
			dstName := filepath.Base(dst)
			dstFileBaseName = strings.TrimSuffix(dstName, filepath.Ext(dstName))
		}
	}

	// 将文件写入 zFileWriter 对象 ，可能会有很多个目录及文件，递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		if strings.HasSuffix(path, ".zip") {
			return
		}
		//创建文件头
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		if baseDir != "" {
			// 如果原目录是以"/"结尾，表示打包指定目录时不包含该目录
			if strings.HasSuffix(src, "/") {
				header.Name = strings.TrimPrefix(path, src)
			} else {
				header.Name = filepath.Join(dstFileBaseName, baseDir, strings.TrimPrefix(path, src))
			}
		}

		if fi.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		var fh io.Writer
		if encrypt {
			// 加密方式
			var encryption = zip.AES256Encryption
			switch enc {
			case "Standard":
				encryption = zip.StandardEncryption
			case "AES128":
				encryption = zip.AES128Encryption
			case "AES192":
				encryption = zip.AES192Encryption
			}
			// 写入文件头信息，并配置密码
			fh, err = zFileWriter.Encrypt(header, password, encryption)
		} else {
			// 写入文件头信息
			fh, err = zFileWriter.CreateHeader(header)
		}

		if err != nil {
			return err
		}

		// 判断是否是标准文件
		if !header.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// 将文件对象拷贝到 writer 结构中
		ret, err := io.Copy(fh, file)
		if err != nil {
			return err
		}
		log.Printf("added： %s, total: %d\n", path, ret)
		return nil
	})
}
