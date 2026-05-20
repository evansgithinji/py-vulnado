package entity

// LdapConfig holds LDAP server configuration for a domain.
type LdapConfig struct {
	Domain        string
	BaseDN        string
	Host          string
	Port          int
	BindDN        string
	AdminPassword string
}

// LdapAuthRequest is a DTO for LDAP authentication requests.
type LdapAuthRequest struct {
	Username string
	Password string
	Domain   string
}

// LdapUser represents a user entry in the LDAP directory.
type LdapUser struct {
	UID      string `json:"uid"`
	CN       string `json:"cn"`
	Mail     string `json:"mail"`
	Password string `json:"userPassword"`
	OU       string `json:"ou"`
}
