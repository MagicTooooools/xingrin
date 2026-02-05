"use client"

import * as React from "react"
import { toast } from "sonner"
import {
  IconActivity,
  IconAlertTriangle,
  IconCalendar,
  IconCoffee,
  IconCpu,
  IconHeart,
  IconLock,
  IconMoon,
  IconSun,
  IconTerminal,
  IconUser,
} from "@/components/icons"
import { useNudgeToast, type NudgeToastVariant } from "@/hooks/use-nudge-toast"

// 持久化 Key，记录上次触发时间，确保每天只弹一次
const STORAGE_KEY = "lunafox:care-nudge:last-seen"

// --- 场景定义 ---

// 1. 深夜 (23:00 - 05:00)
const LATE_NIGHT_VARIANTS: NudgeToastVariant[] = [
  {
    icon: <IconMoon className="size-6 text-indigo-500" />,
    title: "凌晨修仙警告",
    description: "指挥官，现在的流量最干净，但你的发际线正在报警。守护 Root 权限的同时，别忘了守护发量。",
    primaryAction: { label: "睡了睡了" },
    secondaryAction: { label: "再挖一个洞" },
  },
  {
    icon: <IconTerminal className="size-6 text-indigo-400" />,
    title: "It's Late...",
    description: "午夜钟声敲响，灰姑娘丢了水晶鞋，而你刚刚打开了 Burp Suite。",
    primaryAction: { label: "继续抓包" },
    secondaryAction: { label: "休息一下" },
  },
]

// 2. 清晨 (05:00 - 09:00)
const MORNING_VARIANTS: NudgeToastVariant[] = [
  {
    icon: <IconSun className="size-6 text-orange-500" />,
    title: "早起的黑客有洞挖",
    description: "当别人还在睡梦中，你已经完成了第一波资产测绘。今天的目标：拿下一个 Shell。",
    primaryAction: { label: "Get Shell!" },
  },
]

// 3. 下午茶 (14:00 - 16:00)
const COFFEE_VARIANTS: NudgeToastVariant[] = [
  {
    icon: <IconCoffee className="size-6 text-amber-700" />,
    title: "咖啡因浓度过低",
    description: "SQL 注入失败？也许不是 Payload 的问题，是你血液里的咖啡因不足了。来杯 Java 提提神。",
    primaryAction: { label: "去倒咖啡" },
    secondaryAction: { label: "红牛万岁" },
  },
]

// 4. 饭点 (11:30-13:30, 17:30-19:30)
const FOOD_VARIANTS: NudgeToastVariant[] = [
  {
    icon: <IconCpu className="size-6 text-rose-500" />,
    title: "缓冲区下溢警告",
    description: "检测到你的胃部缓冲区即将 Underflow。请立即执行 eat() 函数，防止身体 Crash。",
    primaryAction: { label: "去干饭" },
    secondaryAction: { label: "还能再扛" },
  },
]

// 5. 周五 & 假日 (Friday & Holiday)
const FRIDAY_VARIANTS: NudgeToastVariant[] = [
  {
    icon: <IconCalendar className="size-6 text-purple-500" />,
    title: "Read-Only Friday",
    description: "周五不改生产库，这是江湖规矩。除非你想在这个周末收到报警短信，否则请放下手中的 DROP TABLE。",
    primaryAction: { label: "遵命" },
    secondaryAction: { label: "头铁" },
  },
  {
    icon: <IconHeart className="size-6 text-pink-500" />,
    title: "节假日加班警告",
    description: "大过节的还在挖洞？你的敬业程度让我 CPU 温度都升高了。记得给自己发三倍工资（或者三倍快乐水）。",
    primaryAction: { label: "搞完就撤" },
    secondaryAction: { label: "我爱加班" },
  },
]

// 6. 熬夜加班 (Late Night Overtime 23:00 - 02:00)
const OVERTIME_VARIANTS: NudgeToastVariant[] = [
  {
    icon: <IconMoon className="size-6 text-indigo-500" />,
    title: "还在肝？",
    description: "听说深夜的代码 Bug 更少？也许吧，但发际线后移的速度肯定更快。保重啊，指挥官。",
    primaryAction: { label: "最后一行" },
    secondaryAction: { label: "再战三百回" },
  },
  {
    icon: <IconCoffee className="size-6 text-amber-700" />,
    title: "深夜食堂",
    description: "这个点还在屏幕前，是不是饿了？泡面虽好，可不要贪杯。或者去阳台看看星星（如果有的话）。",
    primaryAction: { label: "去觅食" },
    secondaryAction: { label: "我不饿" },
  },
]

// 6. 随机梗 (通用)
const MEME_VARIANTS: NudgeToastVariant[] = [
  {
    icon: <IconAlertTriangle className="size-6 text-red-600" />,
    title: "rm -rf /*",
    description: "这是一个危险的命令... 哪怕你是在虚拟机里。时刻保持敬畏，记得备份，记得看清你在哪个窗口。",
    primaryAction: { label: "已经在跑路了" },
  },
  {
    icon: <IconLock className="size-6 text-blue-500" />,
    title: "密码是 123456？",
    description: "希望不是。否则我的字典生成器第一秒就能跑出来。去改个强密码吧，比如 P@$$w0rd (开玩笑的)。",
    primaryAction: { label: "我很强" },
    secondaryAction: { label: "这就改" },
  },
  {
    icon: <IconUser className="size-6 text-pink-500" />,
    title: "面向对象编程",
    description: "找不到对象？没关系，new 一个就行了。但在现实里，可能需要你合上电脑出去走走。",
    primaryAction: { label: "我去 new 一个" },
    secondaryAction: { label: "代码就是恋人" },
  },
  {
    icon: <IconActivity className="size-6 text-emerald-500" />,
    title: "姿势不对，Shell 不会",
    description: "长期低头会压迫颈椎，导致 Payload 构造思路受阻。建议立即执行：stand_up(); stretch();",
    primaryAction: { label: "活动一下" },
    secondaryAction: { label: "再坐五百年" },
  },
  {
    icon: <IconHeart className="size-6 text-red-500" />,
    title: "别扫了，扫我吧",
    description: "我是你的小狐狸助手，一直在默默守护你的控制台。要不要给我点个 Star？",
    primaryAction: { 
      label: "去点 Star",
      onClick: () => window.open("https://github.com/yyhuni/xingrin", "_blank") 
    },
    secondaryAction: { label: "下次一定" },
  },
]

interface UseCareNudgeOptions {
  /**
   * 触发概率 (0-1)，用于控制随机弹出的频率
   * @default 0.3 (30% 概率触发，避免太烦人)
   */
  probability?: number
  /**
   * 延迟触发时间 (ms)
   * @default 3000
   */
  delay?: number
}

/**
 * 智能关怀 Hook
 * 根据时间、日期、随机事件，给予用户 Hacker 风格的关怀提示
 * 每天最多只触发一次 (依靠 cooldownMs = 24h)
 */
export function useCareNudge(options: UseCareNudgeOptions = {}) {
  const { probability = 0.3, delay = 3000 } = options

  // 根据当前环境计算合适的 variants
  const variants = React.useMemo(() => {
    const now = new Date()
    const hour = now.getHours()
    const day = now.getDay() // 0 = Sunday, 5 = Friday

    let pool: NudgeToastVariant[] = [...MEME_VARIANTS]

    // 1. 时间判断
    if (hour >= 23 || hour <= 4) {
      // 深夜：高优先级，清空其他，只弹深夜/熬夜关怀
      // 混合 LATE_NIGHT 和 OVERTIME
      return [...LATE_NIGHT_VARIANTS, ...OVERTIME_VARIANTS]
    } else if (hour >= 5 && hour <= 9) {
      pool = [...pool, ...MORNING_VARIANTS]
    } else if ((hour >= 11 && hour <= 13) || (hour >= 17 && hour <= 19)) {
      pool = [...pool, ...FOOD_VARIANTS]
    } else if (hour >= 14 && hour <= 16) {
      pool = [...pool, ...COFFEE_VARIANTS]
    }

    // 2. 日期判断
    if (day === 5 || day === 0 || day === 6) {
      // 周五或周末：加入假日/加班专属梗
      pool = [...pool, ...FRIDAY_VARIANTS]
    }

    // TODO: 节日判断 (1024, 春节等) - 暂略，可扩展

    return pool
  }, [])

  // 这里的 cooldownMs 设为 20小时左右，保证每天大概能看到一次（如果用户每天都来的话）
  // 或者设为 12小时，涵盖早晚
  const COOLDOWN_MS = 16 * 60 * 60 * 1000

  // AI Proxy API URL (hardcoded)
  const aiApiUrl = "https://lunafox-ai-proxy-fzqn2tz4eb4f.yyhunisec.deno.net/"

  const { trigger, reset } = useNudgeToast({
    storageKey: STORAGE_KEY,
    cooldownMs: COOLDOWN_MS,
    probability,
    delay,
    duration: 8000,
    position: "bottom-right",
    variants, // 本地 variants 作为 fallback (useNudgeToast 内部如果拿到 variants 就会用)
  })

  // 包装 trigger：优先尝试 AI 生成，失败则降级到本地静态
  const triggerWithAi = React.useCallback(async () => {
    if (!aiApiUrl) {
      trigger() // 使用本地 variants
      return
    }

    try {
      // 检查 CD (简化版)
      const lastSeen = localStorage.getItem(STORAGE_KEY)
      if (lastSeen && Number(lastSeen) > Date.now()) return // 还在冷却中

      const now = new Date()
      const context = {
        hour: now.getHours(),
        day: now.getDay(),
        event: "daily_care",
      }

      const res = await fetch(aiApiUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ context }),
      })

      if (!res.ok) throw new Error("AI API failed")

      const aiData = await res.json()
      
      // AI 成功，渲染 Toast
      const { NudgeToastCard } = await import("@/components/nudges/nudge-toast-card")
      
      toast.custom((t) => (
        <NudgeToastCard
          title={aiData.title}
          description={aiData.description}
          icon={<span className="text-2xl">{aiData.icon || "🤖"}</span>}
          primaryAction={{
            label: aiData.primaryAction?.label || "OK",
            onClick: () => {
              toast.dismiss(t)
              localStorage.setItem(STORAGE_KEY, String(Date.now() + COOLDOWN_MS))
            }
          }}
          secondaryAction={aiData.secondaryAction ? {
            label: aiData.secondaryAction.label,
            buttonVariant: "outline",
            onClick: () => {
              toast.dismiss(t)
              localStorage.setItem(STORAGE_KEY, String(Date.now() + COOLDOWN_MS))
            }
          } : undefined}
          onDismiss={() => {
            toast.dismiss(t)
            localStorage.setItem(STORAGE_KEY, String(Date.now() + COOLDOWN_MS))
          }}
        />
      ), { duration: 8000, position: "bottom-right" })

    } catch (err) {
      console.warn("AI Nudge failed, falling back to static:", err)
      trigger() // 降级到本地静态
    }

  }, [aiApiUrl, trigger, COOLDOWN_MS])

  // 自动触发：组件挂载后尝试触发
  React.useEffect(() => {
    // 延迟一点点执行，避免和服务端渲染冲突
    const timer = setTimeout(() => {
      triggerWithAi()
    }, delay)
    return () => clearTimeout(timer)
  }, [triggerWithAi, delay])

  return { trigger: triggerWithAi, reset }
}
