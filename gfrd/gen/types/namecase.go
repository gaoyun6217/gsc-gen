package types

import "regexp"

// NameCase 名称转换工具
type NameCase struct{}

var nc = &NameCase{}

// ToCamel 转为小驼峰 (user_name -> userName)
func (n *NameCase) ToCamel(s string) string {
	if s == "" {
		return ""
	}
	return toCamelCase(s, false)
}

// ToPascal 转为大驼峰 (user_name -> UserName)
func (n *NameCase) ToPascal(s string) string {
	if s == "" {
		return ""
	}
	return toCamelCase(s, true)
}

// ToKebab 转为短横线 (user_name -> user-name)
func (n *NameCase) ToKebab(s string) string {
	if s == "" {
		return ""
	}
	return s
}

// ToSnake 转为下划线 (userName -> user_name)
func (n *NameCase) ToSnake(s string) string {
	if s == "" {
		return ""
	}
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	return re.ReplaceAllString(s, "${1}_${2}")
}

func toCamelCase(s string, pascal bool) string {
	parts := splitBySeparator(s)
	result := ""
	for i, part := range parts {
		if i == 0 && !pascal {
			result += part
		} else {
			if len(part) > 0 {
				result += toTitle(part)
			}
		}
	}
	return result
}

func splitBySeparator(s string) []string {
	// 支持下划线和短横线
	re := regexp.MustCompile("[_-]+")
	parts := re.Split(s, -1)
	result := []string{}
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func toTitle(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = toUpper(runes[0])
	return string(runes)
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 'a' + 'A'
	}
	return r
}

// ToCamel 包级别函数
func ToCamel(s string) string {
	return nc.ToCamel(s)
}

// ToPascal 包级别函数
func ToPascal(s string) string {
	return nc.ToPascal(s)
}

// ToKebab 包级别函数
func ToKebab(s string) string {
	return nc.ToKebab(s)
}

// ToSnake 包级别函数
func ToSnake(s string) string {
	return nc.ToSnake(s)
}
