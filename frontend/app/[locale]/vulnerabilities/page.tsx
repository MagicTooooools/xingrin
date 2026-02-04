"use client"

import React from "react"
import { useTranslations } from "next-intl"
import { VulnerabilitiesDetailView } from "@/components/vulnerabilities"
import { PageHeader } from "@/components/common/page-header"

/**
 * All vulnerabilities page
 * Displays all vulnerabilities in the system
 */
export default function VulnerabilitiesPage() {
  const t = useTranslations("vulnerabilities")

  return (
    <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
      <PageHeader
        code="VUL-01"
        title={t("title")}
        description={t("description")}
      />

      {/* Vulnerability list */}
      <div className="px-4 lg:px-6">
        <VulnerabilitiesDetailView />
      </div>
    </div>
  )
}
