package rule

import (
	"github.com/apex/log"
)

// A Rule func returns if entry passes the rule.
type Rule func(*log.Entry) bool

type ruleHandler struct {
	next  log.Handler
	rules []Rule
}

// Implements log.Handler interface.
func (rh *ruleHandler) HandleLog(entry *log.Entry) error {
	for _, rule := range rh.rules {
		if !rule(entry) {
			return nil
		}
	}

	return rh.next.HandleLog(entry)
}

// Constructs a RuleHandler which passes entries to next Handler
// only if all Rules are satisfied.
func New(next log.Handler, rules ...Rule) *ruleHandler {
	return &ruleHandler{next, rules}
}

// Helper function that merges a set of Rules, returning a Rule which is
// satisfied if at least one of underlying Rules is satisfied.
func Or(rules ...Rule) Rule {
	return func(entry *log.Entry) bool {
		for _, rule := range rules {
			if rule(entry) {
				return true
			}
		}

		return false
	}
}
