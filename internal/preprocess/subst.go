package preprocess

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/danieljhkim/hive-duck/internal/config"
)

var re = regexp.MustCompile(`\$\{(hiveconf|hivevar|env):([A-Za-z0-9_.\-]+)\}`)

func Substitute(sql string, cfg *config.Config) (string, error) {
	var missing []string

	out := re.ReplaceAllStringFunc(sql, func(m string) string {
		sub := re.FindStringSubmatch(m)
		kind, key := sub[1], sub[2]

		var val string
		var ok bool
		switch kind {
		case "hiveconf":
			val, ok = cfg.HiveConf[key]
		case "hivevar":
			val, ok = cfg.HiveVar[key]
		case "env":
			val, ok = cfg.Env[key]
		}

		if !ok {
			missing = append(missing, m)
			return m
		}
		// Conservative: substitute raw; callers should quote in SQL when needed.
		return val
	})

	if cfg.StrictVars && len(missing) > 0 {
		return "", fmt.Errorf("missing variables: %s", strings.Join(missing, ", "))
	}
	return out, nil
}
