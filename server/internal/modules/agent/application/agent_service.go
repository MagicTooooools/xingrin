package application

type QueryService struct {
	agentStore AgentStore
}

func NewQueryService(agentStore AgentStore) *QueryService {
	return &QueryService{agentStore: agentStore}
}

type CommandService struct {
	agentStore AgentStore
	tokenStore RegistrationTokenStore
}

func NewCommandService(agentStore AgentStore, tokenStore RegistrationTokenStore) *CommandService {
	return &CommandService{agentStore: agentStore, tokenStore: tokenStore}
}
