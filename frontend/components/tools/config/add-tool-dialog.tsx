"use client"

import React, { useState, useEffect } from "react"
import { Wrench, AlertTriangle } from "@/components/icons"
import { IconPlus } from "@/components/icons"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useTranslations } from "next-intl"

// Import UI components
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"
import { LoadingSpinner } from "@/components/loading-spinner"
import { IconX } from "@/components/icons"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"

// Import React Query Hook
import { useCreateTool, useUpdateTool } from "@/hooks/use-tools"

// Import type definitions
import type { Tool } from "@/types/tool.types"
import { CategoryNameMap } from "@/types/tool.types"

// Component props type definition
interface AddToolDialogProps {
  tool?: Tool                   // Tool data to edit (optional, edit mode when provided)
  onAdd?: (tool: Tool) => void  // Callback function on successful add (optional)
  open?: boolean                // External control for dialog open state
  onOpenChange?: (open: boolean) => void  // External control callback for dialog open state
}

/**
 * Auto-generate version query command based on tool name and install command
 */
function generateVersionCommand(toolName: string, installCommand: string): string {
  if (!toolName) return ""

  const lowerName = toolName.toLowerCase().trim()
  const lowerInstall = installCommand.toLowerCase()

  // Python tools
  if (lowerInstall.includes("python") || lowerInstall.includes(".py")) {
    return `python ${lowerName}.py -v`
  }

  // Go tools
  if (lowerInstall.includes("go install") || lowerInstall.includes("go get")) {
    return `${lowerName} -version`
  }

  // Default: try common version commands
  return `${lowerName} --version`
}

/**
 * Add tool dialog component (using React Query)
 */
export function AddToolDialog({
  tool,
  onAdd,
  open: externalOpen,
  onOpenChange: externalOnOpenChange
}: AddToolDialogProps) {
  const t = useTranslations("tools.config")
  
  // Determine if in edit mode or add mode
  const isEditMode = !!tool

  // Dialog open state - supports external control
  const [internalOpen, setInternalOpen] = useState(false)
  const open = externalOpen !== undefined ? externalOpen : internalOpen
  const setOpen = externalOnOpenChange || setInternalOpen

  // Use predefined category list
  const availableCategories = Object.keys(CategoryNameMap)

  // Use React Query create and update tool mutations
  const createTool = useCreateTool()
  const updateTool = useUpdateTool()

  // Form validation Schema
  const formSchema = z.object({
    name: z.string()
      .min(2, { message: t("toolNameMin") })
      .max(255, { message: t("toolNameMax") }),
    repoUrl: z.string().optional().or(z.literal("")),
    version: z.string().max(100).optional().or(z.literal("")),
    description: z.string().max(1000).optional().or(z.literal("")),
    categoryNames: z.array(z.string()),
    installCommand: z.string().min(1, { message: t("installCommandRequired") }),
    updateCommand: z.string().min(1, { message: t("updateCommandRequired") }),
    versionCommand: z.string().min(1, { message: t("versionCommandRequired") }),
  })

  type FormValues = z.infer<typeof formSchema>

  // Initialize form
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: tool?.name || "",
      repoUrl: tool?.repoUrl || "",
      version: tool?.version || "",
      description: tool?.description || "",
      categoryNames: tool?.categoryNames || [],
      installCommand: tool?.installCommand || "",
      updateCommand: tool?.updateCommand || "",
      versionCommand: tool?.versionCommand || "",
    },
  })

  // Reset form when tool changes
  useEffect(() => {
    if (tool) {
      form.reset({
        name: tool.name || "",
        repoUrl: tool.repoUrl || "",
        version: tool.version || "",
        description: tool.description || "",
        categoryNames: tool.categoryNames || [],
        installCommand: tool.installCommand || "",
        updateCommand: tool.updateCommand || "",
        versionCommand: tool.versionCommand || "",
      })
    }
  }, [tool, form])

  // Watch form value changes
  const watchName = form.watch("name")
  const watchInstallCommand = form.watch("installCommand")
  const watchVersionCommand = form.watch("versionCommand")
  const watchCategoryNames = form.watch("categoryNames")

  // Auto-generate version command
  useEffect(() => {
    if (watchName && watchInstallCommand && !watchVersionCommand) {
      const generatedCmd = generateVersionCommand(watchName, watchInstallCommand)
      form.setValue("versionCommand", generatedCmd)
    }
  }, [watchName, watchInstallCommand, watchVersionCommand, form])

  // Handle form submission
  const onSubmit = (values: FormValues) => {
    const toolData = {
      name: values.name.trim(),
      type: 'opensource' as const,
      repoUrl: values.repoUrl?.trim() || undefined,
      version: values.version?.trim() || undefined,
      description: values.description?.trim() || undefined,
      categoryNames: values.categoryNames.length > 0 ? values.categoryNames : undefined,
      installCommand: values.installCommand.trim(),
      updateCommand: values.updateCommand.trim(),
      versionCommand: values.versionCommand.trim(),
    }

    const onSuccessCallback = (response: { tool?: Tool }) => {
      form.reset()
      setOpen(false)
      if (onAdd && response?.tool) {
        onAdd(response.tool)
      }
    }

    if (isEditMode && tool?.id) {
      updateTool.mutate(
        { id: tool.id, data: toolData },
        { onSuccess: onSuccessCallback }
      )
    } else {
      createTool.mutate(toolData, { onSuccess: onSuccessCallback })
    }
  }

  // Handle category tag click
  const handleCategoryToggle = (categoryName: string) => {
    const current = form.getValues("categoryNames")
    const isSelected = current.includes(categoryName)
    form.setValue(
      "categoryNames",
      isSelected
        ? current.filter(c => c !== categoryName)
        : [...current, categoryName]
    )
  }

  // Remove category tag
  const handleCategoryRemove = (categoryName: string) => {
    const current = form.getValues("categoryNames")
    form.setValue("categoryNames", current.filter(c => c !== categoryName))
  }

  // Handle dialog close
  const handleOpenChange = (newOpen: boolean) => {
    if (!createTool.isPending && !updateTool.isPending) {
      setOpen(newOpen)
      if (!newOpen) {
        form.reset()
      }
    }
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      {externalOpen === undefined && (
        <DialogTrigger asChild>
          <Button>
            <IconPlus className="h-5 w-5" />
            {t("addTool")}
          </Button>
        </DialogTrigger>
      )}

      <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Wrench />
            <span>{isEditMode ? t("editTool") : t("addNewTool")}</span>
          </DialogTitle>
          <DialogDescription>
            {t("dialogDesc")}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="grid gap-6 py-4">
              <div className="space-y-4">
                <h3 className="text-sm font-semibold text-muted-foreground">{t("basicInfo")}</h3>

                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("toolName")} <span className="text-destructive">*</span></FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t("toolNamePlaceholder")}
                          disabled={createTool.isPending || updateTool.isPending}
                          maxLength={255}
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>{t("characters", { count: field.value.length, max: 255 })}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="repoUrl"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("repoUrl")}</FormLabel>
                      <FormControl>
                        <Input
                          type="url"
                          placeholder={t("repoUrlPlaceholder")}
                          disabled={createTool.isPending || updateTool.isPending}
                          maxLength={512}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="version"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("currentVersion")}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t("versionPlaceholder")}
                          disabled={createTool.isPending || updateTool.isPending}
                          maxLength={100}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("toolDesc")}</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={t("toolDescPlaceholder")}
                          disabled={createTool.isPending || updateTool.isPending}
                          rows={3}
                          maxLength={1000}
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>{t("characters", { count: (field.value || "").length, max: 1000 })}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid gap-2">
                  <FormLabel>{t("categoryTags")}</FormLabel>

                  {watchCategoryNames.length > 0 && (
                    <div className="flex flex-wrap gap-2 p-3 border rounded-md bg-muted/50">
                      {watchCategoryNames.map((categoryName) => (
                        <Badge
                          key={categoryName}
                          variant="default"
                          className="flex items-center gap-1 px-2 py-1"
                        >
                          {CategoryNameMap[categoryName] || categoryName}
                          <button
                            type="button"
                            onClick={() => handleCategoryRemove(categoryName)}
                            disabled={createTool.isPending || updateTool.isPending}
                            className="ml-1 hover:bg-primary/20 rounded-full p-0.5"
                          >
                            <IconX className="h-3 w-3" />
                          </button>
                        </Badge>
                      ))}
                    </div>
                  )}

                  <div className="flex flex-wrap gap-2 p-3 border rounded-md">
                    {availableCategories.length > 0 ? (
                      availableCategories.map((categoryName) => {
                        const isSelected = watchCategoryNames.includes(categoryName)
                        return (
                          <Badge
                            key={categoryName}
                            variant={isSelected ? "secondary" : "outline"}
                            className="cursor-pointer hover:bg-secondary/80 transition-colors"
                            onClick={() => handleCategoryToggle(categoryName)}
                          >
                            {CategoryNameMap[categoryName] || categoryName}
                          </Badge>
                        )
                      })
                    ) : (
                      <p className="text-sm text-muted-foreground">{t("noCategories")}</p>
                    )}
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <h3 className="text-sm font-semibold text-muted-foreground">{t("commandConfig")}</h3>

                <FormField
                  control={form.control}
                  name="installCommand"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("installCommand")} <span className="text-destructive">*</span></FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={t("installCommandPlaceholder")}
                          disabled={createTool.isPending || updateTool.isPending}
                          rows={3}
                          className="font-mono text-sm"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription className="space-y-1">
                        <span className="block"><strong>{t("installCommandHint")}</strong></span>
                        <span className="block">• {t("installCommandGit")} <code className="bg-muted px-1 py-0.5 rounded">git clone https://github.com/user/tool</code></span>
                        <span className="block">• {t("installCommandGo")} <code className="bg-muted px-1 py-0.5 rounded">go install -v github.com/tool@latest</code></span>
                        <span className="flex items-center gap-1 text-amber-600">
                          <AlertTriangle className="h-3.5 w-3.5" />
                          {t("installCommandNote")}
                        </span>
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="updateCommand"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("updateCommand")} <span className="text-destructive">*</span></FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={t("updateCommandPlaceholder")}
                          disabled={createTool.isPending || updateTool.isPending}
                          rows={2}
                          className="font-mono text-sm"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription className="space-y-1">
                        <span className="block">• {t("updateCommandGitHint")} <code className="bg-muted px-1 py-0.5 rounded">git pull</code></span>
                        <span className="block">• {t("updateCommandGoHint")}</span>
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="versionCommand"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t("versionCommand")} <span className="text-destructive">*</span>
                        {field.value && (
                          <span className="ml-2 text-xs text-muted-foreground font-normal">
                            {t("versionCommandAutoGenerated")}
                          </span>
                        )}
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t("versionCommandPlaceholder")}
                          disabled={createTool.isPending || updateTool.isPending}
                          maxLength={500}
                          className="font-mono text-sm"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription className="space-y-1">
                        <span className="block">{t("versionCommandHint")}</span>
                        <span className="block">• <code className="bg-muted px-1 py-0.5 rounded">toolname -v</code></span>
                        <span className="block">• <code className="bg-muted px-1 py-0.5 rounded">toolname -V</code></span>
                        <span className="block">• <code className="bg-muted px-1 py-0.5 rounded">toolname --version</code></span>
                        <span className="block">• <code className="bg-muted px-1 py-0.5 rounded">python tool_name.py -v</code></span>
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => handleOpenChange(false)}
                disabled={createTool.isPending || updateTool.isPending}
              >
                {t("cancel")}
              </Button>
              <Button
                type="submit"
                disabled={createTool.isPending || updateTool.isPending || !form.formState.isValid}
              >
                {(createTool.isPending || updateTool.isPending) ? (
                  <>
                    <LoadingSpinner />
                    {isEditMode ? t("saving") : t("creating")}
                  </>
                ) : (
                  <>
                    <IconPlus className="h-5 w-5" />
                    {isEditMode ? t("saveChanges") : t("createTool")}
                  </>
                )}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
