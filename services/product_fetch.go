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
)

// FetchAndCacheProductsForLocale 抓取多个产品页面中的 productSelectionData，并缓存为 user_config/products_<locale>.json
func FetchAndCacheProductsForLocale(locale string, urls []string) error {
    cleaned := make([]string, 0, len(urls))
    for _, u := range urls {
        u = strings.TrimSpace(u)
        if u != "" {
            cleaned = append(cleaned, u)
        }
    }
    if len(cleaned) == 0 {
        return errors.New("未提供任何产品页面链接")
    }

    var jsonObjects []json.RawMessage
    var preOrderErrors []string
    
    for _, link := range cleaned {
        // 验证URL是否为Apple官网
        if !strings.Contains(link, "apple.com") {
            return fmt.Errorf("请输入Apple官网链接: %s", link)
        }
        
        _, body, errs := gorequest.New().
            Set("user-agent", defaultUA()).
            Set("referer", "https://www.apple.com/shop/buy-iphone").
            Timeout(8 * time.Second).
            Get(link).End()
        if len(errs) > 0 {
            return fmt.Errorf("抓取失败: %s: %v", link, errs[0])
        }
        obj, err := extractProductSelectionData(body)
        if err != nil {
            // 收集预售相关错误
            if strings.Contains(err.Error(), "尚未开放预购") {
                preOrderErrors = append(preOrderErrors, fmt.Sprintf("%s: %s", link, err.Error()))
                continue
            }
            return fmt.Errorf("解析 productSelectionData 失败: %s: %w", link, err)
        }
        jsonObjects = append(jsonObjects, json.RawMessage(obj))
    }
    
    // 如果所有链接都是预售状态
    if len(preOrderErrors) > 0 && len(jsonObjects) == 0 {
        return fmt.Errorf("所有产品均未开放预购:\n%s", strings.Join(preOrderErrors, "\n"))
    }
    
    // 如果没有成功获取任何数据
    if len(jsonObjects) == 0 {
        return errors.New("未能获取任何产品数据")
    }

    // 输出为数组 JSON
    out, err := json.MarshalIndent(jsonObjects, "", "  ")
    if err != nil {
        return err
    }

    if err := os.MkdirAll("user_config", 0755); err != nil {
        return err
    }
    file := filepath.Join("user_config", fmt.Sprintf("products_%s.json", locale))
    if err := os.WriteFile(file, out, 0644); err != nil {
        return err
    }
    return nil
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

