package quota

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"wicklight/config"
	"wicklight/logger"
)

const (
	gb int64 = 1024 * 1024 * 1024
	mb int64 = 1024 * 1024
)

var cache map[string]int64
var quota map[string]int64

// InitQuota init quota cache
func InitQuota() {
	cache = make(map[string]int64)
	quota = make(map[string]int64)

	for _, u := range config.Conf.Users {
		if u.Quota != 0 {
			cache[u.Username] = 0
			quota[u.Username] = u.Quota * gb
		}
	}

	logger.Debug("[quota] read usage from file", config.Conf.Log.QuotaFile)
	if config.Conf.Log.QuotaFile != "" {
		fp, err := os.Open(config.Conf.Log.QuotaFile)
		if err == nil {
			bufReader := bufio.NewReader(fp)
			for {
				strs, _, err := bufReader.ReadLine()
				if err != nil {
					break
				}
				raw := strings.Split(string(strs), " ")
				if len(raw) == 2 {
					if usageI, err := strconv.ParseInt(raw[1], 10, 64); err == nil {
						cache[raw[0]] = usageI
					}
				}
			}
		}
		fp.Close()
	}
}

// StoreQuota store quota cache
func StoreQuota() {
	logger.Debug("[quota] store usage to file", config.Conf.Log.QuotaFile)
	if config.Conf.Log.QuotaFile == "" {
		return
	}
	fp, err := os.Create(config.Conf.Log.QuotaFile)
	if err != nil {
		logger.Error("[quota][quote] can not create usage file")
	}
	bufWrtier := bufio.NewWriter(fp)

	for k, v := range cache {
		bufWrtier.WriteString(fmt.Sprintf("%v %v\n", k, v))
	}
	bufWrtier.WriteString("\n")
	bufWrtier.Flush()
	fp.Close()
}

// UpdateQuota update quota
func UpdateQuota(username string, usage int64) {
	if oldUsage, ok := cache[username]; ok {
		cache[username] = oldUsage + usage
	} else {
		cache[username] = usage
	}
}

// CheckQuota is to check quota
func CheckQuota(username string) bool {
	if _, ok := quota[username]; !ok {
		return true
	}
	if _, ok := cache[username]; !ok {
		cache[username] = 0
		return true
	}

	return quota[username] >= cache[username]
}

// PrintQuota print quota and usage
func PrintQuota(username string) string {
	q := "INF"
	u := "0KB"
	if tq, ok := quota[username]; ok {
		q = formatUsage(tq)
	}
	if tu, ok := cache[username]; ok {
		u = formatUsage(tu)
	}
	return u + "/" + q
}

func formatUsage(usage int64) string {
	list := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	idx := 0
	usageF := float64(usage)
	for usageF > 1024.0 && idx <= 5 {
		usageF = usageF / 1024.0
		idx = idx + 1
	}
	return fmt.Sprintf("%.2f %v", usageF, list[idx])
}
