"use client"

import { useEffect, useState } from "react"
import type { LucideIcon } from "@/components/icons"
import { cn } from "@/lib/utils"

interface BauhausPageHeaderProps {
  /** 页面代号，如 "TGT-01" */
  code: string
  /** 副标题/分类，如 "Asset Management" */
  subtitle: string
  /** 主标题 */
  title: string
  /** 描述信息（可选） */
  description?: string
  /** 是否显示描述（默认 false） */
  showDescription?: boolean
  /** 图标组件 */
  icon?: LucideIcon
  /** 是否显示图标（默认 false） */
  showIcon?: boolean
  /** 状态文本，默认 "ACTIVE" */
  statusText?: string
  /** 是否在线状态，默认 true */
  isOnline?: boolean
  /** 外层容器的自定义 class */
  className?: string
}

/**
 * Bauhaus 风格的页面头部组件
 * 仅在 Bauhaus 主题下显示，提供统一的工业风格视觉效果
 */
export function BauhausPageHeader({
  code,
  subtitle,
  title,
  description,
  showDescription = false,
  icon: Icon,
  showIcon = false,
  statusText = "ACTIVE",
  isOnline = true,
  className,
}: BauhausPageHeaderProps) {
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
    updateTime()
    const interval = setInterval(updateTime, 1000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className={cn("hidden [[data-theme=bauhaus]_&]:block px-4 lg:px-6", className)}>
      <div className="bg-card border border-border border-t-2 border-t-primary p-5">
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
          {/* 左侧标题区 */}
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <span className="px-1.5 py-0.5 text-[10px] font-mono bg-primary text-primary-foreground tracking-wider">
                {code}
              </span>
              <p className="text-xs tracking-[0.2em] text-muted-foreground uppercase">
                {subtitle}
              </p>
            </div>
            <h1 className="text-2xl font-bold tracking-tight uppercase flex items-center gap-3">
              {showIcon && Icon ? <Icon className="h-6 w-6" /> : null}
              {title}
            </h1>
            {showDescription && description ? (
              <p className="text-sm text-muted-foreground">{description}</p>
            ) : null}
          </div>

          {/* 右侧状态栏 */}
          <div className="flex gap-2">
            <div className="px-3 py-1.5 flex items-center gap-2 text-xs font-mono bg-secondary border border-border">
              <span
                className={`w-1.5 h-1.5 rounded-sm ${
                  isOnline ? "bg-[var(--success)]" : "bg-[var(--error)]"
                }`}
              />
              STATUS: {statusText}
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
