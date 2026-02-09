"use client"

import { cn } from "@/lib/utils"

interface PageHeaderProps {
  /** 页面代号，如 "TGT-01" */
  code?: string
  /** 主标题 */
  title: string
  /** 描述信息（可选） */
  description?: string
  /** 面包屑（兼容旧页面参数） */
  breadcrumbItems?: Array<{ label: string; href: string }>
  /** 外层容器的自定义 class */
  className?: string
  /** 右侧操作区域（可选） */
  action?: React.ReactNode
}

/**
 * 工业风页面头部组件 - 极简下划线风格 (Option C)
 * 仅保留底部强调线，去除背景和边框，强调内容本身的层次
 */
export function PageHeader({
  code,
  title,
  description,
  breadcrumbItems,
  className,
  action,
}: PageHeaderProps) {
  const displayCode = code ?? "PAGE"

  return (
    <div className={cn("px-4 lg:px-6 mb-2", className)}>
      {breadcrumbItems && breadcrumbItems.length > 0 ? (
        <div className="mb-2 text-xs text-muted-foreground flex items-center gap-2">
          {breadcrumbItems.map((item, index) => (
            <span key={item.href + item.label} className="flex items-center gap-2">
              {index > 0 ? <span>/</span> : null}
              <span>{item.label}</span>
            </span>
          ))}
        </div>
      ) : null}
      <div className="flex items-end gap-2 mb-2">
        <div className="flex items-baseline gap-3 border-b-2 border-primary pb-2">
          <h1 className="text-2xl font-bold tracking-tight uppercase">
            {title}
          </h1>
          <span className="font-mono text-xs text-muted-foreground font-medium tracking-wide">
            /{displayCode}
          </span>
        </div>
        <div className="flex-1 h-1.5 bg-[repeating-linear-gradient(45deg,transparent,transparent_4px,currentColor_4px,currentColor_5px)] text-primary/10" />
        {action}
      </div>
      {description && (
        <p className="text-sm text-muted-foreground">
          {description}
        </p>
      )}
    </div>
  )
}
