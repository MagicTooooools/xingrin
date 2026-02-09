"use client"

import { useEffect, useState } from "react"

/**
 * Bauhaus 风格的 Dashboard Header
 * 仅在 Bauhaus 主题下显示，模仿 dashboard-demo 原型的视觉效果
 */
export function BauhausDashboardHeader() {
  const [currentTime, setCurrentTime] = useState<string>("")

  useEffect(() => {
    const updateTime = () => {
      const now = new Date()
      setCurrentTime(
        now.toLocaleTimeString("en-US", {
          hour12: false,
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        })
      )
    }

    let intervalId: number | null = null

    const shouldRun = () =>
      !document.hidden && document.documentElement.getAttribute("data-theme") === "bauhaus"

    const start = () => {
      updateTime()
      intervalId = window.setInterval(updateTime, 1000)
    }

    const stop = () => {
      if (intervalId) {
        clearInterval(intervalId)
        intervalId = null
      }
    }

    const sync = () => {
      if (shouldRun()) {
        if (!intervalId) start()
      } else {
        stop()
      }
    }

    sync()

    const observer = new MutationObserver(sync)
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ["data-theme"] })
    document.addEventListener("visibilitychange", sync)

    return () => {
      document.removeEventListener("visibilitychange", sync)
      observer.disconnect()
      stop()
    }
  }, [])

  return (
    <div className="hidden [[data-theme=bauhaus]_&]:block px-4 lg:px-6">
      <div className="bg-card border border-border border-t-2 border-t-primary p-5">
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
          {/* 左侧标题区 */}
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <span className="px-1.5 py-0.5 text-[10px] font-mono bg-primary text-primary-foreground tracking-wider">
                DASH-01
              </span>
              <p className="text-xs tracking-[0.2em] text-muted-foreground uppercase">
                Operations Command
              </p>
            </div>
            <h1 className="text-2xl font-bold tracking-tight uppercase">
              System Overview
            </h1>
          </div>

          {/* 右侧状态栏 */}
          <div className="flex gap-2">
            <div className="px-3 py-1.5 flex items-center gap-2 text-xs font-mono bg-secondary border border-border">
              <span className="w-1.5 h-1.5 bg-[var(--success)] rounded-sm" />
              NETWORK: ONLINE
            </div>
            <div className="px-3 py-1.5 flex items-center gap-2 text-xs font-mono bg-secondary border border-border">
              <span className="text-muted-foreground">CYCLE:</span>
              {currentTime}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
