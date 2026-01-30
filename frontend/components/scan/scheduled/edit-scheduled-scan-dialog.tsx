"use client"

import React from "react"
import { toast } from "sonner"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import { Checkbox } from "@/components/ui/checkbox"
import { IconX, IconLoader2 } from "@tabler/icons-react"
import { cn } from "@/lib/utils"
import { useUpdateScheduledScan } from "@/hooks/use-scheduled-scans"
import { useTargets } from "@/hooks/use-targets"
import { useEngines } from "@/hooks/use-engines"
import { useTranslations } from "next-intl"
import type { ScheduledScan, UpdateScheduledScanRequest } from "@/types/scheduled-scan.types"
import type { ScanEngine } from "@/types/engine.types"
import type { Target } from "@/types/target.types"

interface EditScheduledScanDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  scheduledScan: ScheduledScan | null
  onSuccess?: () => void
}

export function EditScheduledScanDialog({
  open,
  onOpenChange,
  scheduledScan,
  onSuccess,
}: EditScheduledScanDialogProps) {
  const { mutate: updateScheduledScan, isPending } = useUpdateScheduledScan()
  const { data: targetsData } = useTargets()
  const { data: enginesData } = useEngines()
  const t = useTranslations("scan.scheduled")

  const CRON_PRESETS = [
    { label: t("presets.everyMinute"), value: "* * * * *" },
    { label: t("presets.every5Minutes"), value: "*/5 * * * *" },
    { label: t("presets.everyHour"), value: "0 * * * *" },
    { label: t("presets.daily2am"), value: "0 2 * * *" },
    { label: t("presets.daily4am"), value: "0 4 * * *" },
    { label: t("presets.weekly"), value: "0 2 * * 1" },
    { label: t("presets.monthly"), value: "0 2 1 * *" },
  ]

  const [name, setName] = React.useState("")
  const [engineIds, setEngineIds] = React.useState<number[]>([])
  const [selectedTargetId, setSelectedTargetId] = React.useState<number | null>(null)
  const [cronExpression, setCronExpression] = React.useState("")

  React.useEffect(() => {
    if (scheduledScan && open) {
      setName(scheduledScan.name)
      setEngineIds(scheduledScan.engineIds || [])
      setSelectedTargetId(scheduledScan.targetId || null)
      setCronExpression(scheduledScan.cronExpression || "0 2 * * *")
    }
  }, [scheduledScan, open])

  const handleEngineToggle = (engineId: number, checked: boolean) => {
    if (checked) {
      setEngineIds((prev) => [...prev, engineId])
    } else {
      setEngineIds((prev) => prev.filter((id) => id !== engineId))
    }
  }

  const handleTargetSelect = (targetId: number) => {
    setSelectedTargetId(selectedTargetId === targetId ? null : targetId)
  }

  const validateCron = (cron: string): boolean => {
    const parts = cron.trim().split(/\s+/)
    return parts.length === 5
  }

  const handleSubmit = () => {
    if (!scheduledScan) return

    if (!name.trim()) {
      toast.error(t("form.taskNameRequired"))
      return
    }
    if (engineIds.length === 0) {
      toast.error(t("form.scanEngineRequired"))
      return
    }
    if (scheduledScan.scanMode === 'target' && !selectedTargetId) {
      toast.error(t("toast.selectTarget"))
      return
    }
    if (!validateCron(cronExpression)) {
      toast.error(t("form.cronRequired"))
      return
    }

    const request: UpdateScheduledScanRequest = {
      name: name.trim(),
      engineIds: engineIds,
      cronExpression: cronExpression.trim(),
    }

    if (scheduledScan.scanMode === 'target' && selectedTargetId) {
      request.targetId = selectedTargetId
    }

    updateScheduledScan(
      { id: scheduledScan.id, data: request },
      {
        onSuccess: () => {
          onOpenChange(false)
          onSuccess?.()
        },
        onError: (err: unknown) => {
          const error = err as { response?: { data?: { error?: { code?: string; message?: string } } } }
          if (error?.response?.data?.error?.code === 'CONFIG_CONFLICT') {
            toast.error(t("toast.configConflict"), {
              description: error.response.data.error.message,
            })
          }
        },
      }
    )
  }

  const targets: Target[] = targetsData?.targets || []
  const engines: ScanEngine[] = enginesData || []

  if (!scheduledScan) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{t("editTitle")}</DialogTitle>
          <DialogDescription>{t("editDesc")}</DialogDescription>
        </DialogHeader>

        <div className="grid gap-6 py-4">
          <div className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="edit-name">{t("form.taskName")} *</Label>
              <Input
                id="edit-name"
                placeholder={t("form.taskNamePlaceholder")}
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
          </div>

          <div className="grid gap-2">
            <Label>{t("form.scanEngine")} *</Label>
            {engineIds.length > 0 && (
              <p className="text-xs text-muted-foreground">{t("form.selectedEngines", { count: engineIds.length })}</p>
            )}
            <div className="border rounded-md p-3 max-h-[150px] overflow-y-auto space-y-2">
              {engines.length === 0 ? (
                <p className="text-sm text-muted-foreground">{t("form.noEngine")}</p>
              ) : (
                engines.map((engine) => (
                  <label
                    key={engine.id}
                    htmlFor={`edit-engine-${engine.id}`}
                    className={cn(
                      "flex items-center gap-3 p-2 rounded-lg cursor-pointer transition-all",
                      engineIds.includes(engine.id)
                        ? "bg-primary/10 border border-primary/30"
                        : "hover:bg-muted/50 border border-transparent"
                    )}
                  >
                    <Checkbox
                      id={`edit-engine-${engine.id}`}
                      checked={engineIds.includes(engine.id)}
                      onCheckedChange={(checked) => handleEngineToggle(engine.id, checked as boolean)}
                    />
                    <span className="text-sm">{engine.name}</span>
                  </label>
                ))
              )}
            </div>
          </div>

          <div className="grid gap-2">
            <Label>{t("form.scanScope")}</Label>
            {scheduledScan.scanMode === 'organization' ? (
              <div className="border rounded-md p-3 bg-muted/50">
                <div className="flex items-center gap-2">
                  <Badge variant="secondary">{t("form.organizationMode")}</Badge>
                  <span className="font-medium">{scheduledScan.organizationName}</span>
                </div>
                <p className="text-xs text-muted-foreground mt-2">
                  {t("form.organizationModeHint")}
                </p>
              </div>
            ) : (
              <>
                <div className="border rounded-md p-3 max-h-[150px] overflow-y-auto">
                  {targets.length === 0 ? (
                    <p className="text-sm text-muted-foreground">{t("form.noAvailableTarget")}</p>
                  ) : (
                    <div className="flex flex-wrap gap-2">
                      {targets.map((target) => (
                        <Badge
                          key={target.id}
                          variant={selectedTargetId === target.id ? "default" : "outline"}
                          className="cursor-pointer"
                          onClick={() => handleTargetSelect(target.id)}
                        >
                          {target.name}
                          {selectedTargetId === target.id && (
                            <IconX className="h-3 w-3 ml-1" />
                          )}
                        </Badge>
                      ))}
                    </div>
                  )}
                </div>
                {selectedTargetId && (
                  <p className="text-xs text-muted-foreground">
                    {t("form.selected")}: {targets.find(t => t.id === selectedTargetId)?.name}
                  </p>
                )}
              </>
            )}
          </div>

          <div className="grid gap-4">
            <div className="grid gap-2">
              <Label>{t("form.cronExpression")} *</Label>
              <Input
                placeholder={t("form.cronPlaceholder")}
                value={cronExpression}
                onChange={(e) => setCronExpression(e.target.value)}
                className="font-mono"
              />
              <p className="text-xs text-muted-foreground">
                {t("form.cronFormat")}
              </p>
            </div>

            <div className="grid gap-2">
              <Label className="text-xs text-muted-foreground">{t("form.quickSelect")}</Label>
              <div className="flex flex-wrap gap-2">
                {CRON_PRESETS.map((preset) => (
                  <Badge
                    key={preset.value}
                    variant={cronExpression === preset.value ? "default" : "outline"}
                    className="cursor-pointer"
                    onClick={() => setCronExpression(preset.value)}
                  >
                    {preset.label}
                  </Badge>
                ))}
              </div>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            {t("buttons.cancel")}
          </Button>
          <Button onClick={handleSubmit} disabled={isPending}>
            {isPending && <IconLoader2 className="h-4 w-4 animate-spin" />}
            {t("buttons.saveChanges")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
