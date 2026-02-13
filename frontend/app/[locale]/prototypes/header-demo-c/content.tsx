"use client"

import { AllTargetsDetailView } from "@/components/target/all-targets-detail-view"
import { Button } from "@/components/ui/button"
import { Settings, Search } from "@/components/icons"

/**
 * Demo C：对接式界面
 * 设计理念：Header 底边框与表格顶边框合为一条线，形成无缝对接
 * 关键 CSS：Header border-b-2，表格 border-t-0，紧贴无间隙
 */
export default function DemoPageC() {
  return (
    <div className="flex flex-col h-full">
      {/* 说明区域 */}
      <div className="p-6 border-b bg-muted/30">
        <h1 className="text-xl font-bold">方案 C：对接式界面</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Header 与表格无缝对接，像一个控制台面板。
        </p>
      </div>

      {/* 核心 Demo 区域 */}
      <div className="flex-1 p-6 flex flex-col min-h-0">
        
        {/* Header 区域 - 底部粗边框，作为分割线 */}
        <div className="flex items-end justify-between border-b-2 border-primary pb-3 shrink-0 px-2">
          <div className="flex flex-col">
            <span className="text-[10px] font-mono text-muted-foreground mb-1 uppercase tracking-wider">
              Context: TGT-01
            </span>
            <div className="flex items-center gap-3">
              <h1 className="text-2xl font-bold leading-none">目标管理</h1>
              <span className="bg-[var(--success)]/10 text-[var(--success)] px-1.5 py-0.5 text-[10px] font-bold border border-[var(--success)]/20">
                ACTIVE
              </span>
            </div>
          </div>
          
          {/* Tab 风格的按钮 - 贴底对齐，制造对接效果 */}
          <div className="flex gap-0.5 translate-y-[2px]">
             <Button 
               variant="secondary" 
               size="sm" 
               className="rounded-none rounded-t border-x border-t border-b-0 border-border bg-muted text-muted-foreground h-8"
             >
               <Settings className="h-3.5 w-3.5 mr-1"/> Config
             </Button>
             <Button 
               variant="default" 
               size="sm" 
               className="rounded-none rounded-t border-x border-t border-b-2 border-b-primary shadow-none h-8"
             >
               <Search className="h-3.5 w-3.5 mr-1"/> List View
             </Button>
          </div>
        </div>
        
        {/* 表格区域 - 移除顶部边框，紧贴 Header */}
        <div className="flex-1 overflow-auto border-x-2 border-b-2 border-primary bg-card">
          <AllTargetsDetailView
            className="space-y-0 px-2 pb-2"
            tableClassName="border-0 rounded-none"
          />
        </div>
      </div>
    </div>
  )
}
