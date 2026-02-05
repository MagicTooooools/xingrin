"use client"

import React from "react"
import { useParams } from "next/navigation"
import { SubdomainsDetailView } from "@/components/subdomains/subdomains-detail-view"
export default function ScanHistorySubdomainPage() {
  const { id } = useParams<{ id: string }>()

  return (
    <div className="px-4 lg:px-6">
      <SubdomainsDetailView scanId={parseInt(id)} />
    </div>
  )
}
