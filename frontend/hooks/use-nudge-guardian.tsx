"use client"

import * as React from "react"
import { toast } from "sonner"
import { NudgeToastCard } from "@/components/nudges/nudge-toast-card"
import { 
  IconMoon, 
  IconCoffee, 
  IconActivity, 
  IconSun, 
  IconCpu,
  IconBug // 新增：用于代表赛博生物
} from "@/components/icons"

// 存储 Key
const KEY_LAST_NUDGE = "lunafox:last-nudge" // { [ruleId]: timestamp }
const SESSION_START = "lunafox:session-start"

type RuleId = "late_night" | "long_session" | "cyber_fox"

// --- 文案库 ---
const VARIANTS = {
  // 1. 赛博狐狸 (卖萌求关注)
  cyber_fox: [
    {
      title: "系统监测到毛茸茸生物",
      desc: "🦊 嘤？LunaFox 的吉祥物跑出来了！它似乎对你的鼠标指针很好奇。工作累了就陪它玩会儿吧。",
      icon: IconBug, 
      color: "text-orange-500",
      primary: "陪它玩",
      secondary: "摸摸头"
    }
  ],
  // 2. 深夜修仙
  late_night: [
    {
      title: "深夜修仙警告",
      desc: "指挥官，流量最干净的时候确实适合挖洞，但发际线也在报警。守护 Root 权限的同时，别忘了守护发量。",
      icon: IconMoon,
      color: "text-indigo-500",
      primary: "再熬半小时",
      secondary: "这就睡"
    },
    {
      title: "夜深了",
      desc: "这个点还在屏幕前，是不是饿了？泡面虽好，可不要贪杯。或者去阳台看看星星（如果有的话），然后关机。",
      icon: IconCoffee,
      color: "text-amber-600",
      primary: "去看星星",
      secondary: "关机"
    },
    {
      title: "凌晨 1 点的月狐",
      desc: "你和服务器是现在唯一醒着的伙伴。别太累了，服务器有 UPS 撑着，你可没有。",
      icon: IconMoon,
      color: "text-purple-500",
      primary: "陪它聊会儿",
      secondary: "去休息"
    }
  ],
  // 3. 久坐提醒
  long_session: [
    {
      title: "姿势不对，Shell 不会",
      desc: "长期低头会压迫颈椎，导致 Payload 构造思路受阻。建议立即执行：stand_up(); stretch();",
      icon: IconActivity,
      color: "text-emerald-500",
      primary: "起来活动",
      secondary: "再坐会儿"
    },
    {
      title: "视觉模块过热",
      desc: "你盯着屏幕太久了，视网膜可能出现了残影。去窗边看看有没有真的狐狸，或者哪怕只是看看那棵树。",
      icon: IconSun,
      color: "text-amber-500",
      primary: "去看狐狸",
      secondary: "看看树"
    },
    {
      title: "久坐提醒",
      desc: "已经连续战斗 2 小时了。起来接杯水，活动一下腰背，为了更长远的黑客生涯。",
      icon: IconCpu,
      color: "text-indigo-500",
      primary: "接水去",
      secondary: "活动腰背"
    },
    {
      title: "Drink Water",
      desc: "多喝水，多排毒。身体是革命的本钱，也是挖洞的本钱。",
      icon: IconCoffee,
      color: "text-cyan-500",
      primary: "去接水",
      secondary: "一会儿喝"
    }
  ]
}

export function useNudgeGuardian() {
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
      } catch {}

      // 辅助：检查冷却时间 (默认 16 小时内不重复触发同一规则)
      const isCoolingDown = (id: RuleId, cooldownHours = 16) => {
        const last = lastNudges[id] || 0
        return Date.now() - last < cooldownHours * 60 * 60 * 1000
      }

      // 辅助：随机选择文案并触发
      const trigger = (id: RuleId, cooldownOverride?: number) => {
        if (isCoolingDown(id, cooldownOverride)) return

        // 随机选择一条文案
        const variants = VARIANTS[id]
        const content = variants[Math.floor(Math.random() * variants.length)]

        // 记录触发时间
        lastNudges[id] = Date.now()
        localStorage.setItem(KEY_LAST_NUDGE, JSON.stringify(lastNudges))

        // 弹窗
        toast.dismiss() 
        toast.custom((t) => (
          <NudgeToastCard
            title={content.title}
            description={content.desc}
            icon={<content.icon className={`size-8 ${content.color}`} />}
            primaryAction={{ 
              label: content.primary || "收到", 
              onClick: () => toast.dismiss(t) 
            }}
            secondaryAction={
              content.secondary
                ? {
                    label: content.secondary,
                    onClick: () => toast.dismiss(t),
                    buttonVariant: "outline",
                  }
                : undefined
            }
            onDismiss={() => toast.dismiss(t)}
          />
        ), { duration: Infinity, position: "bottom-right" })
      }

      // --- 规则判断 ---

      // 1. 深夜修仙 (23:00 - 04:00)
      if ((hour >= 23 || hour < 4)) {
        trigger("late_night")
      }

      // 2. 久坐提醒 (连续 2 小时) - 冷却时间 3 小时
      // 为了测试方便，您可以暂时把 120 改成 1 (1分钟)
      if (sessionDurationMinutes > 120) {
        trigger("long_session", 3)
      }

      // 3. 赛博狐狸 (15:00-16:00, 1% 概率)
      // 冷却时间 20 小时，保证每天最多偶遇一次
      if (hour === 15 && Math.random() < 0.01) {
        trigger("cyber_fox", 20)
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
  }, [])
}
