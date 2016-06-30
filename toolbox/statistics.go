package toolbox

// Statistics struct
import (
	"fmt"
	"sync"
	"time"
)

// short string format
func toS(d time.Duration) string {

	u := uint64(d)
	if u < uint64(time.Second) {
		switch {
		case u == 0:
			return "0"
		case u < uint64(time.Microsecond):
			return fmt.Sprintf("%.2fns", float64(u))
		case u < uint64(time.Millisecond):
			return fmt.Sprintf("%.2fus", float64(u)/1000)
		default:
			return fmt.Sprintf("%.2fms", float64(u)/1000/1000)
		}
	} else {
		switch {
		case u < uint64(time.Minute):
			return fmt.Sprintf("%.2fs", float64(u)/1000/1000/1000)
		case u < uint64(time.Hour):
			return fmt.Sprintf("%.2fm", float64(u)/1000/1000/1000/60)
		default:
			return fmt.Sprintf("%.2fh", float64(u)/1000/1000/1000/60/60)
		}
	}

}

type Statistics struct {
	Actor      string
	Method     string
	RequestNum int64
	MinTime    time.Duration
	MaxTime    time.Duration
	TotalTime  time.Duration
}

// URLMap contains several statistics struct to log different data
type URLMap struct {
	lock        sync.RWMutex
	LengthLimit int //limit the urlmap's length if it's equal to 0 there's no limit
	urlmap      map[string]map[string]*Statistics
}

// AddStatistics add statistics task.
// it needs request method, request url, request controller and statistics time duration
func (m *URLMap) AddStatistics(actorName, methodName string, requesttime time.Duration) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if actor, ok := m.urlmap[actorName]; ok {
		if s, ok := actor[methodName]; ok {
			s.RequestNum++
			if s.MaxTime < requesttime {
				s.MaxTime = requesttime
			}
			if s.MinTime > requesttime {
				s.MinTime = requesttime
			}
			s.TotalTime += requesttime
		} else {
			nb := &Statistics{
				Actor:      actorName,
				Method:     methodName,
				RequestNum: 1,
				MinTime:    requesttime,
				MaxTime:    requesttime,
				TotalTime:  requesttime,
			}
			m.urlmap[actorName][methodName] = nb
		}

	} else {
		if m.LengthLimit > 0 && m.LengthLimit <= len(m.urlmap) {
			return
		}
		methodmap := make(map[string]*Statistics)
		nb := &Statistics{
			Actor:      actorName,
			Method:     methodName,
			RequestNum: 1,
			MinTime:    requesttime,
			MaxTime:    requesttime,
			TotalTime:  requesttime,
		}
		methodmap[methodName] = nb
		m.urlmap[actorName] = methodmap
	}
}

// GetMap put url statistics result in io.Writer
func (m *URLMap) GetMap() map[string]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var fields = []string{"actor", "method", "times", "used", "max used", "min used", "avg used"}

	var resultLists [][]string
	content := make(map[string]interface{})
	content["Fields"] = fields

	for k, v := range m.urlmap {
		for kk, vv := range v {
			result := []string{
				fmt.Sprintf("% -16s", k),
				fmt.Sprintf("% -16s", kk),
				fmt.Sprintf("% -16d", vv.RequestNum),
				fmt.Sprintf("% -16s", toS(vv.TotalTime)),
				fmt.Sprintf("% -16s", toS(vv.MaxTime)),
				fmt.Sprintf("% -16s", toS(vv.MinTime)),
				fmt.Sprintf("% -16s", toS(time.Duration(int64(vv.TotalTime)/vv.RequestNum))),
			}
			resultLists = append(resultLists, result)
		}
	}
	content["Data"] = resultLists
	return content
}

func (m *URLMap) PrintMap() {
	content := m.GetMap()
	// fmt.Println(content)

	fields := content["Fields"].([]string)
	fmt.Println(
		[]string{
			fmt.Sprintf("% -16s", fields[0]),
			fmt.Sprintf("% -16s", fields[1]),
			fmt.Sprintf("% -16s", fields[2]),
			fmt.Sprintf("% -16s", fields[3]),
			fmt.Sprintf("% -16s", fields[4]),
			fmt.Sprintf("% -16s", fields[5]),
			fmt.Sprintf("% -16s", fields[6]),
		},
	)

	resultLists := content["Data"].([][]string)
	for _, v := range resultLists {
		fmt.Println(v)
	}
}

// GetMapData return all mapdata
func (m *URLMap) GetMapData() []map[string]interface{} {

	var resultLists []map[string]interface{}

	for k, v := range m.urlmap {
		for kk, vv := range v {
			result := map[string]interface{}{
				"actor":      k,
				"method":     kk,
				"times":      vv.RequestNum,
				"total_time": toS(vv.TotalTime),
				"max_time":   toS(vv.MaxTime),
				"min_time":   toS(vv.MinTime),
				"avg_time":   toS(time.Duration(int64(vv.TotalTime) / vv.RequestNum)),
			}
			resultLists = append(resultLists, result)
		}
	}
	return resultLists
}

// StatisticsMap hosld global statistics data map
var StatisticsMap *URLMap

func init() {
	StatisticsMap = &URLMap{
		urlmap: make(map[string]map[string]*Statistics),
	}
}
