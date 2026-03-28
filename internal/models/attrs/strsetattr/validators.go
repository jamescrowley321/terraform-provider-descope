package strsetattr

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
)

var CommaSeparatedValidator = setvalidator.ValueStringsAre(
	stringattr.NonEmptyValidator,
	stringvalidator.RegexMatches(regexp.MustCompile(`^[^,]*$`), "must not contain commas"),
)
