"use client"

import React, { useState, useRef } from "react"
import { Plus, Target as TargetIcon, Building2, Loader2, Check, ChevronsUpDown } from "lucide-react"
import { IconChevronLeft, IconChevronRight, IconChevronsLeft, IconChevronsRight } from "@tabler/icons-react"
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
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { cn } from "@/lib/utils"
import { LoadingSpinner } from "@/components/loading-spinner"
import { TargetValidator } from "@/lib/target-validator"

// Import React Query Hooks
import { useOrganizations } from "@/hooks/use-organizations"
import { useBatchCreateTargets } from "@/hooks/use-targets"
import type { BatchCreateTargetsRequest } from "@/types/target.types"

// Component props type definition
interface AddTargetDialogProps {
  onAdd?: () => void                                             // Success callback after adding
  open?: boolean                                                 // External control for dialog open state
  onOpenChange?: (open: boolean) => void                         // External control for dialog open callback
  prefetchEnabled?: boolean                                      // Whether to prefetch organization list
}

/**
 * Add target dialog component (supports organization selection)
 * 
 * Features:
 * 1. Batch input targets
 * 2. Optional organization selection
 * 3. Auto-create non-existent targets
 * 4. Auto-manage submission state
 * 5. Auto error handling and success notifications
 */
export function AddTargetDialog({ 
  onAdd,
  open: externalOpen, 
  onOpenChange: externalOnOpenChange,
  prefetchEnabled,
}: AddTargetDialogProps) {
  const t = useTranslations("target.dialog")
  const tCommon = useTranslations("common.actions")
  const tPagination = useTranslations("common.pagination")
  
  // Dialog open state - supports external control
  const [internalOpen, setInternalOpen] = useState(false)
  const open = externalOpen !== undefined ? externalOpen : internalOpen
  const setOpen = externalOnOpenChange || setInternalOpen
  const [orgPickerOpen, setOrgPickerOpen] = useState(false)
  
  // Form data state
  const [formData, setFormData] = useState({
    targets: "",  // Target list, one per line
    organizationId: "",  // Selected organization ID
  })
  
  // Organization picker state
  const [orgSearchQuery, setOrgSearchQuery] = useState("")
  const [orgPage, setOrgPage] = useState(1)
  const [orgPageSize, setOrgPageSize] = useState(20)  // Default 20 items per page
  const pageSizeOptions = [20, 50, 200, 500, 1000]
  
  // Validation error state
  const [invalidTargets, setInvalidTargets] = useState<Array<{ index: number; originalTarget: string; error: string; type?: string }>>([])
  
  // Use batch create targets mutation
  const batchCreateTargets = useBatchCreateTargets()
  
  // Refs for line numbers and textarea (for synchronized scrolling)
  const lineNumbersRef = useRef<HTMLDivElement | null>(null)
  const textareaRef = useRef<HTMLTextAreaElement | null>(null)
  
  // Get organization list (supports pagination)
  const shouldEnableOrgsQuery = Boolean(prefetchEnabled || orgPickerOpen)
  const { data: organizationsData, isLoading: isLoadingOrganizations } = useOrganizations(
    {
      page: orgPage,
      pageSize: orgPageSize,  // Dynamic page size
    },
    { enabled: shouldEnableOrgsQuery }
  )

  // Handle input change
  const handleInputChange = (field: keyof typeof formData, value: string) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }))

    if (field === "targets") {
      const lines = value
        .split("\n")
        .map((s) => s.trim())
        .filter((s) => s.length > 0)

      if (lines.length === 0) {
        setInvalidTargets([])
        return
      }

      const results = TargetValidator.validateTargetBatch(lines)
      const invalid = results
        .filter((r) => !r.isValid)
        .map((r) => ({ index: r.index, originalTarget: r.originalTarget, error: r.error || t("invalidFormat"), type: r.type }))
      setInvalidTargets(invalid)
    }
  }
  
  // Calculate target count
  const targetCount = formData.targets
    .split("\n")
    .map(line => line.trim())
    .filter(line => line.length > 0).length

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    // Form validation
    if (!formData.targets.trim()) {
      return
    }

    if (invalidTargets.length > 0) {
      return
    }

    // Parse target list (one target per line)
    const targetList = formData.targets
      .split("\n")
      .map(line => line.trim())
      .filter(line => line.length > 0)
      .map(name => ({
        name,
      }))

    if (targetList.length === 0) {
      return
    }

    // Assemble request data (organization is optional)
    const payload: BatchCreateTargetsRequest = {
      targets: targetList,
    }

    if (formData.organizationId) {
      payload.organizationId = parseInt(formData.organizationId, 10)
    }

    // Call batch create API
    batchCreateTargets.mutate(
      payload,
      {
        onSuccess: () => {
          // Reset form
          setFormData({
            targets: "",
            organizationId: "",
          })
          setInvalidTargets([])
          setOrgSearchQuery("")
          setOrgPage(1)
          setOrgPageSize(20)
          
          // Close dialog
          setOpen(false)
          
          // Call external callback (if provided)
          if (onAdd) {
            onAdd()
          }
        }
      }
    )
  }

  // Handle dialog close
  const handleOpenChange = (newOpen: boolean) => {
    if (!batchCreateTargets.isPending) {
      setOpen(newOpen)
      if (!newOpen) {
        // Reset form when closing
        setFormData({
          targets: "",
          organizationId: "",
        })
        setInvalidTargets([])
        setOrgSearchQuery("")
        setOrgPage(1)
        setOrgPageSize(20)  // Reset to default value
      }
    }
  }

  // Form validation
  const isFormValid = formData.targets.trim().length > 0 && invalidTargets.length === 0
  
  // Synchronize textarea and line numbers scrolling
  const handleTextareaScroll = (e: React.UIEvent<HTMLTextAreaElement>) => {
    if (lineNumbersRef.current) {
      lineNumbersRef.current.scrollTop = e.currentTarget.scrollTop
    }
  }

  // Get selected organization name
  const [selectedOrgName, setSelectedOrgName] = useState("")
  const selectedOrganization = organizationsData?.organizations.find(
    org => org.id.toString() === formData.organizationId
  )
  
  // Update selected organization name
  React.useEffect(() => {
    if (selectedOrganization) {
      setSelectedOrgName(selectedOrganization.name)
    }
  }, [selectedOrganization])
  
  // Filter organization list
  const filteredOrganizations = React.useMemo(() => {
    if (!organizationsData?.organizations) return []
    if (!orgSearchQuery) return organizationsData.organizations
    return organizationsData.organizations.filter(org => 
      org.name.toLowerCase().includes(orgSearchQuery.toLowerCase())
    )
  }, [organizationsData?.organizations, orgSearchQuery])
  
  // Handle organization selection
  const handleSelectOrganization = (orgId: string, orgName: string) => {
    handleInputChange("organizationId", orgId)
    setSelectedOrgName(orgName)
    setOrgPickerOpen(false)
    setOrgSearchQuery("")
    setOrgPage(1)
    setOrgPageSize(20)  // Reset to default value
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      {/* Trigger button - only shown when not externally controlled */}
      {externalOpen === undefined && (
        <DialogTrigger asChild>
          <Button size="sm">
            <Plus />
            {t("addTarget")}
          </Button>
        </DialogTrigger>
      )}
      
      {/* Dialog content */}
      <DialogContent className="sm:max-w-[650px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <TargetIcon />
            <span>{t("addTitle")}</span>
          </DialogTitle>
          <DialogDescription>
            {t("addDesc")}
          </DialogDescription>
        </DialogHeader>
        
        {/* Form */}
        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            {/* Target input (supports multiple lines) */}
            <div className="grid gap-2">
              <Label htmlFor="targets">
                {t("targetList")} <span className="text-destructive">*</span>
              </Label>
              <div className="flex border rounded-md overflow-hidden h-[180px]">
                {/* Line numbers column - fixed width */}
                <div className="flex-shrink-0 w-12 border-r bg-muted/50">
                  <div 
                    ref={lineNumbersRef}
                    className="py-3 px-2 text-right font-mono text-xs text-muted-foreground leading-[1.4] h-full overflow-y-auto scrollbar-hide"
                  >
                    {Array.from({ length: Math.max(formData.targets.split('\n').length, 8) }, (_, i) => (
                      <div key={i + 1} className="h-[20px]">
                        {i + 1}
                      </div>
                    ))}
                  </div>
                </div>
                {/* Input area - takes remaining space */}
                <div className="flex-1 overflow-hidden">
                  {/* Input - fixed height showing 8 lines */}
                  <Textarea
                    ref={textareaRef}
                    id="targets"
                    value={formData.targets}
                    onChange={(e) => handleInputChange("targets", e.target.value)}
                    onScroll={handleTextareaScroll}
                    placeholder={t("targetPlaceholder")}
                    disabled={batchCreateTargets.isPending}
                    className="font-mono h-full overflow-y-auto resize-none border-0 focus-visible:ring-0 focus-visible:ring-offset-0 leading-[1.4] text-sm py-3"
                    style={{ lineHeight: '20px' }}
                  />
                </div>
              </div>
              <div className="text-xs text-muted-foreground">
                {t("targetCount", { count: targetCount })}
              </div>
              {invalidTargets.length > 0 && (
                <div className="text-xs text-destructive">
                  {t("invalidCount", {
                    count: invalidTargets.length,
                    line: invalidTargets[0].index + 1,
                    target: invalidTargets[0].originalTarget,
                    error: invalidTargets[0].error,
                  })}
                </div>
              )}
            </div>

            {/* Organization (optional, searchable, paginated) */}
            <div className="grid gap-2">
              <Label htmlFor="organization">
                {t("linkOrganization")}
              </Label>
              <Button
                type="button"
                variant="outline"
                role="combobox"
                className="w-full justify-between"
                onClick={() => setOrgPickerOpen(true)}
                disabled={batchCreateTargets.isPending || isLoadingOrganizations}
              >
                {isLoadingOrganizations ? (
                  <span className="flex items-center gap-2">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    {t("loading")}
                  </span>
                ) : formData.organizationId ? (
                  <span className="flex items-center gap-2">
                    <Building2 className="h-4 w-4" />
                    <span className="truncate">{selectedOrgName}</span>
                  </span>
                ) : (
                  t("selectOrganization")
                )}
                <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
              </Button>
              <CommandDialog
                open={orgPickerOpen}
                onOpenChange={(o) => {
                  setOrgPickerOpen(o)
                  if (!o) {
                    setOrgSearchQuery("")
                    setOrgPage(1)
                    setOrgPageSize(20)
                  }
                }}
              >
                <CommandInput
                  placeholder={t("searchOrganization")}
                  value={orgSearchQuery}
                  onValueChange={(v) => setOrgSearchQuery(v)}
                />
                <CommandList className="max-h-[300px] overflow-y-auto overscroll-contain">
                  {isLoadingOrganizations ? (
                    <div className="py-6 text-center text-sm">
                      <Loader2 className="mx-auto h-4 w-4 animate-spin" />
                    </div>
                  ) : filteredOrganizations.length === 0 ? (
                    <CommandEmpty>{t("noOrganization")}</CommandEmpty>
                  ) : (
                    <CommandGroup>
                      <div className="grid grid-cols-2 gap-1 p-1">
                        {filteredOrganizations.map((org) => (
                          <CommandItem
                            key={org.id}
                            value={org.id.toString()}
                            onSelect={() => handleSelectOrganization(org.id.toString(), org.name)}
                            className="cursor-pointer"
                          >
                            <Check
                              className={cn(
                                "mr-1 h-3.5 w-3.5 flex-shrink-0",
                                formData.organizationId === org.id.toString()
                                  ? "opacity-100"
                                  : "opacity-0"
                              )}
                            />
                            <Building2 className="mr-1 h-3.5 w-3.5 flex-shrink-0" />
                            <span className="font-medium text-sm truncate">{org.name}</span>
                          </CommandItem>
                        ))}
                      </div>
                    </CommandGroup>
                  )}
                </CommandList>
                {organizationsData && (
                  <div className="flex items-center justify-between border-t p-2 bg-muted/50">
                    <div className="text-xs text-muted-foreground">
                      {t("orgPagination", {
                        total: organizationsData.pagination.total,
                        page: organizationsData.pagination.page,
                        totalPages: organizationsData.pagination.totalPages,
                      })}
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="flex items-center gap-1">
                        <span className="text-xs text-muted-foreground">{t("perPage")}</span>
                        <Select value={orgPageSize.toString()} onValueChange={(value) => {
                          setOrgPageSize(Number(value))
                          setOrgPage(1)
                        }}>
                          <SelectTrigger className="h-7 w-16 text-xs">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            {pageSizeOptions.map((size) => (
                              <SelectItem key={size} value={size.toString()}>
                                {size}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Button
                          type="button"
                          variant="outline"
                          className="hidden h-8 w-8 p-0 lg:flex"
                          onClick={() => setOrgPage(1)}
                          disabled={orgPage === 1 || isLoadingOrganizations}
                        >
                          <span className="sr-only">{tPagination("first")}</span>
                          <IconChevronsLeft />
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          className="h-8 w-8 p-0"
                          onClick={() => setOrgPage(prev => Math.max(1, prev - 1))}
                          disabled={orgPage === 1 || isLoadingOrganizations}
                        >
                          <span className="sr-only">{tPagination("previous")}</span>
                          <IconChevronLeft />
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          className="h-8 w-8 p-0"
                          onClick={() => setOrgPage(prev => Math.min(organizationsData.pagination.totalPages, prev + 1))}
                          disabled={orgPage === organizationsData.pagination.totalPages || isLoadingOrganizations}
                        >
                          <span className="sr-only">{tPagination("next")}</span>
                          <IconChevronRight />
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          className="hidden h-8 w-8 p-0 lg:flex"
                          onClick={() => setOrgPage(organizationsData.pagination.totalPages)}
                          disabled={orgPage === organizationsData.pagination.totalPages || isLoadingOrganizations}
                        >
                          <span className="sr-only">{tPagination("last")}</span>
                          <IconChevronsRight />
                        </Button>
                      </div>
                    </div>
                  </div>
                )}
              </CommandDialog>
            </div>
          </div>
          
          {/* Dialog footer buttons */}
          <DialogFooter>
            <Button 
              type="button" 
              variant="outline" 
              onClick={() => handleOpenChange(false)}
              disabled={batchCreateTargets.isPending}
            >
              {tCommon("cancel")}
            </Button>
            <Button 
              type="submit" 
              disabled={batchCreateTargets.isPending || !isFormValid}
            >
              {batchCreateTargets.isPending ? (
                <>
                  <LoadingSpinner/>
                  {t("creating")}
                </>
              ) : (
                <>
                  <Plus />
                  {t("addTarget")}
                </>
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
