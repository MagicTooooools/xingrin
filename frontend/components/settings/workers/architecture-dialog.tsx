"use client"

import { useTranslations } from "next-intl"
import { IconInfoCircle } from "@/components/icons"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { ScrollArea } from "@/components/ui/scroll-area"
import { ArchitectureFlow } from "./architecture-flow"

export function ArchitectureDialog() {
  const t = useTranslations("pages.workers")
  const labels = {
    location: t("flowTableLocation"),
    comms: t("flowTableComms"),
    responsibilities: t("flowTableResponsibilities"),
  }
  const roleDetails = [
    {
      id: "server",
      title: t("flowServerTitle"),
      location: t("flowServerLocation"),
      comms: t("flowServerComms"),
      responsibilities: [
        t("flowServerItem1"),
        t("flowServerItem2"),
        t("flowServerItem3"),
      ],
    },
    {
      id: "agent",
      title: t("flowAgentTitle"),
      location: t("flowAgentLocation"),
      comms: t("flowAgentComms"),
      responsibilities: [
        t("flowAgentItem1"),
        t("flowAgentItem2"),
        t("flowAgentItem3"),
      ],
    },
    {
      id: "worker",
      title: t("flowWorkerTitle"),
      location: t("flowWorkerLocation"),
      comms: t("flowWorkerComms"),
      responsibilities: [
        t("flowWorkerItem1"),
        t("flowWorkerItem2"),
        t("flowWorkerItem3"),
      ],
    },
  ]

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          <IconInfoCircle className="h-4 w-4 mr-2" />
          {t("viewArchitecture")}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-5xl max-h-[90vh]">
        <DialogHeader>
          <DialogTitle>{t("flowTitle")}</DialogTitle>
          <DialogDescription>{t("flowDesc")}</DialogDescription>
        </DialogHeader>
        <ScrollArea className="max-h-[calc(90vh-120px)] pr-4">
          <div className="space-y-6">
            {/* 流程图 */}
            <div className="space-y-2">
              <div>
                <p className="text-sm font-medium">{t("flowDiagramTitle")}</p>
                <p className="text-xs text-muted-foreground">
                  {t("flowDiagramDesc")}
                </p>
              </div>
              <ArchitectureFlow />
            </div>

            <Separator />

            {/* 角色详情 */}
            <div className="space-y-3">
              <div>
                <p className="text-sm font-medium">{t("flowRolesTitle")}</p>
                <p className="text-xs text-muted-foreground">
                  {t("flowRolesDesc")}
                </p>
              </div>
              <div className="grid gap-3 md:grid-cols-3">
                {roleDetails.map((role) => (
                  <div
                    key={role.id}
                    className="rounded-md border bg-muted/10 p-3"
                  >
                    <p className="text-sm font-medium">{role.title}</p>
                    <div className="mt-2 space-y-1.5 text-xs text-muted-foreground">
                      <p>
                        <span className="font-medium text-foreground">
                          {labels.location}:
                        </span>{" "}
                        {role.location}
                      </p>
                      <p>
                        <span className="font-medium text-foreground">
                          {labels.comms}:
                        </span>{" "}
                        {role.comms}
                      </p>
                      <div>
                        <span className="font-medium text-foreground">
                          {labels.responsibilities}:
                        </span>
                        <ul className="mt-1 list-disc space-y-0.5 pl-4">
                          {role.responsibilities.map((item, index) => (
                            <li key={index}>{item}</li>
                          ))}
                        </ul>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <Separator />

            {/* 工作流程步骤 */}
            <div className="space-y-3">
              <div>
                <p className="text-sm font-medium">{t("flowStepsTitle")}</p>
                <p className="text-xs text-muted-foreground">
                  {t("flowStepsDesc")}
                </p>
              </div>
              <ol className="space-y-2 text-sm text-muted-foreground">
                <li className="flex gap-2">
                  <span className="font-semibold text-foreground">1.</span>
                  <span>{t("flowStep1")}</span>
                </li>
                <li className="flex gap-2">
                  <span className="font-semibold text-foreground">2.</span>
                  <span>{t("flowStep2")}</span>
                </li>
                <li className="flex gap-2">
                  <span className="font-semibold text-foreground">3.</span>
                  <span>{t("flowStep3")}</span>
                </li>
              </ol>
            </div>
          </div>
        </ScrollArea>
      </DialogContent>
    </Dialog>
  )
}
