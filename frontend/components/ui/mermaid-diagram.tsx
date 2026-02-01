"use client"

import { useEffect, useRef, useState } from "react"
import mermaid from "mermaid"
import { useTheme } from "next-themes"

interface MermaidDiagramProps {
  chart: string
  className?: string
}

export function MermaidDiagram({ chart, className = "" }: MermaidDiagramProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const [svg, setSvg] = useState<string>("")
  const { theme } = useTheme()

  useEffect(() => {
    const renderDiagram = async () => {
      if (!containerRef.current) return

      try {
        const rootStyle = getComputedStyle(document.documentElement)
        const resolveColor = (value: string, fallback: string) => {
          if (!value || !document.body) return fallback
          const probe = document.createElement("span")
          probe.style.color = value
          probe.style.position = "absolute"
          probe.style.opacity = "0"
          probe.style.pointerEvents = "none"
          document.body.appendChild(probe)
          const resolved = getComputedStyle(probe).color || fallback
          probe.remove()
          return resolved
        }
        const colorFromVar = (name: string, fallback: string) => {
          const raw = rootStyle.getPropertyValue(name).trim()
          if (!raw) return fallback
          const value = raw.includes("(") ? raw : `hsl(${raw})`
          return resolveColor(value, fallback)
        }

        const background = colorFromVar("--background", "#ffffff")
        const card = colorFromVar("--card", "#ffffff")
        const foreground = colorFromVar("--foreground", "#111827")
        const mutedForeground = colorFromVar("--muted-foreground", "#6b7280")
        const border = colorFromVar("--border", "#e5e7eb")

        // 配置 Mermaid
        mermaid.initialize({
          startOnLoad: false,
          theme: theme === "dark" ? "dark" : "base",
          themeVariables: {
            fontFamily:
              "var(--font-sans, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, \"Segoe UI\", sans-serif)",
            primaryColor: card,
            primaryTextColor: foreground,
            primaryBorderColor: border,
            lineColor: border,
            secondaryColor: card,
            tertiaryColor: card,
            background: "transparent",
            mainBkg: card,
            secondBkg: card,
            tertiaryBkg: card,
            nodeBorder: border,
            nodeTextColor: foreground,
            clusterBkg: background,
            clusterBorder: border,
            defaultLinkColor: border,
            titleColor: foreground,
            edgeLabelBackground: background,
            textColor: mutedForeground,
          },
          flowchart: {
            useMaxWidth: true,
            htmlLabels: false,
            curve: "linear",
          },
        })

        // 生成唯一 ID
        const id = `mermaid-${Math.random().toString(36).substr(2, 9)}`

        // 渲染图表
        const { svg: renderedSvg } = await mermaid.render(id, chart)
        setSvg(renderedSvg)
      } catch (error) {
        console.error("Mermaid rendering error:", error)
      }
    }

    renderDiagram()
  }, [chart, theme])

  return (
    <div
      ref={containerRef}
      className={`mermaid-container ${className}`}
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  )
}
