package filter

import (
	"proxymap/pkg/db"
	"strings"
)

type BackList struct {
	WafBanner        []string
	BlackCode        []int
	BlackExt         []string
	BlackContentType []string
	BlackMethod      []string
	Blacklength      int
}

var BList = BackList{
	WafBanner: []string{"您的访问请求可能对网站造成安全威胁", "您的请求包含恶意行为", "已被服务器拒绝拦截时间", "Request blocked",
		"The incident ID is", "非法操作！", "网站防火墙", "已被网站管理员设置拦截", "拦截提示", "入侵防御系统",
		"可能托管恶意活动", "您正在试图非法攻击", "玄武盾", "检测到可疑访问", "当前访问疑似黑客攻击", "创宇盾", "安全拦截",
		"您的访问可能对网站造成威胁", "此次访问可能会对网站造成安全威胁", "have been blocked", "您的访问被阻断",
		"可能对网站造成安全威胁", "您的请求是黑客攻击", "WAF拦截"},
	BlackCode: []int{0, 400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418, 419, 420, 421,
		422, 423, 424, 425, 426, 427, 428, 429, 430, 431, 432, 433, 434, 435, 436, 437, 438, 439, 440, 441, 442, 443, 444, 445, 446,
		447, 448, 449, 450, 451, 452, 453, 454, 455, 456, 457, 458, 459, 460, 461, 462, 463, 464, 465, 466, 467, 468, 469, 470, 471,
		472, 473, 474, 475, 476, 477, 478, 479, 480, 481, 482, 483, 484, 485, 486, 487, 488, 489, 490, 491, 492, 493, 494, 495, 496,
		497, 498, 499},
	BlackExt: []string{".bmp", ".jpg", ".png", ".tif", ".gif", ".pcx", ".tga", ".exif", ".fpx", ".svg", ".psd", ".css",
		".cdr", ".pcd", ".dxf", ".ufo", ".eps", ".ai", ".raw", ".WMF", ".webp", ".avif", ".apng", ".avi", ".rmvb",
		".rm", ".asf", ".divx", ".mpg", ".mpeg", ".mpe", ".wmv", ".mp4", ".mkv", ".vob", ".mov", ".flv", ".wma", ".mp3"},
	BlackContentType: []string{"image/jpeg", "image/png", "image/gif", "text/css", "audio/mpeg", "video/mp4"},
	BlackMethod:      []string{"CONNECT", "OPTIONS"},
	Blacklength:      1,
}

func (bl *BackList) IsPass(sinfo *db.StoreInfo) bool {
	for _, item := range bl.WafBanner {
		if strings.Contains(sinfo.ReqBody, item) {
			return false
		}
	}
	for _, item := range bl.BlackCode {
		if sinfo.StatusCode == item {
			return false
		}
	}
	for _, item := range bl.BlackExt {
		if sinfo.Extension == item {
			return false
		}
	}
	for _, item := range bl.BlackContentType {
		if sinfo.ContentType == item {
			return false
		}
	}
	for _, item := range bl.BlackMethod {
		if sinfo.Method == item {
			return false
		}
	}
	if bl.Blacklength > sinfo.Length {
		return false
	}
	return true
}
