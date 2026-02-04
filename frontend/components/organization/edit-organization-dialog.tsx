"use client"

import React, { useEffect } from "react"
import { Edit, Building2 } from "@/components/icons"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useTranslations } from "next-intl"

import { Button } from "@/components/ui/button"
import { LoadingSpinner } from "@/components/loading-spinner"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"

import { useUpdateOrganization } from "@/hooks/use-organizations"

import type { Organization } from "@/types/organization.types"

interface EditOrganizationDialogProps {
  organization: Organization
  open: boolean
  onOpenChange: (open: boolean) => void
  onEdit: (organization: Organization) => void
}

export function EditOrganizationDialog({ 
  organization, 
  open, 
  onOpenChange, 
  onEdit 
}: EditOrganizationDialogProps) {
  const t = useTranslations("organization.dialog")
  const tValidation = useTranslations("organization.validation")

  const formSchema = z.object({
    name: z.string()
      .min(2, { message: tValidation("nameMin", { min: 2 }) })
      .max(50, { message: tValidation("nameMax", { max: 50 }) }),
    description: z.string().max(200, { message: tValidation("descMax", { max: 200 }) }).optional(),
  })

  type FormValues = z.infer<typeof formSchema>

  const updateOrganization = useUpdateOrganization()

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: organization?.name || "",
      description: organization?.description || "",
    },
  })

  useEffect(() => {
    if (organization) {
      form.reset({
        name: organization.name || "",
        description: organization.description || "",
      })
    }
  }, [organization, form])

  const hasChanges = form.formState.isDirty

  const onSubmit = (values: FormValues) => {
    updateOrganization.mutate(
      {
        id: Number(organization.id),
        data: {
          name: values.name.trim(),
          description: values.description?.trim() || "",
        }
      },
      {
        onSuccess: (updatedOrganization) => {
          onEdit(updatedOrganization)
          onOpenChange(false)
        }
      }
    )
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!updateOrganization.isPending) {
      onOpenChange(newOpen)
    }
  }

  const handleReset = () => {
    form.reset({
      name: organization.name || "",
      description: organization.description || "",
    })
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Building2 />
            <span>{t("editTitle")}</span>
          </DialogTitle>
          <DialogDescription>{t("editDesc")}</DialogDescription>
        </DialogHeader>
        
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="grid gap-4 py-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {t("orgName")} <span className="text-destructive">*</span>
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("orgNamePlaceholder")}
                        disabled={updateOrganization.isPending}
                        maxLength={50}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      {t("characters", { count: field.value.length, max: 50 })}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              
              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("orgDesc")}</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder={t("orgDescPlaceholder")}
                        disabled={updateOrganization.isPending}
                        rows={3}
                        maxLength={200}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      {t("characters", { count: (field.value || "").length, max: 200 })}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {hasChanges && (
                <div className="text-xs text-amber-600 bg-amber-50 dark:bg-amber-950/20 p-2 rounded">
                  {t("changesDetected")}
                </div>
              )}
            </div>
            
            <DialogFooter className="gap-2">
              <Button 
                type="button" 
                variant="outline" 
                onClick={() => handleOpenChange(false)}
                disabled={updateOrganization.isPending}
              >
                {t("cancel")}
              </Button>
              
              {hasChanges && (
                <Button 
                  type="button" 
                  variant="ghost" 
                  onClick={handleReset}
                  disabled={updateOrganization.isPending}
                >
                  {t("reset")}
                </Button>
              )}
              
              <Button 
                type="submit" 
                disabled={updateOrganization.isPending || !form.formState.isValid || !hasChanges}
              >
                {updateOrganization.isPending ? (
                  <>
                    <LoadingSpinner/>
                    {t("updating")}
                  </>
                ) : (
                  <>
                    <Edit/>
                    {t("update")}
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
