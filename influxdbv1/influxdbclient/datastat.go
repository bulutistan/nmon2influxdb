package influxdbclient

import "sort"

type DataStat struct {
	Name   string
	Min    float64
	Max    float64
	Mean   float64
	Median float64
	Length int
}

type DataStats []DataStat

func (ds *DataStats) FieldSort(field string) {
	switch field {
	case "name":
		sort.Sort(NameDataStats{*ds})
	case "min":
		sort.Sort(MinDataStats{*ds})
	case "max":
		sort.Sort(MaxDataStats{*ds})
	case "median":
		sort.Sort(MedianDataStats{*ds})
	default:
		sort.Sort(MeanDataStats{*ds})
	}
}

func (slice DataStats) Len() int {
	return len(slice)
}

func (slice DataStats) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type NameDataStats struct{ DataStats }

func (slice NameDataStats) Less(i, j int) bool {
	return slice.DataStats[i].Name > slice.DataStats[j].Name
}

type MinDataStats struct{ DataStats }

func (slice MinDataStats) Less(i, j int) bool {
	return slice.DataStats[i].Min > slice.DataStats[j].Min
}

type MaxDataStats struct{ DataStats }

func (slice MaxDataStats) Less(i, j int) bool {
	return slice.DataStats[i].Max > slice.DataStats[j].Max
}

type MeanDataStats struct{ DataStats }

func (slice MeanDataStats) Less(i, j int) bool {
	return slice.DataStats[i].Mean > slice.DataStats[j].Mean
}

type MedianDataStats struct{ DataStats }

func (slice MedianDataStats) Less(i, j int) bool {
	return slice.DataStats[i].Median > slice.DataStats[j].Median
}

func Sum(data []float64) (sum float64) {
	for _, n := range data {
		sum += n
	}
	return sum
}

func Mean(data []float64) (mean float64) {
	sum := Sum(data)
	return sum / float64(len(data))
}

func BuildStats(dsets []*DataSet) (stats DataStats) {
	for _, ds := range dsets {
		for name, data := range ds.Datas {
			length := len(data)

			//sorting data
			sort.Float64s(data)
			var stat DataStat
			if len(ds.Tags) > 0 {
				for _, tagValue := range ds.Tags {
					if len(stat.Name) > 0 {
						stat.Name = stat.Name + "_"
					}
					stat.Name = stat.Name + tagValue
				}
			} else {
				stat.Name = name
			}
			stat.Min = data[0]
			stat.Max = data[length-1]
			stat.Mean = Mean(data)
			stat.Length = length
			if length%2 == 0 {
				stat.Median = Mean(data[length/2-1 : length/2+1])
			} else {
				stat.Median = float64(data[length/2])
			}
			stats = append(stats, stat)
		}
	}
	return
}
