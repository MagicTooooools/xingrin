"use client"

import * as React from "react"
import { toast, type ToasterProps } from "sonner"

import { NudgeToastCard, type NudgeToastCardProps } from "@/components/nudges/nudge-toast-card"
import {
  isLocalStorageAvailable,
  isNudgeSuppressed,
  suppressNudge,
  withNudgeDismiss,
} from "@/lib/nudge-toast-helpers"

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

const withDismiss = withNudgeDismiss

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

  const showVariant = React.useCallback((variant: NudgeToastVariant) => {
    toast.custom(
      (t) => {
        const onDismiss = () => {
          toast.dismiss(t)
          if (storageKey) suppressNudge(storageKey, cooldownMs)
        }

        const primaryAction = withDismiss(variant.primaryAction, onDismiss) ?? {
          label: "OK",
          onClick: onDismiss,
        }

        return (
          <NudgeToastCard
            {...variant}
            onDismiss={onDismiss}
            primaryAction={primaryAction}
            secondaryAction={withDismiss(variant.secondaryAction, onDismiss)}
          />
        )
      },
      {
        duration,
        position,
      }
    )
  }, [cooldownMs, duration, position, storageKey])

  const triggerInternal = React.useCallback((variantOverride?: NudgeToastVariant) => {
    if (typeof window === "undefined") return
    if (!variantOverride && (!variants || variants.length === 0)) return

    // Probability control
    if (Math.random() > probability) return

    // Storage suppression
    if (storageKey && isNudgeSuppressed(storageKey, cooldownMs)) return

    const variant = variantOverride ?? variants[Math.floor(Math.random() * variants.length)]
    if (!variant) return

    // Deduplicate pending triggers
    if (timerRef.current !== null) {
      window.clearTimeout(timerRef.current)
      timerRef.current = null
    }

    timerRef.current = window.setTimeout(() => {
      showVariant(variant)
    }, delay)
  }, [cooldownMs, delay, probability, showVariant, storageKey, variants])

  const trigger = React.useCallback(() => {
    triggerInternal()
  }, [triggerInternal])

  const triggerWithVariant = React.useCallback((variant: NudgeToastVariant) => {
    triggerInternal(variant)
  }, [triggerInternal])

  const isSuppressedNow = React.useCallback(() => {
    if (!isLocalStorageAvailable()) return true
    if (!storageKey) return false

    return isNudgeSuppressed(storageKey, cooldownMs)
  }, [cooldownMs, storageKey])

  const reset = React.useCallback(() => {
    if (!isLocalStorageAvailable()) return
    if (!storageKey) return

    try {
      localStorage.removeItem(storageKey)
    } catch {
      // ignore
    }
  }, [storageKey])

  return { trigger, triggerWithVariant, isSuppressedNow, reset }
}
