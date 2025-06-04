package files

import (
	"douyin-backend/app/global/my_errors"
	"douyin-backend/app/global/variable"
	"mime/multipart"
	"net/http"
	"os"
)

// 返回值说明：
//	7z、exe、doc 类型会返回 application/octet-stream  未知的文件类型
//	jpg	=>	image/jpeg
//	png	=>	image/png
//	ico	=>	image/x-icon
//	bmp	=>	image/bmp
//  xlsx、docx 、zip	=>	application/zip
//  tar.gz	=>	application/x-gzip
//  txt、json、log等文本文件	=>	text/plain; charset=utf-8   备注：就算txt是gbk、ansi编码，也会识别为utf-8

// 通过文件名获取文件mime信息
func GetFilesMimeByFileName(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		variable.ZapLog.Error(my_errors.ErrorsFilesUploadOpenFail + err.Error())
	}
	defer f.Close()

	// 只需要前 32 个字节就可以了
	buffer := make([]byte, 32)
	if _, err := f.Read(buffer); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsFilesUploadReadFail + err.Error())
		return ""
	}

	return http.DetectContentType(buffer)
}

// GetFilesMimeByFp 通过文件指针获取文件mime信息
func GetFilesMimeByFp(fp multipart.File) string {
	// 读取更多字节以提高MIME类型检测的准确性
	buffer := make([]byte, 512)
	if _, err := fp.Read(buffer); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsFilesUploadReadFail + err.Error())
		return ""
	}

	// 重置文件指针位置，以便后续操作
	if _, err := fp.Seek(0, 0); err != nil {
		variable.ZapLog.Error("Failed to reset file pointer: " + err.Error())
		return ""
	}

	mimeType := http.DetectContentType(buffer)

	// 如果是未知类型，尝试通过文件扩展名判断
	if mimeType == "application/octet-stream" {
		// 获取文件扩展名
		header := make([]byte, 512)
		if _, err := fp.Read(header); err != nil {
			variable.ZapLog.Error(my_errors.ErrorsFilesUploadReadFail + err.Error())
			return mimeType
		}

		// 重置文件指针位置
		if _, err := fp.Seek(0, 0); err != nil {
			variable.ZapLog.Error("Failed to reset file pointer: " + err.Error())
			return mimeType
		}

		// 检查文件头是否包含 MP4 文件标识
		if len(header) >= 8 && header[4] == 'f' && header[5] == 't' && header[6] == 'y' && header[7] == 'p' {
			return "video/mp4"
		}

		// 检查文件头是否包含 QuickTime 文件标识
		if len(header) >= 8 && header[4] == 'm' && header[5] == 'o' && header[6] == 'o' && header[7] == 'v' {
			return "video/quicktime"
		}
	}

	return mimeType
}
