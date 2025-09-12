package services

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/parnurzeal/gorequest"

    // Use model.Areas so supported regions come from one source of truth
    "apple-store-helper/model"
)

// isSupportedLocale 检查 locale 是否在内置地区中
func isSupportedLocale(locale string) bool {
    for _, area := range model.Areas {
        if area.Locale == locale {
            return true
        }
    }
    return false
}

// getSupportedLocalesString 返回“支持的地区”字符串，显示 Title 而非 Locale
func getSupportedLocalesString() string {
    titles := make([]string, 0, len(model.Areas))
    for _, area := range model.Areas {
        titles = append(titles, area.Title)
    }
    return strings.Join(titles, ", ")
}

// FetchAndCacheProductsForLocale 抓取多个产品页面中的 productSelectionData，并缓存为 user_config/products_<locale>.json
// 同时抓取对应的门店信息
// 返回实际保存的locale
func FetchAndCacheProductsForLocale(locale string, urls []string) (string, error) {
    cleaned := make([]string, 0, len(urls))
    for _, u := range urls {
        u = strings.TrimSpace(u)
        if u != "" {
            cleaned = append(cleaned, u)
        }
    }
    if len(cleaned) == 0 {
        return "", errors.New("未提供任何产品页面链接")
    }

    // 从URL中提取实际的locale
    actualLocale := ""
    if len(cleaned) > 0 {
        actualLocale = extractLocaleFromURL(cleaned[0])
        if actualLocale == "" {
            return "", errors.New("无法从URL中提取地区信息")
        }
        // 验证是否为支持的地区
        if !isSupportedLocale(actualLocale) {
            return "", fmt.Errorf("不支持当前输入的地区\n仅支持以下地区：%s", getSupportedLocalesString())
        }
    }
    
    var jsonObjects []json.RawMessage
    var preOrderErrors []string
    
    for _, link := range cleaned {
        // 验证URL是否为Apple官网
        if !strings.Contains(link, "apple.com") {
            return "", fmt.Errorf("请输入Apple官网链接: %s", link)
        }
        
        _, body, errs := gorequest.New().
            Set("user-agent", defaultUA()).
            Set("referer", "https://www.apple.com/shop/buy-iphone").
            Timeout(8 * time.Second).
            Get(link).End()
        if len(errs) > 0 {
            return "", fmt.Errorf("抓取失败: %s: %v", link, errs[0])
        }
        obj, err := extractProductSelectionData(body)
        if err != nil {
            // 收集预售相关错误
            if strings.Contains(err.Error(), "尚未开放预购") {
                preOrderErrors = append(preOrderErrors, fmt.Sprintf("%s: %s", link, err.Error()))
                continue
            }
            return "", fmt.Errorf("解析 productSelectionData 失败: %s: %w", link, err)
        }
        jsonObjects = append(jsonObjects, json.RawMessage(obj))
    }
    
    // 如果所有链接都是预售状态
    if len(preOrderErrors) > 0 && len(jsonObjects) == 0 {
        return "", fmt.Errorf("所有产品均未开放预购:\n%s", strings.Join(preOrderErrors, "\n"))
    }
    
    // 如果没有成功获取任何数据
    if len(jsonObjects) == 0 {
        return "", errors.New("未能获取任何产品数据")
    }

    // 输出为数组 JSON
    out, err := json.MarshalIndent(jsonObjects, "", "  ")
    if err != nil {
        return "", err
    }

    if err := os.MkdirAll("user_config", 0755); err != nil {
        return "", err
    }
    // 使用从URL提取的locale而不是传入的locale
    file := filepath.Join("user_config", fmt.Sprintf("products_%s.json", actualLocale))
    if err := os.WriteFile(file, out, 0644); err != nil {
        return "", err
    }
    
    // 注意：门店信息较复杂，建议使用内置数据
    // 用户可以选择内置地区使用对应的门店信息
    
    return actualLocale, nil
}

// extractProductSelectionData 在 HTML 中定位 "productSelectionData"，并提取其后第一个完整的 JSON 对象
func extractProductSelectionData(html string) ([]byte, error) {
    idx := strings.Index(html, "productSelectionData")
    if idx < 0 {
        // 检测是否为预售等待页面
        if strings.Contains(html, "来得早") || strings.Contains(html, "来得好") || 
           strings.Contains(html, "即将") || strings.Contains(html, "北京时间晚") {
            return nil, errors.New("产品尚未开放预购，请等待官方开放预购时间后再试")
        }
        return nil, errors.New("未找到产品数据，可能页面格式已变更或产品尚未发布")
    }
    // 定位到第一个 '{'
    braceRel := strings.Index(html[idx:], "{")
    if braceRel < 0 {
        return nil, errors.New("no opening brace")
    }
    start := idx + braceRel
    depth := 0
    inStr := false
    esc := false
    for i := start; i < len(html); i++ {
        c := html[i]
        if inStr {
            if esc {
                esc = false
            } else {
                if c == '\\' {
                    esc = true
                } else if c == '"' {
                    inStr = false
                }
            }
            continue
        }
        switch c {
        case '"':
            inStr = true
        case '{':
            depth++
        case '}':
            depth--
            if depth == 0 {
                end := i + 1
                return []byte(html[start:end]), nil
            }
        }
    }
    return nil, errors.New("unterminated JSON")
}

func defaultUA() string {
    return "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0 Safari/537.36"
}

// extractLocaleFromURL 从Apple URL中提取locale
// 例如: https://www.apple.com/de/shop/buy-iphone/iphone-17 -> de_DE
// 例如: https://www.apple.com.cn/shop/buy-iphone/iphone-17 -> zh_CN
func extractLocaleFromURL(url string) string {
    // 常见的Apple域名到locale的映射
    localeMap := map[string]string{
        "apple.com.cn": "zh_CN",
        "apple.com.hk": "zh_HK",
        "apple.com.tw": "zh_TW",
        "apple.com.sg": "en_SG",
        "apple.com.my": "en_MY",
        "apple.com.au": "en_AU",
        "apple.co.jp": "ja_JP",
        "apple.co.uk": "en_GB",
        "apple.ca": "en_CA",
        "apple.fr": "fr_FR",
        "apple.de": "de_DE",
        "apple.it": "it_IT",
        "apple.es": "es_ES",
        "apple.nl": "nl_NL",
        "apple.se": "sv_SE",
        "apple.no": "no_NO",
        "apple.dk": "da_DK",
        "apple.fi": "fi_FI",
        "apple.pt": "pt_PT",
        "apple.be": "nl_BE",
        "apple.ch": "de_CH",
        "apple.at": "de_AT",
        "apple.ie": "en_IE",
        "apple.pl": "pl_PL",
        "apple.cz": "cs_CZ",
        "apple.hu": "hu_HU",
        "apple.ro": "ro_RO",
        "apple.bg": "bg_BG",
        "apple.hr": "hr_HR",
        "apple.gr": "el_GR",
        "apple.ru": "ru_RU",
        "apple.tr": "tr_TR",
        "apple.ae": "en_AE",
        "apple.sa": "ar_SA",
        "apple.co.th": "th_TH",
        "apple.co.kr": "ko_KR",
        "apple.co.nz": "en_NZ",
        "apple.co.za": "en_ZA",
        "apple.co.in": "en_IN",
        "apple.com.br": "pt_BR",
        "apple.com.mx": "es_MX",
        "apple.cl": "es_CL",
        "apple.com.co": "es_CO",
        "apple.com.ar": "es_AR",
    }
    
    // 处理 apple.com/xx/ 格式的URL
    if strings.Contains(url, "apple.com/") && !strings.Contains(url, "apple.com.") && !strings.Contains(url, "apple.co.") {
        // 提取路径中的国家代码
        parts := strings.Split(url, "/")
        for i, part := range parts {
            // 检查是否包含apple.com（可能带www前缀）
            if strings.Contains(part, "apple.com") && i+1 < len(parts) {
                countryCode := parts[i+1]
                if len(countryCode) == 2 {
                    // 简单映射，可能需要扩展
                    countryToLocale := map[string]string{
                        "cn": "zh_CN",
                        "hk": "zh_HK",
                        "tw": "zh_TW",
                        "jp": "ja_JP",
                        "kr": "ko_KR",
                        "de": "de_DE",
                        "fr": "fr_FR",
                        "it": "it_IT",
                        "es": "es_ES",
                        "nl": "nl_NL",
                        "se": "sv_SE",
                        "no": "no_NO",
                        "dk": "da_DK",
                        "fi": "fi_FI",
                        "pt": "pt_PT",
                        "be": "nl_BE",
                        "ch": "de_CH",
                        "at": "de_AT",
                        "ie": "en_IE",
                        "uk": "en_GB",
                        "ca": "en_CA",
                        "au": "en_AU",
                        "nz": "en_NZ",
                        "sg": "en_SG",
                        "my": "en_MY",
                        "th": "th_TH",
                        "in": "en_IN",
                        "za": "en_ZA",
                        "ae": "en_AE",
                        "sa": "ar_SA",
                        "tr": "tr_TR",
                        "ru": "ru_RU",
                        "pl": "pl_PL",
                        "cz": "cs_CZ",
                        "hu": "hu_HU",
                        "ro": "ro_RO",
                        "bg": "bg_BG",
                        "hr": "hr_HR",
                        "gr": "el_GR",
                        "br": "pt_BR",
                        "mx": "es_MX",
                        "cl": "es_CL",
                        "co": "es_CO",
                        "ar": "es_AR",
                    }
                    if locale, ok := countryToLocale[countryCode]; ok {
                        return locale
                    }
                }
            }
        }
        // 默认为美国
        return "en_US"
    }
    
    // 从域名中提取
    for domain, locale := range localeMap {
        if strings.Contains(url, domain) {
            return locale
        }
    }
    
    // 默认返回美国
    return "en_US"
}
