"use client"

import React, { useState } from "react"
import { Wrench } from "@/components/icons"
import { IconPlus } from "@/components/icons"
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
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"
import { LoadingSpinner } from "@/components/loading-spinner"
import { IconX } from "@/components/icons"
import { CategoryNameMap, type Tool } from "@/types/tool.types"
import { useCreateTool, useUpdateTool } from "@/hooks/use-tools"

// Component props type definition
interface AddCustomToolDialogProps {
  tool?: Tool                    // Tool data to edit (optional, edit mode when provided)
  onAdd?: (tool: Tool) => void   // Callback function on successful add (optional)
  open?: boolean                 // External control for dialog open state
  onOpenChange?: (open: boolean) => void  // External control callback for dialog open state
}

/**
 * Add/Edit custom tool dialog component
 */
export function AddCustomToolDialog({ 
  tool,
  onAdd, 
  open: externalOpen, 
  onOpenChange: externalOnOpenChange 
}: AddCustomToolDialogProps) {
  const t = useTranslations("tools.config")
  
  // Determine if in edit mode or add mode
  const isEditMode = !!tool
  
  // Dialog open state - supports external control
  const [internalOpen, setInternalOpen] = useState(false)
  const open = externalOpen !== undefined ? externalOpen : internalOpen
  const setOpen = externalOnOpenChange || setInternalOpen
  
  // Form data state - initialize with tool data if in edit mode
  const [formData, setFormData] = useState({
    name: tool?.name || "",
    description: tool?.description || "",
    directory: tool?.directory || "",
    categoryNames: tool?.categoryNames || [] as string[],
  })

  // Use predefined category list
  const availableCategories = Object.keys(CategoryNameMap)

  // Use React Query create and update tool mutations
  const createTool = useCreateTool()
  const updateTool = useUpdateTool()

  // Update form data when tool changes
  React.useEffect(() => {
    if (tool) {
      setFormData({
        name: tool.name || "",
        description: tool.description || "",
        directory: tool.directory || "",
        categoryNames: tool.categoryNames || [],
      })
    }
  }, [tool])

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    // Form validation
    if (!formData.name.trim() || !formData.directory.trim()) {
      return
    }

    const toolData = {
      name: formData.name.trim(),
      type: 'custom' as const, // Custom tool
      description: formData.description.trim() || undefined,
      directory: formData.directory.trim(),
      categoryNames: formData.categoryNames.length > 0 ? formData.categoryNames : undefined,
    }

    const onSuccessCallback = (response: { tool?: Tool }) => {
      // Reset form
      setFormData({
        name: "",
        description: "",
        directory: "",
        categoryNames: [],
      })
      
      // Close dialog
      setOpen(false)
      
      // Call external callback (if provided)
      if (onAdd && response?.tool) {
        onAdd(response.tool)
      }
    }

    // Choose create or update based on mode
    if (isEditMode && tool?.id) {
      // Edit mode: call update API
      updateTool.mutate(
        { id: tool.id, data: toolData },
        { onSuccess: onSuccessCallback }
      )
    } else {
      // Create mode: call create API
      createTool.mutate(toolData, { onSuccess: onSuccessCallback })
    }
  }

  // Handle dialog close - reset form
  const handleOpenChange = (newOpen: boolean) => {
    // Don't allow closing while submitting
    if (!createTool.isPending && !updateTool.isPending) {
      setOpen(newOpen)
      if (!newOpen) {
        // Reset form when dialog closes
        setFormData({
          name: "",
          description: "",
          directory: "",
          categoryNames: [],
        })
      }
    }
  }

  // Handle category tag toggle
  const handleCategoryToggle = (categoryName: string) => {
    setFormData((prev) => {
      const isSelected = prev.categoryNames.includes(categoryName)
      return {
        ...prev,
        categoryNames: isSelected
          ? prev.categoryNames.filter(c => c !== categoryName)
          : [...prev.categoryNames, categoryName]
      }
    })
  }

  // Remove category tag
  const handleCategoryRemove = (categoryName: string) => {
    setFormData((prev) => ({
      ...prev,
      categoryNames: prev.categoryNames.filter(c => c !== categoryName)
    }))
  }

  // Form validation - check required fields
  const isFormValid = 
    formData.name.trim().length > 0 &&
    formData.directory.trim().length > 0

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      {/* Trigger button - only show when not externally controlled */}
      {externalOpen === undefined && (
        <DialogTrigger asChild>
          <Button>
            <IconPlus className="h-5 w-5" />
            {t("addTool")}
          </Button>
        </DialogTrigger>
      )}
      
      {/* Dialog content */}
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Wrench />
            <span>{isEditMode ? t("editCustomTool") : t("addCustomTool")}</span>
          </DialogTitle>
          <DialogDescription>
            {t("customDialogDesc")}
          </DialogDescription>
        </DialogHeader>
        
        {/* Form */}
        <form onSubmit={handleSubmit}>
          <div className="grid gap-6 py-4">
            {/* Tool name */}
            <div className="grid gap-2">
              <Label htmlFor="name">
                {t("toolName")} <span className="text-red-500">*</span>
              </Label>
              <Input
                id="name"
                placeholder={t("customToolNamePlaceholder")}
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                disabled={createTool.isPending || updateTool.isPending}
                required
              />
            </div>

            {/* Tool description */}
            <div className="grid gap-2">
              <Label htmlFor="description">{t("toolDesc")}</Label>
              <Textarea
                id="description"
                placeholder={t("customToolDescPlaceholder")}
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                disabled={createTool.isPending || updateTool.isPending}
                rows={3}
              />
            </div>

            {/* Tool path */}
            <div className="grid gap-2">
              <Label htmlFor="directory">
                {t("toolPath")} <span className="text-red-500">*</span>
              </Label>
              <Input
                id="directory"
                placeholder={t("toolPathPlaceholder")}
                value={formData.directory}
                onChange={(e) => setFormData({ ...formData, directory: e.target.value })}
                disabled={createTool.isPending || updateTool.isPending}
                required
              />
              <p className="text-xs text-muted-foreground">
                {t("toolPathHint")}
              </p>
            </div>

            {/* Category tags */}
            <div className="grid gap-2">
              <Label>{t("categoryTags")}</Label>
              
              {/* Selected tags */}
              {formData.categoryNames.length > 0 && (
                <div className="flex flex-wrap gap-2 p-3 border rounded-md bg-muted/50">
                  {formData.categoryNames.map((categoryName) => (
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

              {/* Available tags */}
              <div className="flex flex-wrap gap-2 p-3 border rounded-md">
                {availableCategories.length > 0 ? (
                  availableCategories.map((categoryName) => {
                    const isSelected = formData.categoryNames.includes(categoryName)
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
              disabled={createTool.isPending || updateTool.isPending || !isFormValid}
            >
              {(createTool.isPending || updateTool.isPending) ? (
                <>
                  <LoadingSpinner/>
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
      </DialogContent>
    </Dialog>
  )
}
