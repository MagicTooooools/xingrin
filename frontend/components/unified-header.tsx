"use client"

import { Button } from "@/components/ui/button"
import { SidebarTrigger } from "@/components/ui/sidebar"
import { NotificationDrawer } from "@/components/notifications"
import { QuickScanDialog } from "@/components/scan/quick-scan-dialog"
import { LanguageSwitcher } from "@/components/language-switcher"
import { Link } from "@/i18n/navigation"
import { useTranslations } from "next-intl"

/**
 * 统一顶栏组件
 * 包含 Logo、侧边栏触发器、快捷操作按钮
 * 横跨整个页面宽度，Logo 在最左侧
 */
export function UnifiedHeader() {
  const t = useTranslations("navigation")
  const logoSrc = "/images/icon-64.png"

  return (
    <header
      data-slot="unified-header"
      className="flex h-(--header-height) shrink-0 items-center border-b bg-background"
    >
      {/* Logo 区域 - 固定宽度，与侧边栏宽度一致 */}
      <div className="flex h-full w-(--sidebar-width) shrink-0 items-center px-4">
        <Link href="/" className="flex items-center gap-2">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img src={logoSrc} alt="Logo" className="size-6" />
          <span className="text-base font-semibold">
            {t("appName")}
          </span>
        </Link>
      </div>

      {/* 右侧内容区 */}
      <div className="flex flex-1 items-center gap-2 px-4">
        {/* 侧边栏触发器 - 仅移动端显示 */}
        <SidebarTrigger className="-ml-1 md:hidden" />

        {/* 右侧按钮区 */}
        <div className="ml-auto flex items-center gap-2">
          <QuickScanDialog />
          <NotificationDrawer />
          <LanguageSwitcher />
          <Button variant="ghost" asChild size="sm" className="hidden sm:flex">
            <a
              href="https://github.com/yyhuni/xingrin"
              rel="noopener noreferrer"
              target="_blank"
            >
              GitHub
            </a>
          </Button>
        </div>
      </div>
    </header>
  )
}
