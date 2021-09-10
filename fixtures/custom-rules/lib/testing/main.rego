package lib.testing

import data.lib

assert_response_set(result_set, test_case) {
	total_violations := {res |
		result := result_set[index]
		result.publicId == test_case.publicId
		trace(sprintf("[%s][%d] Issue msg : %s", [test_case.publicId, test_case.index, result.msg]))
		res := index
	}

	trace(sprintf("[%s][%s] Number of issues identified: want %d, got %d", [test_case.publicId, test_case.fixture, count(test_case.want_msgs), count(total_violations)]))
	count(total_violations) == count(test_case.want_msgs)

	violation_match := {res |
		result := total_violations[index]
		result.msg == test_case.want_msgs[_]
		trace(sprintf("[%s][%d] Violation msg : %s", [test_case.publicId, test_case.index, result.msg]))
		res := index
	}

	trace(sprintf("[%s][%s] Number of issues with correct `msg` value: want %d, got %d", [test_case.publicId, test_case.fixture, count(test_case.want_msgs), count(violation_match)]))
	count(violation_match) == count(test_case.want_msgs)
	trace(sprintf("[%s] Fixture %d passed", [test_case.publicId, test_case.index]))
} else = false {
	true
}

parse_fixture_file(fixture_file) = fixture {
	fixture := lib.normalize_to_array(yaml.unmarshal_file(fixture_file))
}

get_result_set(fixture) = result_set {
	result_set := data.rules.deny with input as fixture
}

evaluate_test_cases(publicId, test_cases) {
	passed_tests := {res |
		tc := lib.merge_objects(test_cases[index], {"publicId": publicId, "index": index})
		result_set := get_result_set(tc.fixture)
		assert_response_set(result_set, tc)
		res := index
	}

	trace(sprintf("[%s] Number of test cases passed: want %d, got %d", [publicId, count(test_cases), count(passed_tests)]))
	count(passed_tests) == count(test_cases)
} else = false {
	true
}
