package lib.testing

import data.lib

parse_fixture_file(fixture_file) = fixture {
	endswith(fixture_file, "json.tfplan")
	fixture := lib.normalize_to_array(tfplan.unmarshal_file(fixture_file))
}
