package store

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

var clapLock sync.Mutex

type ClapRecord struct {
	IP   string `json:"ip"`
	Time int64  `json:"time"`
}
type ClapStat struct {
	Count   int          `json:"count"`
	Records []ClapRecord `json:"records"`
}

var clapFile = "data/clap.json"

func AddClap(ip string) (int, bool, error) {
	clapLock.Lock()
	defer clapLock.Unlock()
	stat := &ClapStat{}
	_ = loadClap(stat)
	for _, r := range stat.Records {
		if r.IP == ip {
			return stat.Count, false, nil // 已点过
		}
	}
	stat.Count++
	stat.Records = append(stat.Records, ClapRecord{IP: ip, Time: time.Now().Unix()})
	saveClap(stat)
	return stat.Count, true, nil
}
func GetClapCount() int {
	clapLock.Lock()
	defer clapLock.Unlock()
	stat := &ClapStat{}
	_ = loadClap(stat)
	return stat.Count
}
func loadClap(stat *ClapStat) error {
	data, err := os.ReadFile(clapFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, stat)
}
func saveClap(stat *ClapStat) error {
	data, _ := json.Marshal(stat)
	return os.WriteFile(clapFile, data, 0644)
}
