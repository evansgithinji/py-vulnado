package persistence

import (
	"strings"

	"goapp/internal/domain/entity"
)

// LdapFilterBuilder builds LDAP filter strings.
type LdapFilterBuilder struct{}

// NewLdapFilterBuilder creates a new LdapFilterBuilder.
func NewLdapFilterBuilder() *LdapFilterBuilder {
	return &LdapFilterBuilder{}
}

// BuildAuthFilter builds an LDAP authentication filter.
// VULNERABLE: LDAP Injection (CWE-90) - direct string concatenation
func (b *LdapFilterBuilder) BuildAuthFilter(username, password string) string {
	// VULNERABLE: Direct string concatenation into LDAP filter
	return "(&(uid=" + username + ")(userPassword=" + password + "))"
}

// BuildSearchFilter builds an LDAP search filter.
// VULNERABLE: LDAP Injection (CWE-90) - direct string concatenation
func (b *LdapFilterBuilder) BuildSearchFilter(query string) string {
	// VULNERABLE: Direct string concatenation into LDAP filter
	return "(|(uid=" + query + ")(cn=" + query + ")(mail=" + query + "))"
}

// LdapConnection represents a connection to an LDAP directory.
type LdapConnection struct {
	config *entity.LdapConfig
	users  []entity.LdapUser
}

// Search executes a filter against the in-memory user store.
func (c *LdapConnection) Search(filter string) []entity.LdapUser {
	var results []entity.LdapUser
	for _, user := range c.users {
		if c.matchFilter(user, filter) {
			results = append(results, user)
		}
	}
	return results
}

// matchFilter evaluates a naive LDAP filter against a user.
func (c *LdapConnection) matchFilter(user entity.LdapUser, filter string) bool {
	filter = strings.TrimSpace(filter)
	for strings.HasPrefix(filter, "(") && strings.HasSuffix(filter, ")") {
		filter = filter[1 : len(filter)-1]
		filter = strings.TrimSpace(filter)
	}

	// AND filter: &(cond1)(cond2)
	if strings.HasPrefix(filter, "&") {
		inner := filter[1:]
		conditions := c.splitConditions(inner)
		for _, cond := range conditions {
			if !c.matchFilter(user, cond) {
				return false
			}
		}
		return true
	}

	// OR filter: |(cond1)(cond2)
	if strings.HasPrefix(filter, "|") {
		inner := filter[1:]
		conditions := c.splitConditions(inner)
		for _, cond := range conditions {
			if c.matchFilter(user, cond) {
				return true
			}
		}
		return false
	}

	// Simple attribute=value comparison
	parts := strings.SplitN(filter, "=", 2)
	if len(parts) != 2 {
		// Malformed filter - match all (vulnerable behavior)
		return true
	}

	attr := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])

	// Wildcard matches all
	if val == "*" {
		return true
	}

	userVal := c.getAttr(user, attr)
	return strings.EqualFold(userVal, val)
}

func (c *LdapConnection) splitConditions(s string) []string {
	var conditions []string
	depth := 0
	start := -1
	for i, ch := range s {
		if ch == '(' {
			if depth == 0 {
				start = i
			}
			depth++
		} else if ch == ')' {
			depth--
			if depth == 0 && start >= 0 {
				conditions = append(conditions, s[start:i+1])
				start = -1
			}
		}
	}
	return conditions
}

func (c *LdapConnection) getAttr(user entity.LdapUser, attr string) string {
	switch strings.ToLower(attr) {
	case "uid":
		return user.UID
	case "cn":
		return user.CN
	case "mail":
		return user.Mail
	case "userpassword":
		return user.Password
	case "ou":
		return user.OU
	}
	return ""
}

// LdapConnectionAdapter creates LDAP connections for a given config and user set.
type LdapConnectionAdapter struct{}

// NewLdapConnectionAdapter creates a new LdapConnectionAdapter.
func NewLdapConnectionAdapter() *LdapConnectionAdapter {
	return &LdapConnectionAdapter{}
}

// Connect creates an LdapConnection bound to the given config and users.
func (a *LdapConnectionAdapter) Connect(config *entity.LdapConfig, users []entity.LdapUser) *LdapConnection {
	return &LdapConnection{
		config: config,
		users:  users,
	}
}
