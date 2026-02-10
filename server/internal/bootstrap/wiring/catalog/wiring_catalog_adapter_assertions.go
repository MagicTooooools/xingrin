package catalogwiring

import catalogapp "github.com/yyhuni/lunafox/server/internal/modules/catalog/application"

var _ catalogapp.EngineQueryStore = (*catalogEngineStoreAdapter)(nil)
var _ catalogapp.EngineCommandStore = (*catalogEngineStoreAdapter)(nil)
var _ catalogapp.EngineStore = (*catalogEngineStoreAdapter)(nil)

var _ catalogapp.TargetQueryStore = (*catalogTargetStoreAdapter)(nil)
var _ catalogapp.TargetCommandStore = (*catalogTargetStoreAdapter)(nil)
var _ catalogapp.TargetStore = (*catalogTargetStoreAdapter)(nil)

var _ catalogapp.OrganizationStore = (*catalogOrganizationStoreAdapter)(nil)

var _ catalogapp.WordlistQueryStore = (*catalogWordlistStoreAdapter)(nil)
var _ catalogapp.WordlistCommandStore = (*catalogWordlistStoreAdapter)(nil)
var _ catalogapp.WordlistStore = (*catalogWordlistStoreAdapter)(nil)
