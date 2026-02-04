"use client"

import { AgentList, ArchitectureDialog } from "@/components/settings/workers"
import { useTranslations } from "next-intl"
import { PageHeader } from "@/components/common/page-header"

export default function WorkersPage() {
  const t = useTranslations("pages.workers")

  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <div className="flex items-center justify-between">
        <PageHeader
          code="WRK-01"
          title={t("title")}
          description={t("description")}
          className="px-0"
        />
        <ArchitectureDialog />
      </div>
      <AgentList />
    </div>
  )
}
