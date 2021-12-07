# File was generated automatically by the snyk-iac-rules tool
# It contains the testing framework for Snyk custom rules with extra functionality for Terraform Plan - modify at your own risk!
package lib.testing

import data.lib

parse_fixture_file(fixture_file) = fixture {
	endswith(fixture_file, "json.tfplan")
	fixture := lib.normalize_to_array(tfplan.unmarshal_file(fixture_file))
}
