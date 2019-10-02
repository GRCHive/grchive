package render

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func CreateRedirectParams(title string, subtitle string, redirectUrl string) map[string]interface{} {
	newMap := core.StructToMap(*core.LoadTemplateConfig())
	newMap["Title"] = title
	newMap["Subtitle"] = title
	newMap["Redirect"] = redirectUrl
	return newMap
}
