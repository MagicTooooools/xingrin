"use client"

import * as React from "react"
import { toast, type ToasterProps } from "sonner"

import {
  NudgeToastCard,
  type NudgeToastAction,
  type NudgeToastCardProps,
} from "@/components/nudges/nudge-toast-card"

export type NudgeToastVariant = Omit<NudgeToastCardProps, "onDismiss">

export interface UseNudgeToastOptions {
  /**
   * 持久化抑制 key。
   * - 不传：不做持久化抑制（每次 trigger 都可能弹）
   * - 传入：会在用户关闭/点击按钮后写入 localStorage
   */
  storageKey?: string

  /**
   * 冷却时间 (ms)。
   * - 不传：写入 "true"，永久不再弹（直到 reset）
   * - 传入：写入 nextAllowedAt 时间戳，到期后可再次弹
   */
  cooldownMs?: number

  /**
   * 触发概率 (0-1)，用于测试或灰度
   * @default 1
   */
  probability?: number

  /**
   * 延迟触发时间 (ms)
   * @default 1500
   */
  delay?: number

  /**
   * Toast 持续时间（Infinity 表示不自动关闭）
   * @default Infinity
   */
  duration?: number

  /**
   * Toast 位置
   * @default "bottom-right"
   */
  position?: ToasterProps["position"]

  /**
   * 可选的不同文案/样式变体
   */
  variants: NudgeToastVariant[]
}

function isSuppressed(storageKey: string, cooldownMs?: number): boolean {
  try {
    const raw = localStorage.getItem(storageKey)
    if (!raw) return false

    // Forever suppression
    if (!cooldownMs) return true

    // Cooldown suppression
    const nextAllowedAt = Number(raw)
    if (!Number.isFinite(nextAllowedAt)) return true

    if (nextAllowedAt > Date.now()) return true

    // Expired cooldown - allow again
    localStorage.removeItem(storageKey)
    return false
  } catch {
    // If storage is blocked, be conservative: treat as suppressed
    return true
  }
}

function suppress(storageKey: string, cooldownMs?: number) {
  try {
    if (!cooldownMs) {
      localStorage.setItem(storageKey, "true")
      return
    }

    localStorage.setItem(storageKey, String(Date.now() + cooldownMs))
  } catch {
    // ignore
  }
}

function withDismiss(
  action: NudgeToastAction | undefined,
  onDismiss: () => void
): NudgeToastAction | undefined {
  if (!action) return undefined

  const original = action.onClick
  return {
    ...action,
    onClick: () => {
      original?.()
      onDismiss()
    },
  }
}

export function useNudgeToast({
  storageKey,
  cooldownMs,
  probability = 1,
  delay = 1500,
  duration = Infinity,
  position = "bottom-right",
  variants,
}: UseNudgeToastOptions) {
  const timerRef = React.useRef<number | null>(null)

  React.useEffect(() => {
    return () => {
      if (timerRef.current !== null) {
        window.clearTimeout(timerRef.current)
        timerRef.current = null
      }
    }
  }, [])

  const trigger = React.useCallback(() => {
    if (typeof window === "undefined") return
    if (!variants || variants.length === 0) return

    // Probability control
    if (Math.random() > probability) return

    // Storage suppression
    if (storageKey && isSuppressed(storageKey, cooldownMs)) return

    const variant = variants[Math.floor(Math.random() * variants.length)]

    // Deduplicate pending triggers
    if (timerRef.current !== null) {
      window.clearTimeout(timerRef.current)
      timerRef.current = null
    }

    timerRef.current = window.setTimeout(() => {
      toast.custom(
        (t) => {
          const onDismiss = () => {
            toast.dismiss(t)
            if (storageKey) suppress(storageKey, cooldownMs)
          }

          return (
            <NudgeToastCard
              {...variant}
              onDismiss={onDismiss}
              primaryAction={withDismiss(variant.primaryAction, onDismiss)!}
              secondaryAction={withDismiss(variant.secondaryAction, onDismiss)}
            />
          )
        },
        {
          duration,
          position,
        }
      )
    }, delay)
  }, [cooldownMs, delay, duration, position, probability, storageKey, variants])

  const reset = React.useCallback(() => {
    if (typeof window === "undefined") return
    if (!storageKey) return

    try {
      localStorage.removeItem(storageKey)
    } catch {
      // ignore
    }
  }, [storageKey])

  return { trigger, reset }
}
