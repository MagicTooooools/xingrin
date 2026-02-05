"use client"

import * as React from "react"
import {
  IconActivity,
  IconHeart,
  IconMoon,
  IconTrophy,
} from "@/components/icons"
import { useNudgeToast, type NudgeToastVariant } from "@/hooks/use-nudge-toast"

const STORAGE_FIRST_SEEN_KEY = "lunafox:milestone:first-seen"
const STORAGE_MILESTONE_PREFIX = "lunafox:milestone:triggered:"

// 定义里程碑（天数 -> 变体）
const MILESTONES: Record<number, NudgeToastVariant> = {
  1: {
    icon: <IconHeart className="size-6 text-pink-500" />,
    title: "Hello World! 👋",
    description: "这是我们在月狐控制台共度的第一天。很高兴认识你，指挥官。",
    primaryAction: { label: "你好呀" },
  },
  7: {
    icon: <IconTrophy className="size-6 text-yellow-500" />,
    title: "第一周达成！🏅",
    description: "已经过去一周了。你的资产库是不是也跟着胖了一圈？保持这个节奏！",
    primaryAction: { label: "继续冲" },
  },
  30: {
    icon: <IconMoon className="size-6 text-indigo-500" />,
    title: "满月纪念 🌕",
    description: "30 天的陪伴。这一个月里，感谢你为了互联网安全所做的每一次扫描。",
    primaryAction: { label: "干杯" },
  },
  100: {
    icon: <IconActivity className="size-6 text-emerald-500" />,
    title: "百日修仙达成 💯",
    description: "100 天的坚持。今天的你，一定比 100 天前更强了。",
    primaryAction: { label: "确实" },
    secondaryAction: { label: "还得练" },
  },
}

/**
 * 纯前端里程碑 Hook
 * 记录用户首次访问时间，并在特定天数（1, 7, 30, 100）触发一次性纪念弹窗
 */
export function useMilestoneNudge() {
  const [targetVariant, setTargetVariant] = React.useState<NudgeToastVariant[]>([])
  const [storageKey, setStorageKey] = React.useState<string | undefined>(undefined)

  // 初始化检查
  React.useEffect(() => {
    if (typeof window === "undefined") return

    const now = Date.now()
    const firstSeenStr = localStorage.getItem(STORAGE_FIRST_SEEN_KEY)

    // 1. 如果是第一次来，记录时间
    if (!firstSeenStr) {
      localStorage.setItem(STORAGE_FIRST_SEEN_KEY, now.toString())
      // 第 1 天的里程碑可以立刻触发（或者等下一次刷新）
      // 这里选择立刻触发 Day 1
      const variant = MILESTONES[1]
      const key = `${STORAGE_MILESTONE_PREFIX}1`
      if (!localStorage.getItem(key)) {
        setTargetVariant([variant])
        setStorageKey(key)
      }
      return
    }

    // 2. 计算已使用天数
    const firstSeen = parseInt(firstSeenStr, 10)
    const daysPassed = Math.floor((now - firstSeen) / (1000 * 60 * 60 * 24)) + 1

    // 3. 检查是否命中了某个里程碑
    const milestone = MILESTONES[daysPassed]
    if (milestone) {
      const key = `${STORAGE_MILESTONE_PREFIX}${daysPassed}`
      // 检查这个里程碑是否已经弹过（永久不弹，所以不用 cooldown，直接用 storageKey 控制）
      if (!localStorage.getItem(key)) {
        setTargetVariant([milestone])
        setStorageKey(key)
      }
    }
  }, [])

  // 使用通用 Hook 触发
  const { trigger } = useNudgeToast({
    storageKey, // 这里的 key 是动态的（e.g. triggered:30），弹过一次就会写入 true
    delay: 2000, // 稍微晚一点弹，让页面先加载完
    probability: 1, // 里程碑是硬性触发，不随机
    variants: targetVariant,
  })

  // 当检测到有待触发的里程碑时，执行
  React.useEffect(() => {
    if (targetVariant.length > 0 && storageKey) {
      trigger()
    }
  }, [targetVariant, storageKey, trigger])

  return null // 不需要暴露方法，全自动
}
