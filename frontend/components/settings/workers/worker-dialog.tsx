"use client"

import { useEffect } from "react"
import { useForm } from "react-hook-form"
import { useTranslations } from "next-intl"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { useUpdateAgentConfig } from "@/hooks/use-agents"
import type { Agent } from "@/types/agent.types"

const formSchema = z.object({
  maxTasks: z.coerce.number().int().min(1).max(20),
  cpuThreshold: z.coerce.number().int().min(1).max(100),
  memThreshold: z.coerce.number().int().min(1).max(100),
  diskThreshold: z.coerce.number().int().min(1).max(100),
})

type FormValues = z.infer<typeof formSchema>

interface AgentConfigDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  agent?: Agent | null
}

export function AgentConfigDialog({ open, onOpenChange, agent }: AgentConfigDialogProps) {
  const t = useTranslations("settings.workers")
  const tCommon = useTranslations("common.actions")
  const updateAgentConfig = useUpdateAgentConfig()

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema) as never,
    defaultValues: {
      maxTasks: 5,
      cpuThreshold: 85,
      memThreshold: 85,
      diskThreshold: 90,
    },
  })

  useEffect(() => {
    if (!open) return
    form.reset({
      maxTasks: agent?.maxTasks ?? 5,
      cpuThreshold: agent?.cpuThreshold ?? 85,
      memThreshold: agent?.memThreshold ?? 85,
      diskThreshold: agent?.diskThreshold ?? 90,
    })
  }, [open, agent, form])

  const onSubmit = async (values: FormValues) => {
    if (!agent) return
    try {
      await updateAgentConfig.mutateAsync({
        id: agent.id,
        data: values,
      })
      onOpenChange(false)
    } catch {
      // handled by hook
    }
  }

  const isPending = updateAgentConfig.isPending

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[420px]">
        <DialogHeader>
          <DialogTitle>{t("config.title")}</DialogTitle>
          <DialogDescription>{t("config.desc")}</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="maxTasks"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("config.maxTasks")}</FormLabel>
                  <FormControl>
                    <Input type="number" min={1} max={20} {...field} />
                  </FormControl>
                  <FormDescription>{t("config.maxTasksDesc")}</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="grid grid-cols-3 gap-3">
              <FormField
                control={form.control}
                name="cpuThreshold"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("config.cpuThreshold")}</FormLabel>
                    <FormControl>
                      <Input type="number" min={1} max={100} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="memThreshold"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("config.memThreshold")}</FormLabel>
                    <FormControl>
                      <Input type="number" min={1} max={100} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="diskThreshold"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("config.diskThreshold")}</FormLabel>
                    <FormControl>
                      <Input type="number" min={1} max={100} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isPending}
              >
                {tCommon("cancel")}
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending ? t("config.saving") : t("config.save")}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
