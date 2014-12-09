package mser

import "testing"

var mserTestData = []struct {
	size   int
	params Params
	input  *ExtremalRegion
	output []*ExtremalRegion
}{{
	size:   1,
	params: Params{},
	input:  newER(fromFields{level: 1, area: 1}),
	output: []*ExtremalRegion{newER(fromFields{level: 1, area: 1})},
}, {
	// test MinArea & MaxArea
	size:   20,
	params: Params{MinArea: 0.1, MaxArea: 0.2},
	input: newER(fromFields{level: 10, area: 20},
		newER(fromFields{level: 1, area: 1}),
		newER(fromFields{level: 2, area: 2}),
		newER(fromFields{level: 4, area: 4}),
		newER(fromFields{level: 5, area: 5})),
	output: []*ExtremalRegion{
		newER(fromFields{level: 2, area: 2}),
		newER(fromFields{level: 4, area: 4})},
}, {
	// test Delta
	size:   10,
	params: Params{Delta: 1},
	input: newER(fromFields{level: 10, area: 10}, // 0
		newER(fromFields{level: 9, area: 9}, // 1/9 = 0.11...
			newER(fromFields{level: 5, area: 5}, // 0
				newER(fromFields{level: 4, area: 4}, // 1/4 = 0.25
					newER(fromFields{level: 2, area: 2}))))), // 0
	output: []*ExtremalRegion{newER(fromFields{level: 10, area: 10},
		newER(fromFields{level: 5, area: 5},
			newER(fromFields{level: 2, area: 2})))},
}, {
	// test Delta
	size:   10,
	params: Params{Delta: 2},
	input: newER(fromFields{level: 10, area: 10}, // 0
		newER(fromFields{level: 9, area: 9}, // 1/9 = 0.11..
			newER(fromFields{level: 5, area: 5}, // 0
				newER(fromFields{level: 4, area: 4}, // 1/4 = 0.25
					newER(fromFields{level: 2, area: 2}))))), // 2/2 = 1
	output: []*ExtremalRegion{newER(fromFields{level: 10, area: 10},
		newER(fromFields{level: 5, area: 5}))},
}, {
	// test MaxVariation
	size:   10,
	params: Params{Delta: 2, MaxArea: 0.9, MaxVariation: 0.20},
	input: newER(fromFields{level: 10, area: 10}, // 0
		newER(fromFields{level: 8, area: 8}, // 2/8 = 0.25
			newER(fromFields{level: 7, area: 7}, // 1/7 = 0.14...
				newER(fromFields{level: 5, area: 5}, // 2/5 = 0.4
					newER(fromFields{level: 4, area: 4}))))), // 1/4 = 0.25
	output: []*ExtremalRegion{newER(fromFields{level: 7, area: 7})},
}, {
	// test MinDiversity
	size:   10,
	params: Params{Delta: 2, MinDiversity: 0.5},
	input: newER(fromFields{level: 10, area: 10}, // 0
		newER(fromFields{level: 8, area: 8}, // 2/8 = 0.25
			newER(fromFields{level: 7, area: 7}, // 1/7 = 0.14...
				newER(fromFields{level: 5, area: 5}, // 2/5 = 0.4
					newER(fromFields{level: 4, area: 4}))))), // 1/4 = 0.25
	output: []*ExtremalRegion{newER(fromFields{level: 10, area: 10},
		newER(fromFields{level: 4, area: 4}))},
}}

func TestSelectMSER(t *testing.T) {
	for i, tc := range mserTestData {
		forest := selectMSER(tc.input, tc.size, tc.params)

		expected := newER(fromFields{}, tc.output...)
		actual := newER(fromFields{}, forest...)
		if !assertERsEqual(t, expected, actual) {
			expStr, actStr := erStringer{region: expected, indent: 2}, erStringer{region: actual, indent: 2}
			t.Errorf(`test case %d:
input:
  size = %v
  params = %+v
expected:
%v
actual:
%v`, i+1, tc.size, tc.params, expStr, actStr)
		}
	}
}
