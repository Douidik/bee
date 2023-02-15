package main
import "strings"

func UnescapeStr(s string) string {
	rp := strings.NewReplacer(
		`\a`, "\a",
		`\b`, "\b",
		`\f`, "\f",
		`\n`, "\n",
		`\r`, "\r",
		`\t`, "\t",
		`\v`, "\v",
		`\'`, "'",
		`\\`, "\\",
	)
	
	return rp.Replace(s)
}
