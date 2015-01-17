package oval

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	_ = iota
	EVROP_LESS_THAN
	EVROP_EQUALS
	EVROP_UNKNOWN
)

func evr_lookup_operation(s string) int {
	switch s {
	case "less than":
		return EVROP_LESS_THAN
	}
	return EVROP_UNKNOWN
}

func evr_operation_str(val int) string {
	switch val {
	case EVROP_LESS_THAN:
		return "<"
	case EVROP_EQUALS:
		return "="
	default:
		return "?"
	}
}

// Asset an epoch is present within a version string, if not a modified
// string is returned including a default epoch value (0)
func evr_epoch_assert(s string) string {
	f, _ := regexp.MatchString("^\\d+\\:", s)
	if !f {
		return "0:" + s
	}
	return s
}

func evr_extract(s string) (string, string, string) {
	var epoch string
	var version string
	var release string

	s0 := strings.Split(s, ":")
	if len(s0) < 2 {
		panic("evr_extract: can't extract epoch")
	}
	epoch = s0[0]

	// If we have a + character in the vr component, we treat this as a
	// dpkg style package, otherwise rpm
	if strings.Contains(s0[1], "+") {
		s0 = strings.Split(s0[1], "+")
		if len(s0) < 2 {
			panic("evr_extract: + tokenize failure")
		}
		version = s0[0]
		release = s0[1]
	} else {
		version = s0[1]
		release = ""
	}

	debug_prt("[evr_extract] epoch=%v, version=%v, revision=%v\n", epoch, version, release)
	return epoch, version, release
}

func evr_e_compare(actual string, check string) int {
	ai, err := strconv.Atoi(actual)
	if err != nil {
		panic("evr_e_compare: atoi actual")
	}
	ci, err := strconv.Atoi(check)
	if err != nil {
		panic("evr_e_compare: atoi actual")
	}
	if ai > ci {
		return 1
	} else if ai < ci {
		return -1
	}
	return 0
}

//
// Compare a component of a version string containing an integer followed
// by a character
//
func evr_v_compare_numalpha(actual string, check string) (string, string, bool) {
	return "", "", true
}

func evr_v_compare(actual string, check string) int {
	if len(actual) == 0 || len(check) == 0 {
		panic("evr_v_compare: empty version string")
	}
	debug_prt("[evr_v_compare] %v %v\n", actual, check)
	dashbuf_actual := strings.Split(actual, "-")
	dashbuf_check := strings.Split(check, "-")

	ret := 0
	for x, actdash := range dashbuf_actual {
		if x >= len(dashbuf_check) {
			// The actual string has more dash components then the
			// comparison string does, return what we have so far
			// and ignore the rest
			return ret
		}
		// sigma represents the component of the dash buffer from the
		// check value for this cycle
		sigma := dashbuf_check[x]

		dot_act := strings.Split(actdash, ".")
		dot_sig := strings.Split(sigma, ".")

		// Loop through each dot component in the version string;
		// regular integer values are handled simply, if the component
		// has other types of characters we pass them off to extended
		// handling functions
		for y, actdot := range dot_act {
			if y >= len(dot_sig) {
				// There are more version components in this
				// string then in the check version, treat this
				// as greater if we have gotten this far
				return 1
			}
			ai, err_a := strconv.Atoi(actdot)
			ci, err_c := strconv.Atoi(dot_sig[y])

			// If the conversion failed for either one, try a few
			// other extended comparison methods for the component
			extend := true
			if err_a != nil || err_c != nil {
				extend = false
				ai, ci, valid := evr_v_compare_numalpha(actdot, dot_sig[y])
				if valid {
					extend = true
					if ai > ci {
						return 1
					}
				}
			}
			if !extend {
				panic("evr_v_compare: conversion and extended methods failed")
			}

			if ai > ci {
				return 1
			} else if ai < ci {
				return -1
			}
			// Otherwise the components were equal, continue on with the next
			// one
		}
	}
	return ret
}

func evr_r_compare(actual string, check string) int {
	return 0
}

func evr_compare(op int, actual string, check string) bool {
	debug_prt("[evr_compare] %v %v %v\n", actual, evr_operation_str(op), check)

	actual = evr_epoch_assert(actual)
	check = evr_epoch_assert(check)
	a_e, a_v, a_r := evr_extract(actual)
	c_e, c_v, c_r := evr_extract(check)

	res_epoch := evr_e_compare(a_e, c_e)
	res_version := evr_v_compare(a_v, c_v)
	res_release := evr_r_compare(a_r, c_r)
	debug_prt("[evr_compare] [%v:%v:%v] \n", res_epoch, res_version, res_release)

	switch op {
	case EVROP_EQUALS:
		if res_epoch == 0 &&
			res_version == 0 &&
			res_release == 0 {
			return true
		}
		return false
	case EVROP_LESS_THAN:
		switch res_epoch {
		case -1:
			return true
		case 1:
			return false
		}
		switch res_version {
		case -1:
			return true
		case 1:
			return false
		}
		switch res_release {
		case -1:
			return true
		case 1:
			return false
		}
		return false
	default:
		panic("unknown evr comparison operation")
	}
	return false
}

func Test_evr_compare(op int, actual string, check string) bool {
	return evr_compare(op, actual, check)
}
