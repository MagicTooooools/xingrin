/**
 * Agent management hooks
 */

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { agentService } from '@/services/agent.service'
import { useToastMessages } from '@/lib/toast-helpers'
import { getErrorCode, getErrorResponseData } from '@/lib/response-parser'
import type { UpdateAgentConfigRequest } from '@/types/agent.types'

export const agentKeys = {
  all: ['agents'] as const,
  lists: () => [...agentKeys.all, 'list'] as const,
  list: (page: number, pageSize: number, status?: string) =>
    [...agentKeys.lists(), { page, pageSize, status }] as const,
  details: () => [...agentKeys.all, 'detail'] as const,
  detail: (id: number) => [...agentKeys.details(), id] as const,
}

export function useAgents(page = 1, pageSize = 10, status?: string) {
  return useQuery({
    queryKey: agentKeys.list(page, pageSize, status),
    queryFn: () => agentService.getAgents(page, pageSize, status),
    refetchInterval: 15000,
  })
}

export function useAgent(id: number) {
  return useQuery({
    queryKey: agentKeys.detail(id),
    queryFn: () => agentService.getAgent(id),
    enabled: id > 0,
  })
}

export function useCreateRegistrationToken() {
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: () => agentService.createRegistrationToken(),
    onSuccess: () => {
      toastMessages.success('toast.agent.token.success', {}, 'agent-token')
    },
    onError: (error: unknown) => {
      toastMessages.errorFromCode(getErrorCode(getErrorResponseData(error)), 'toast.agent.token.error', 'agent-token')
    },
  })
}

export function useUpdateAgentConfig() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateAgentConfigRequest }) =>
      agentService.updateAgentConfig(id, data),
    onSuccess: (_: unknown, { id }: { id: number; data: UpdateAgentConfigRequest }) => {
      queryClient.invalidateQueries({ queryKey: agentKeys.lists() })
      queryClient.invalidateQueries({ queryKey: agentKeys.detail(id) })
      toastMessages.success('toast.agent.config.success')
    },
    onError: (error: unknown) => {
      toastMessages.errorFromCode(getErrorCode(getErrorResponseData(error)), 'toast.agent.config.error')
    },
  })
}


export function useDeleteAgent() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: (id: number) => agentService.deleteAgent(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: agentKeys.lists(), refetchType: 'active' })
      toastMessages.success('toast.agent.delete.success')
    },
    onError: (error: unknown) => {
      toastMessages.errorFromCode(getErrorCode(getErrorResponseData(error)), 'toast.agent.delete.error')
    },
  })
}
