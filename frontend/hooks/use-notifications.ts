/**
 * Notification-related React Query hooks
 */

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { NotificationService } from '@/services/notification.service'
import type {
  GetNotificationsRequest,
} from '@/types/notification.types'

/**
 * Get notification list
 */
export function useNotifications(params?: GetNotificationsRequest) {
  return useQuery({
    queryKey: ['notifications', params],
    queryFn: () => NotificationService.getNotifications(params),
  })
}

/**
 * Get unread notification count
 */
export function useUnreadCount() {
  return useQuery({
    queryKey: ['notifications', 'unread-count'],
    queryFn: () => NotificationService.getUnreadCount(),
    refetchInterval: 30000, // Auto refresh every 30 seconds
  })
}

/**
 * Mark all notifications as read
 */
export function useMarkAllAsRead() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => NotificationService.markAllAsRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
    },
    onError: () => {
    },
  })
}
