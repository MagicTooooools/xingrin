"use client"

import { AgentList, ArchitectureDialog } from "@/components/settings/workers"
import { useTranslations } from "next-intl"

export default function WorkersPage() {
  const t = useTranslations("pages.workers")

  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">{t("title")}</h1>
          <p className="text-muted-foreground">
            {t("description")}
          </p>
        </div>
        <ArchitectureDialog />
      </div>
      <AgentList />
    </div>
  )
}
