package transformer

import (
	"encoding/json"
	"fmt"
	"github.com/oneliang/util-golang/common"
	"testing"
)

func TestName(t *testing.T) {
	mergeData := &common.MergerConfig{
		MasterDataList: []*map[string]any{
			{
				"A": "1A_VALUE",
				"B": "1B_VALUE",
			},
			{
				"A": "2_VALUE",
				"B": "2_VALUE",
			},
		},
		SlaveDataList: []*common.SlaveData{
			&common.SlaveData{
				DataList: []*map[string]any{
					{
						"A": "1A_VALUE",
						"B": "1B_VALUE",
						"C": "1C_VALUE",
					},
				},
				MergeKeys: []string{"A", "B"},
			},
			&common.SlaveData{
				DataList: []*map[string]any{
					{
						"C": "1C_VALUE",
						"D": "1D_VALUE",
						"E": "1E_VALUE",
					},
				},
				MergeKeys: []string{"C"},
			},
		},
		StaticDataList: nil,
	}
	outputMergedData := common.Merge(mergeData)
	jsonBytes, _ := json.Marshal(outputMergedData)
	fmt.Println(string(jsonBytes))
}
