import * as React from "react"
import { act, renderHook } from "@testing-library/react"
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { useNudgeToast, type NudgeToastVariant } from "@/hooks/use-nudge-toast"

const sonnerMocks = vi.hoisted(() => ({
  custom: vi.fn(),
  dismiss: vi.fn(),
}))

vi.mock("sonner", () => ({
  toast: sonnerMocks,
}))

const baseVariant: NudgeToastVariant = {
  title: "Nudge title",
  description: "Nudge description",
  icon: <span>!</span>,
  primaryAction: {
    label: "Open",
    onClick: vi.fn(),
  },
}

describe("useNudgeToast", () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it("triggerWithVariant 会触发 toast.custom 且回调返回可渲染节点", () => {
    const { result } = renderHook(() =>
      useNudgeToast({
        variants: [baseVariant],
        delay: 0,
      })
    )

    act(() => {
      result.current.triggerWithVariant(baseVariant)
    })
    act(() => {
      vi.runAllTimers()
    })

    expect(sonnerMocks.custom).toHaveBeenCalledTimes(1)
    const [renderContent, options] = sonnerMocks.custom.mock.calls[0] as [
      (id: string) => React.ReactNode,
      { duration: number; position: string }
    ]
    const node = renderContent("toast-1")

    expect(React.isValidElement(node)).toBe(true)
    expect(options).toMatchObject({
      duration: Infinity,
      position: "bottom-right",
    })
  })

  it("primaryAction 缺失时仍返回组件并注入兜底动作", () => {
    const unsafeVariant = {
      ...baseVariant,
      primaryAction: undefined,
    } as unknown as NudgeToastVariant

    const { result } = renderHook(() =>
      useNudgeToast({
        variants: [unsafeVariant],
        delay: 0,
      })
    )

    act(() => {
      result.current.triggerWithVariant(unsafeVariant)
    })
    act(() => {
      vi.runAllTimers()
    })

    const [renderContent] = sonnerMocks.custom.mock.calls[0] as [(id: string) => React.ReactNode]
    const node = renderContent("toast-2")

    expect(React.isValidElement(node)).toBe(true)
    const props = (node as React.ReactElement<{ primaryAction: { label: string } }>).props
    expect(props.primaryAction.label).toBe("OK")
  })
})
