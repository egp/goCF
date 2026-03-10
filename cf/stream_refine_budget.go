// stream_refine_budget.go v1
package cf

import "fmt"

func consumeRefineBudget(
	engineName string,
	refinesThisDigit *int,
	refinesTotal *int,
	maxRefinesPerDigit int,
	maxTotalRefines int,
) error {
	*refinesThisDigit = *refinesThisDigit + 1
	*refinesTotal = *refinesTotal + 1

	if maxRefinesPerDigit >= 0 && *refinesThisDigit > maxRefinesPerDigit {
		return fmt.Errorf("%s exceeded MaxRefinesPerDigit=%d", engineName, maxRefinesPerDigit)
	}
	if maxTotalRefines >= 0 && *refinesTotal > maxTotalRefines {
		return fmt.Errorf("%s exceeded MaxTotalRefines=%d", engineName, maxTotalRefines)
	}
	return nil
}

// stream_refine_budget.go v1
