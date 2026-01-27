"use client"

import { useMemo } from "react"
import { cn } from "@/lib/utils"
import { Progress } from "@/components/ui/progress"
import { IconAlertTriangle } from "@tabler/icons-react"

interface MetricProgressProps {
  label: string
  value: number
  threshold?: number
  unit?: string
  className?: string
  showWarning?: boolean
}

export function MetricProgress({
  label,
  value,
  threshold,
  unit = "%",
  className,
  showWarning = true,
}: MetricProgressProps) {
  const percentage = Math.min(100, Math.max(0, value))

  const status = useMemo(() => {
    if (!threshold) return "normal"
    if (percentage >= threshold) return "critical"
    if (percentage >= threshold * 0.8) return "warning"
    return "normal"
  }, [percentage, threshold])

  const progressColor = useMemo(() => {
    if (status === "critical") return "bg-red-500"
    if (status === "warning") return "bg-amber-500"
    return "bg-emerald-500"
  }, [status])

  const textColor = useMemo(() => {
    if (status === "critical") return "text-red-600"
    if (status === "warning") return "text-amber-600"
    return "text-foreground"
  }, [status])

  return (
    <div className={cn("space-y-1.5", className)}>
      <div className="flex items-center justify-between text-xs">
        <span className="text-muted-foreground flex items-center gap-1">
          {label}
          {showWarning && status !== "normal" && (
            <IconAlertTriangle className="h-3 w-3 text-amber-500" />
          )}
        </span>
        <span className={cn("font-medium tabular-nums", textColor)}>
          {percentage.toFixed(0)}{unit}
        </span>
      </div>
      <Progress
        value={percentage}
        className={cn(
          "h-1.5",
          status === "critical" && "bg-red-500/20",
          status === "warning" && "bg-amber-500/20"
        )}
      >
        <div
          className={cn("h-full transition-all duration-300", progressColor)}
          style={{ width: `${percentage}%` }}
        />
      </Progress>
      {threshold && (
        <div className="text-[10px] text-muted-foreground">
          阈值: {threshold}{unit}
        </div>
      )}
    </div>
  )
}
