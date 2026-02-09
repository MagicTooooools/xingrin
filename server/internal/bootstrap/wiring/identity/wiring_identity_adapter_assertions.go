package identitywiring

import identityapp "github.com/yyhuni/lunafox/server/internal/modules/identity/application"

var _ identityapp.UserStore = (*identityUserStoreAdapter)(nil)
var _ identityapp.AuthUserStore = (*identityUserStoreAdapter)(nil)
var _ identityapp.OrganizationStore = (*identityOrganizationStoreAdapter)(nil)
