package utils

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Add global scope flags to a URL request
func AddListGlobalScope(cmd *cobra.Command, url string) string {
	global, _ := cmd.Flags().GetBool("global")

	if global {
		return url + fmt.Sprintf("?globalScope=true")
	}

	return url
}

func AddRequestGlobalScope(cmd *cobra.Command, url string) string {
	global, _ := cmd.Flags().GetBool("global")

	if global {
		return url + fmt.Sprintf("?metadata.requestMeta.globalScope=true")
	}

	return url
}
