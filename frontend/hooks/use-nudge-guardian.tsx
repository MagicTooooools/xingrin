"use client"

import * as React from "react"
import {
  GUARDIAN_RULES,
  type GuardianRuleContext,
  type GuardianRuleId,
  type GuardianVariantTemplate,
} from "@/hooks/nudge-rules"
import { useNudgeToast, type NudgeToastVariant } from "@/hooks/use-nudge-toast"

// 存储 Key
const KEY_LAST_NUDGE = "lunafox:last-nudge" // { [ruleId]: timestamp }
const SESSION_START = "lunafox:session-start"

export function useNudgeGuardian() {
  const { triggerWithVariant } = useNudgeToast({
    variants: [],
    probability: 1,
    delay: 0,
    duration: Infinity,
    position: "bottom-right",
  })

  const toToastVariant = React.useCallback((content: GuardianVariantTemplate): NudgeToastVariant => {
    return {
      title: content.title,
      description: content.desc,
      icon: <content.icon className={`size-8 ${content.color}`} />,
      primaryAction: {
        label: content.primary || "收到",
      },
      secondaryAction: content.secondary
        ? {
            label: content.secondary,
            buttonVariant: "outline",
          }
        : undefined,
    }
  }, [])

  React.useEffect(() => {
    // 1. 初始化会话时间
    if (!sessionStorage.getItem(SESSION_START)) {
      sessionStorage.setItem(SESSION_START, String(Date.now()))
    }

    // 检查逻辑
    const checkRules = () => {
      const now = new Date()
      const hour = now.getHours()

      const sessionStart = Number(sessionStorage.getItem(SESSION_START) || Date.now())
      const sessionDurationMinutes = (Date.now() - sessionStart) / 1000 / 60

      // 获取上次触发记录
      let lastNudges: Record<string, number> = {}
      try {
        lastNudges = JSON.parse(localStorage.getItem(KEY_LAST_NUDGE) || "{}")
      } catch {
        lastNudges = {}
      }

      // 辅助：检查冷却时间
      const isCoolingDown = (id: GuardianRuleId, cooldownHours: number) => {
        const last = lastNudges[id] || 0
        return Date.now() - last < cooldownHours * 60 * 60 * 1000
      }

      const context: GuardianRuleContext = {
        hour,
        sessionDurationMinutes,
      }

      for (const rule of GUARDIAN_RULES) {
        if (!rule.shouldTrigger(context)) continue
        if (isCoolingDown(rule.id, rule.cooldownHours)) continue
        if (rule.variants.length === 0) continue

        const content = rule.variants[Math.floor(Math.random() * rule.variants.length)]
        if (!content) continue

        // 记录触发时间
        lastNudges[rule.id] = Date.now()
        localStorage.setItem(KEY_LAST_NUDGE, JSON.stringify(lastNudges))

        triggerWithVariant(toToastVariant(content))
      }
    }

    // 启动定时器：每分钟检查一次
    const timer = setInterval(checkRules, 60 * 1000)

    // 首次加载延迟 3 秒检查一次
    const initialTimer = setTimeout(checkRules, 3000)

    return () => {
      clearInterval(timer)
      clearTimeout(initialTimer)
    }
  }, [toToastVariant, triggerWithVariant])
}
