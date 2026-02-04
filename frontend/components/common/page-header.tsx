"use client"

import { cn } from "@/lib/utils"

interface PageHeaderProps {
  /** 页面代号，如 "TGT-01" */
  code: string
  /** 主标题 */
  title: string
  /** 描述信息（可选） */
  description?: string
  /** 外层容器的自定义 class */
  className?: string
}

/**
 * 工业风页面头部组件 - 极简下划线风格 (Option C)
 * 仅保留底部强调线，去除背景和边框，强调内容本身的层次
 */
export function PageHeader({
  code,
  title,
  description,
  className,
}: PageHeaderProps) {
  return (
    <div className={cn("px-4 lg:px-6 mb-2", className)}>
      <div className="flex items-end gap-2 mb-2">
        <div className="flex items-baseline gap-3 border-b-2 border-primary pb-2">
          <h1 className="text-2xl font-bold tracking-tight uppercase">
            {title}
          </h1>
          <span className="font-mono text-xs text-muted-foreground font-medium tracking-wide">
            /{code}
          </span>
        </div>
        <div className="flex-1 h-1.5 bg-[repeating-linear-gradient(45deg,transparent,transparent_4px,currentColor_4px,currentColor_5px)] text-primary/10 mb-0.5" />
      </div>
      {description && (
        <p className="text-sm text-muted-foreground">
          {description}
        </p>
      )}
    </div>
  )
}
