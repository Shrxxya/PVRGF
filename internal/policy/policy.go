package policy

type Policy struct {
	MinLength      int
	MinUppercase   int
	MinLowercase   int
	MinNumbers     int
	MinSpecial     int
	AllowedSpecial string
}

var DomainPolicies = map[string]Policy{
	"gmail.com": {
		MinLength:      12,
		MinUppercase:   1,
		MinLowercase:   1,
		MinNumbers:     1,
		MinSpecial:     1,
		AllowedSpecial: "!@#$%^&*",
	},

	"facebook.com": {
		MinLength:      10,
		MinUppercase:   1,
		MinLowercase:   1,
		MinNumbers:     1,
		MinSpecial:     1,
		AllowedSpecial: "!@#$%^&*",
	},

	"linkedin.com": {
		MinLength:      8,
		MinUppercase:   1,
		MinLowercase:   1,
		MinNumbers:     1,
		MinSpecial:     0,
		AllowedSpecial: "!@#$%^&*",
	},

	"instagram.com": {
		MinLength:      8,
		MinUppercase:   1,
		MinLowercase:   1,
		MinNumbers:     1,
		MinSpecial:     0,
		AllowedSpecial: "!@#$%^&*",
	},
}

func GetPolicy(domain string) (Policy, bool) {
	policy, exists := DomainPolicies[domain]
	return policy, exists
}
