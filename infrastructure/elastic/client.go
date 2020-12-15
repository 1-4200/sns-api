package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"sns-api/logger"
	"time"
)

var tweetIndex = "sns"
var userIndex = "user"
var urlIndex = "url"
var mediaIndex = "media"

func buildQuery(query *[]map[string]interface{}, value interface{}, attr, column string) {
	var q map[string]interface{}
	switch v := value.(type) {
	case string:
		if len(v) > 0 {
			if attr == "created_at" {
				q = map[string]interface{}{
					"range": map[string]interface{}{
						attr: map[string]interface{}{
							column: v,
						},
					},
				}
			} else {
				q = map[string]interface{}{
					attr: map[string]interface{}{
						column: v,
					},
				}
			}
			*query = append(*query, q)
		}
	case int:
		if v > 0 || (column == "gte" && v >= 0) {
			q = map[string]interface{}{
				"range": map[string]interface{}{
					attr: map[string]interface{}{
						column: v,
					},
				},
			}
			*query = append(*query, q)
		}
	case float64:
		if v > 0 {
			q = map[string]interface{}{
				"range": map[string]interface{}{
					attr: map[string]interface{}{
						column: v,
					},
				},
			}
			*query = append(*query, q)
		}
	case []int:
		if len(v) > 0 {
			q = map[string]interface{}{
				attr: map[string]interface{}{
					column: v,
				},
			}
			*query = append(*query, q)
		}
	case []string:
		if len(v) > 0 && column == "hashtag" && attr == "wildcard" {
			size := len(v)
			for i := 0; i < size; i++ {
				q = map[string]interface{}{
					attr: map[string]interface{}{
						column: map[string]interface{}{
							"value": fmt.Sprintf("%s%s%s", "*", v[i], "*"),
						},
					},
				}
				*query = append(*query, q)
			}
		} else if len(v) > 0 && column == "hashtag" && attr == "match" {
			size := len(v)
			for i := 0; i < size; i++ {
				q = map[string]interface{}{
					attr: map[string]interface{}{
						column: v[i],
					},
				}
				*query = append(*query, q)
			}
		} else if len(v) > 0 {
			q = map[string]interface{}{
				attr: map[string]interface{}{
					column: v,
				},
			}
			*query = append(*query, q)
		}
	}
}

func buildIndexByTimeAdd(index string, t time.Time, num int) []string {
	var m []string
	if num == 0 {
		m = append(m, fmt.Sprintf("%v-%v.%02v", index, t.Year(), int(t.Month())))
	} else {
		for i := 0; i <= num; i++ {
			d := t.AddDate(0, i, 0)
			m = append(m, fmt.Sprintf("%v-%v.%02v", index, d.Year(), int(d.Month())))
		}
	}
	return m
}

func monthDiff(t1, t2 time.Time) int {
	if t2.After(t1) {
		t1, t2 = t2, t1
	}
	return (t1.Year()-t2.Year())*12 + int(t1.Month()) - int(t2.Month())
}

func timeDiff(t1, t2 time.Time) (int, int, int, int, int, int) {
	if t1.Location() != t2.Location() {
		t2 = t2.In(t1.Location())
	}
	if t1.After(t2) {
		t1, t2 = t2, t1
	}
	y1, M1, d1 := t1.Date()
	y2, M2, d2 := t2.Date()

	h1, m1, s1 := t1.Clock()
	h2, m2, s2 := t2.Clock()

	year := int(y2 - y1)
	month := int(M2 - M1)
	day := int(d2 - d1)
	hour := int(h2 - h1)
	min := int(m2 - m1)
	sec := int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return year, month, day, hour, min, sec
}

func encodeQuery(buf *bytes.Buffer, query map[string]interface{}) error {
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return err
	}
	return nil
}

func search(ctx context.Context, l logger.Logging, es *elasticsearch.Client, index string, buf *bytes.Buffer, size int) (map[string]interface{}, error) {
	var r map[string]interface{}

	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(index),
		es.Search.WithBody(buf),
		es.Search.WithSize(size),
		es.Search.WithTrackTotalHits(true),
		//es.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			l.Errorf(fmt.Sprintf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			))
			return nil, err
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	l.Infof(fmt.Sprintf("[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	))

	return r, nil
}

func convertTime(timeStr string) (string, error) {
	layout := "2006-01-02 15:04:05"
	def := "2006-01-01 00:00:00"
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return def, err
	}

	t, err := time.ParseInLocation(layout, timeStr, time.UTC)
	if err != nil {
		return def, err
	}
	return t.In(loc).Format(layout), nil
}
