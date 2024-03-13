package conductor

import (
	"errors"
	"strings"
)

var ConditionExpressionDivider = "::"
var ConditionExpressionOperatorAny = "any"
var ConditionExpressionOperatorEqual = "equal"
var ConditionExpressionOperatorPrefix = "prefix"

// ConditionExression is a struct that represents a condition expression
// It contains two parts: the operator and the comparison value
// examples:
//   - any::
//   - equal::orders.paid
//   - prefix::orders.
type ConditionExression struct {
	Expression string

	operator string
	value    string
}

func (ce *ConditionExression) Validate() error {
	parts := strings.Split(ce.Expression, ConditionExpressionDivider)
	if len(parts) != 2 {
		return errors.New("CONDUCTOR.CONDITION_EXPRESSION.MALFORMED_EXPRESSION.ERROR")
	}

	ce.operator = parts[0]
	ce.value = parts[1]

	return nil
}

func (ce *ConditionExression) Match(data string) (bool, error) {
	// always validate the expression before matching
	if err := ce.Validate(); err != nil {
		return false, err
	}

	if ce.operator == ConditionExpressionOperatorAny {
		return true, nil
	}

	if ce.operator == ConditionExpressionOperatorEqual {
		return data == ce.value, nil
	}

	if ce.operator == ConditionExpressionOperatorPrefix {
		return strings.HasPrefix(data, ce.value), nil
	}

	return false, errors.New("CONDUCTOR.CONDITION_EXPRESSION.OPERATOR_UNKNOWN")
}
