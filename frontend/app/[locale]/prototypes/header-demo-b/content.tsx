"use client"

import { AllTargetsDetailView } from "@/components/target/all-targets-detail-view"
import { Button } from "@/components/ui/button"

/**
 * Demo B：隐形标题
 * 设计理念：Header 完全无边框，只靠字号和间距建立层次；表格保留自己的边框
 * 关键 CSS：Header 无任何边框装饰
 */
export default function DemoPageB() {
  return (
    <div className="flex flex-col h-full">
      {/* 说明区域 */}
      <div className="p-6 border-b bg-muted/30">
        <h1 className="text-xl font-bold">方案 B：隐形标题</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Header 无边框悬浮，表格独立拥有边框。最透气的设计。
        </p>
      </div>

      {/* 核心 Demo 区域 */}
      <div className="flex-1 p-6 flex flex-col min-h-0">
        
        {/* Header 区域 - 纯文字，无任何边框 */}
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4 pb-4 shrink-0 px-2">
          <div className="space-y-1">
            <div className="flex items-center gap-3">
              <h1 className="text-3xl font-bold tracking-tight text-foreground">目标管理</h1>
              <span className="font-mono text-xs text-muted-foreground bg-muted px-2 py-1">
                /TGT-01
              </span>
            </div>
            <p className="text-muted-foreground text-sm">
              Manage and monitor all your scan targets here.
            </p>
          </div>
          <div className="flex gap-2">
             <Button variant="secondary" size="sm">Export</Button>
             <Button size="sm">New Scan</Button>
          </div>
        </div>

        {/* 表格区域 - 保留表格自带边框，但加粗顶部边框以强调 */}
        <div className="flex-1 overflow-auto">
          <AllTargetsDetailView
            className="space-y-3"
            tableClassName="border-2 border-primary rounded-none"
          />
        </div>
      </div>
    </div>
  )
}
