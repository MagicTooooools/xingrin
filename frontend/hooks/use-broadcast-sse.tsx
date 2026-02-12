"use client"

import { useEffect, useRef, useCallback } from "react"
import { toast } from "sonner"

/**
 * LunaFox SSE 广播推送 Hook
 * 监听固定的 Deno SSE 服务，收到广播消息时弹出 Toast
 */

// 固定的 SSE 服务地址 - 只允许这个 URL
const SSE_URL = "https://lunafox-sse-push.yyhunisec.deno.net/sse"

// localStorage key - 用于 "不再提醒" 功能
const SUPPRESS_KEY = "lunafox:broadcast-suppress"

// 抑制时间（毫秒）- 点击不再提醒后 24 小时内不弹窗
const SUPPRESS_DURATION = 24 * 60 * 60 * 1000

interface BroadcastMessage {
  type: "broadcast" | "heartbeat" | "connected"
  icon?: string
  title?: string
  description?: string
  primaryAction?: { label: string; href?: string }
  secondaryAction?: { label: string; href?: string } | null
  timestamp?: number
}

let nudgeToastCardLoader: Promise<
  (typeof import("@/components/nudges/nudge-toast-card"))["NudgeToastCard"]
> | null = null

function loadNudgeToastCard() {
  if (!nudgeToastCardLoader) {
    nudgeToastCardLoader = import("@/components/nudges/nudge-toast-card").then(
      (mod) => mod.NudgeToastCard
    )
  }
  return nudgeToastCardLoader
}

function isSuppressed(): boolean {
  try {
    const raw = localStorage.getItem(SUPPRESS_KEY)
    if (!raw) return false
    const until = Number(raw)
    if (Date.now() < until) return true
    localStorage.removeItem(SUPPRESS_KEY)
    return false
  } catch {
    return false
  }
}

function suppress() {
  try {
    localStorage.setItem(SUPPRESS_KEY, String(Date.now() + SUPPRESS_DURATION))
  } catch {
    // ignore
  }
}

// 固定的 toast ID，用于新消息覆盖旧弹窗
const BROADCAST_TOAST_ID = "lunafox-broadcast"

export function useBroadcastSSE() {
  const eventSourceRef = useRef<EventSource | null>(null)
  const reconnectTimerRef = useRef<number | null>(null)
  const reconnectAttemptsRef = useRef(0)
  const MAX_RECONNECT_DELAY = 60000 // 最大 60 秒

  const showBroadcastToast = useCallback((data: BroadcastMessage) => {
    // 检查是否被抑制
    if (isSuppressed()) return

    // 先关闭之前的弹窗
    toast.dismiss(BROADCAST_TOAST_ID)

    // 稍微延迟显示新弹窗，确保旧的已关闭
    setTimeout(() => {
      void loadNudgeToastCard().then((NudgeToastCard) => {
        toast.custom(
          (t) => (
            <NudgeToastCard
              title={data.title || "系统通知"}
              description={data.description || ""}
              icon={<span className="text-2xl">{data.icon || "📢"}</span>}
              primaryAction={{
                label: data.primaryAction?.label || "知道了",
                href: data.primaryAction?.href,
                onClick: () => toast.dismiss(t),
              }}
              secondaryAction={
                data.secondaryAction
                  ? {
                      label: data.secondaryAction.label,
                      href: data.secondaryAction.href,
                      buttonVariant: "outline",
                      onClick: () => {
                        suppress()
                        toast.dismiss(t)
                      },
                    }
                  : undefined
              }
              onDismiss={() => toast.dismiss(t)}
            />
          ),
          {
            id: BROADCAST_TOAST_ID,
            duration: 15000,
            position: "bottom-right",
          }
        )
      })
    }, 100)
  }, [])

  const connect = useCallback(() => {
    // 防止重复连接
    if (eventSourceRef.current?.readyState === EventSource.OPEN) {
      return
    }

    // 清理旧连接
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
    }

    try {
      const es = new EventSource(SSE_URL)
      eventSourceRef.current = es

      es.onopen = () => {
        // 连接成功，重置重试计数器
        reconnectAttemptsRef.current = 0
      }

      es.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as BroadcastMessage

          // 只处理广播消息
          if (data.type === "broadcast") {
            showBroadcastToast(data)
          }
          // heartbeat 和 connected 消息忽略
        } catch {
          // JSON 解析失败，忽略
        }
      }

      es.onerror = () => {
        es.close()
        eventSourceRef.current = null

        // 指数退避: 5s, 10s, 20s, 40s, 60s (max)
        const delay = Math.min(
          5000 * Math.pow(2, reconnectAttemptsRef.current),
          MAX_RECONNECT_DELAY
        )
        reconnectAttemptsRef.current++

        reconnectTimerRef.current = window.setTimeout(() => {
          connect()
        }, delay)
      }
    } catch {
      // 连接失败，忽略
    }
  }, [showBroadcastToast])

  const disconnect = useCallback(() => {
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current)
      reconnectTimerRef.current = null
    }

    if (eventSourceRef.current) {
      eventSourceRef.current.close()
      eventSourceRef.current = null
    }
  }, [])

  // 组件挂载时连接，卸载时断开
  useEffect(() => {
    connect()
    return () => disconnect()
  }, [connect, disconnect])

  return { connect, disconnect }
}
