"use client"

import React from "react"
import dynamic from "next/dynamic"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { Link } from "@/i18n/navigation"

type DynamicModuleDemoProps = {
  loader: () => Promise<Record<string, any>>
  props?: Record<string, any>
  title?: string
  description?: string
  fallbackRoute?: string
  className?: string
}

type ErrorBoundaryState = {
  hasError: boolean
  message?: string
}

class DemoErrorBoundary extends React.Component<React.PropsWithChildren, ErrorBoundaryState> {
  state: ErrorBoundaryState = { hasError: false }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, message: error.message }
  }

  render() {
    if (!this.state.hasError) return this.props.children

    return (
      <Alert className="border-destructive/30 text-destructive">
        <AlertTitle>组件渲染失败</AlertTitle>
        <AlertDescription className="text-xs text-destructive/80">
          {this.state.message || "该组件需要额外上下文或参数。请参考对应业务页面。"}
        </AlertDescription>
      </Alert>
    )
  }
}

const pickComponent = (mod: Record<string, any>) => {
  const isReactComponent = (value: any) =>
    typeof value === "function" ||
    (typeof value === "object" && value && value.$$typeof)

  if (isReactComponent(mod.default)) return mod.default

  const named = Object.values(mod).find(isReactComponent)
  return named
}

export function DynamicModuleDemo({
  loader,
  props,
  title,
  description,
  fallbackRoute,
  className,
}: DynamicModuleDemoProps) {
  const DemoComponent = React.useMemo(
    () =>
      dynamic(async () => {
        const mod = await loader()
        const Component = pickComponent(mod)

        if (!Component) {
          return () => (
            <Alert>
              <AlertTitle>无可渲染导出</AlertTitle>
              <AlertDescription className="text-xs text-muted-foreground">
                该模块未导出可直接渲染的组件。
              </AlertDescription>
            </Alert>
          )
        }

        return (componentProps: Record<string, any>) => (
          <Component {...componentProps} />
        )
      }, {
        ssr: false,
        loading: () => (
          <Alert>
            <AlertTitle>组件加载中</AlertTitle>
            <AlertDescription className="text-xs text-muted-foreground">
              正在动态加载组件…
            </AlertDescription>
          </Alert>
        ),
      }),
    [loader]
  )

  return (
    <div className={cn("space-y-4", className)}>
      {title ? (
        <div>
          <h2 className="text-lg font-semibold">{title}</h2>
          {description ? (
            <p className="text-sm text-muted-foreground">{description}</p>
          ) : null}
        </div>
      ) : null}
      <DemoErrorBoundary>
        <DemoComponent {...(props || {})} />
      </DemoErrorBoundary>
      {fallbackRoute ? (
        <Button asChild variant="outline" size="sm">
          <Link href={fallbackRoute}>进入模块页面</Link>
        </Button>
      ) : null}
    </div>
  )
}
