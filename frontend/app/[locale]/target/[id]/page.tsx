import { redirect } from "next/navigation"

/**
 * Target detail default page
 * Automatically redirects to overview page
 */
export default function TargetDetailPage({
  params,
}: {
  params: { id: string }
}) {
  redirect(`/target/${params.id}/overview/`)
}

