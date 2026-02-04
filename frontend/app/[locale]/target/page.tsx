"use client"

import { AllTargetsDetailView } from "@/components/target/all-targets-detail-view"
import { PageHeader } from "@/components/common/page-header"
import { useTranslations } from "next-intl"

export default function AllTargetsPage() {
  const t = useTranslations("pages.target")

  return (
    <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
      <PageHeader
        code="TGT-01"
        title={t("title")}
        description={t("description")}
      />

      {/* Target list */}
      <div className="px-4 lg:px-6">
        <AllTargetsDetailView />
      </div>
    </div>
  )
}
