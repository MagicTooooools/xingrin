/**
 * WebSocket 实时通知 Hook
 */

import { useCallback, useEffect, useState, useRef } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import type { BackendNotification, Notification, BackendNotificationLevel, NotificationSeverity } from '@/types/notification.types'
import { getBackendBaseUrl } from '@/lib/env'
import { useToastMessages } from '@/lib/toast-helpers'

const severityMap: Record<BackendNotificationLevel, NotificationSeverity> = {
  critical: 'critical',
  high: 'high',
  medium: 'medium',
  low: 'low',
}

const inferNotificationType = (message: string, category?: string) => {
  // 优先使用后端返回的 category
  if (category === 'scan' || category === 'vulnerability' || category === 'asset' || category === 'system') {
    return category
  }
  // 后备：通过消息内容推断
  if (message?.includes('扫描') || message?.includes('任务')) {
    return 'scan' as const
  }
  if (message?.includes('漏洞')) {
    return 'vulnerability' as const
  }
  return 'system' as const
}

const formatTimeAgo = (date: Date): string => {
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / (1000 * 60))
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))

  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins} 分钟前`
  if (diffHours < 24) return `${diffHours} 小时前`
  return date.toLocaleDateString()
}

const isBackendNotification = (value: unknown): value is BackendNotification => {
  if (!value || typeof value !== 'object') return false
  const record = value as Record<string, unknown>
  return typeof record.id === 'number' && typeof record.title === 'string' && typeof record.message === 'string'
}

export const transformBackendNotification = (backendNotification: BackendNotification): Notification => {
  const createdAtRaw = backendNotification.createdAt ?? backendNotification.created_at
  const createdDate = createdAtRaw ? new Date(createdAtRaw) : new Date()
  const isRead = backendNotification.isRead ?? backendNotification.is_read
  const notification: Notification = {
    id: backendNotification.id,
    type: inferNotificationType(backendNotification.message, backendNotification.category),
    title: backendNotification.title,
    description: backendNotification.message,
    time: formatTimeAgo(createdDate),
    unread: isRead === true ? false : true,
    severity: severityMap[backendNotification.level] ?? undefined,
    createdAt: createdDate.toISOString(),
  }
  return notification
}

export function useNotificationSSE() {
  const [isConnected, setIsConnected] = useState(false)
  const [notifications, setNotifications] = useState<Notification[]>([])
  const wsRef = useRef<WebSocket | null>(null)
  const queryClient = useQueryClient()
  const reconnectTimerRef = useRef<NodeJS.Timeout | null>(null)
  const heartbeatTimerRef = useRef<NodeJS.Timeout | null>(null)
  const isConnectingRef = useRef(false)
  const reconnectAttempts = useRef(0)
  const maxReconnectAttempts = 10
  const baseReconnectDelay = 1000 // 1秒
  const toastMessages = useToastMessages()

  const markNotificationsAsRead = useCallback((ids?: number[]) => {
    setNotifications(prev => prev.map(notification => {
      if (!ids || ids.includes(notification.id)) {
        return { ...notification, unread: false }
      }
      return notification
    }))
  }, [])

  // 启动心跳
  const startHeartbeat = useCallback(() => {
    // 清除旧的心跳定时器
    if (heartbeatTimerRef.current) {
      clearInterval(heartbeatTimerRef.current)
    }

    // 每 30 秒发送一次心跳
    heartbeatTimerRef.current = setInterval(() => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'ping' }))
      }
    }, 30000) // 30秒
  }, [])

  // 停止心跳
  const stopHeartbeat = useCallback(() => {
    if (heartbeatTimerRef.current) {
      clearInterval(heartbeatTimerRef.current)
      heartbeatTimerRef.current = null
    }
  }, [])

  // 计算重连延迟（指数退避）
  const getReconnectDelay = useCallback(() => {
    const delay = Math.min(baseReconnectDelay * Math.pow(2, reconnectAttempts.current), 30000)
    return delay
  }, [])

  // 连接 WebSocket
  const connect = useCallback(() => {
    // 防止重复连接
    if (isConnectingRef.current) {
      return
    }

    // 如果已经连接，跳过
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return
    }

    isConnectingRef.current = true

    // 关闭旧连接
    if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
      wsRef.current.close()
    }

    // 清除重连定时器
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current)
      reconnectTimerRef.current = null
    }

    try {
      // 构造 WebSocket URL
      const backendUrl = getBackendBaseUrl()
      const wsProtocol = backendUrl.startsWith('https') ? 'wss' : 'ws'
      const wsHost = backendUrl.replace(/^https?:\/\//, '')
      const wsUrl = `${wsProtocol}://${wsHost}/ws/notifications/`


      const ws = new WebSocket(wsUrl)
      wsRef.current = ws

      ws.onopen = () => {
        setIsConnected(true)
        isConnectingRef.current = false
        reconnectAttempts.current = 0 // 重置重连计数
        // 启动心跳
        startHeartbeat()
      }

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as Record<string, unknown>

          const messageType = typeof data.type === 'string' ? data.type : undefined

          if (messageType === 'connected') {
            return
          }

          if (messageType === 'pong') {
            // 心跳响应
            return
          }

          if (messageType === 'error') {
            const errorMessage = typeof data.message === 'string' ? data.message : ''
            toastMessages.error('toast.notification.connection.error', { message: errorMessage })
            return
          }

          // 处理通知消息
          if (messageType === 'notification') {
            if (isBackendNotification(data)) {
              const notification = transformBackendNotification(data)
              setNotifications(prev => {
                const updated = [notification, ...prev.slice(0, 49)]
                return updated
              })

              queryClient.invalidateQueries({ queryKey: ['notifications'] })
            }
            return
          }

          // 备用处理：直接检查通知字段
          if (isBackendNotification(data)) {
            const notification = transformBackendNotification(data)

            // 添加到通知列表
            setNotifications(prev => {
              const updated = [notification, ...prev.slice(0, 49)]
              return updated
            })

            // 刷新通知查询
            queryClient.invalidateQueries({ queryKey: ['notifications'] })
          }
        } catch (error) {
          void error
        }
      }

      ws.onerror = () => {
        // WebSocket onerror 接收的是 Event 对象，不是 Error
        // 实际的错误信息通常不可用，只能记录连接状态
        setIsConnected(false)
        isConnectingRef.current = false
      }

      ws.onclose = (event) => {
        setIsConnected(false)
        isConnectingRef.current = false
        // 停止心跳
        stopHeartbeat()

        // 自动重连（非正常关闭时）
        if (event.code !== 1000) { // 1000 = 正常关闭
          if (reconnectAttempts.current < maxReconnectAttempts) {
            const delay = getReconnectDelay()
            reconnectAttempts.current++
            reconnectTimerRef.current = setTimeout(() => {
              connect()
            }, delay)
          }
        }
      }
    } catch {
      setIsConnected(false)
      isConnectingRef.current = false
    }
  }, [queryClient, startHeartbeat, stopHeartbeat, getReconnectDelay, toastMessages])

  // 断开连接
  const disconnect = useCallback(() => {
    // 停止心跳
    stopHeartbeat()

    // 清除重连定时器
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current)
      reconnectTimerRef.current = null
    }

    // 重置重连计数
    reconnectAttempts.current = 0
    isConnectingRef.current = false

    if (wsRef.current) {
      wsRef.current.close(1000, 'User disconnect') // 1000 = 正常关闭
      wsRef.current = null
    }
    setIsConnected(false)
  }, [stopHeartbeat])

  // 清空通知
  const clearNotifications = () => {
    setNotifications([])
  }

  // 组件挂载时连接，卸载时断开
  // 注意：不依赖 connect/disconnect 避免无限循环
  useEffect(() => {
    connect()

    return () => {
      disconnect()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return {
    isConnected,
    notifications,
    connect,
    disconnect,
    clearNotifications,
    markNotificationsAsRead,
  }
}
