package util

import "github.com/mellolo/common/utils/jsonUtil"

type processBar struct {
	percent int
	failed  bool
	done    bool
}

func GetProcessBarJsonString(percent float64) string {
	if percent >= 1 {
		bar := processBar{
			percent: 100,
		}
		return jsonUtil.GetJsonString(bar)
	} else if percent < 0 {
		bar := processBar{
			percent: 0,
		}
		return jsonUtil.GetJsonString(bar)
	}
	bar := processBar{
		percent: int(percent * 100),
	}
	return jsonUtil.GetJsonString(bar)
}

func GetFailedProcessBarJsonString() string {
	bar := processBar{
		failed: true,
	}
	return jsonUtil.GetJsonString(bar)
}

func GetDoneProcessBarJsonString() string {
	bar := processBar{
		percent: 100,
		done:    true,
	}
	return jsonUtil.GetJsonString(bar)
}
