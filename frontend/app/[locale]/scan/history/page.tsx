"use client"

import { useTranslations } from "next-intl"
import { ScanHistoryList } from "@/components/scan/history/scan-history-list"
import { ScanHistoryStatCards } from "@/components/scan/history/scan-history-stat-cards"
import { PageHeader } from "@/components/common/page-header"

/**
 * Scan history page
 * Displays historical records of all scan tasks
 */
export default function ScanHistoryPage() {
  const t = useTranslations("scan.history")

  return (
    <div className="@container/main flex flex-col gap-4 py-4 md:gap-6 md:py-6">
      <PageHeader
        code="SCN-01"
        title={t("title")}
        description={t("description")}
      />

      {/* Statistics cards */}
      <div className="px-4 lg:px-6">
        <ScanHistoryStatCards />
      </div>

      {/* Scan history list */}
      <div className="px-4 lg:px-6">
        <ScanHistoryList />
      </div>
    </div>
  )
}
