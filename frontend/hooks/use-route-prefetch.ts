import { useEffect, useCallback, useRef } from 'react'
import { useRouter } from '@/i18n/navigation'

const BASE_CRITICAL_ROUTES = ['/dashboard/'] as const
const BASE_SECONDARY_ROUTES = ['/organization/', '/target/'] as const
const BASE_LOW_PRIORITY_ROUTES = ['/scan/history/', '/vulnerabilities/'] as const

const DETAIL_SUB_ROUTES = [
  'subdomain',
  'endpoints',
  'websites',
  'vulnerabilities',
  'directories',
  'ip-addresses',
] as const

type NetworkInformation = {
  saveData?: boolean
  effectiveType?: string
}

type NavigatorWithConnection = Navigator & {
  connection?: NetworkInformation
}

const getNetworkConnection = (): NetworkInformation | undefined => {
  if (typeof navigator === 'undefined') return undefined
  return (navigator as NavigatorWithConnection).connection
}

const canPrefetchSecondaryRoutes = (): boolean => {
  const connection = getNetworkConnection()
  if (!connection) return true
  if (connection.saveData) return false
  const effectiveType = connection.effectiveType
  return effectiveType !== 'slow-2g' && effectiveType !== '2g'
}

const canPrefetchLowPriorityRoutes = (): boolean => {
  const connection = getNetworkConnection()
  if (!connection) return true
  if (connection.saveData) return false
  return connection.effectiveType === '4g'
}

/**
 * 路由预加载 Hook
 * 在页面加载完成后，后台预加载其他页面的 JS/CSS 资源
 * 不会发送 API 请求，只加载页面组件
 * @param currentPath 当前页面路径（可选），如果提供则会智能预加载相关动态路由
 */
export function useRoutePrefetch(currentPath?: string) {
  const router = useRouter()
  const prefetchedRoutesRef = useRef<Set<string>>(new Set())

  const prefetchOnce = useCallback((path: string) => {
    const normalizedPath = path.startsWith("/") ? path : `/${path}`
    if (prefetchedRoutesRef.current.has(normalizedPath)) return
    prefetchedRoutesRef.current.add(normalizedPath)
    void router.prefetch(normalizedPath)
  }, [router])

  useEffect(() => {
    const w = typeof window !== 'undefined'
      ? (window as Window & { __lunafoxRoutePrefetchDone?: boolean })
      : null
    const hasPrefetched = !!w?.__lunafoxRoutePrefetchDone
    const allowSecondaryPrefetch = canPrefetchSecondaryRoutes()
    const allowLowPriorityPrefetch = canPrefetchLowPriorityRoutes()
    const idleTaskIds: number[] = []
    const timeoutIds: Array<ReturnType<typeof setTimeout>> = []

    const prefetchBatch = (routes: readonly string[]) => {
      routes.forEach((route) => {
        prefetchOnce(route)
      })
    }

    // 使用 requestIdleCallback 在浏览器空闲时预加载，不影响当前页面渲染
    const prefetchBaseRoutes = () => {
      prefetchBatch(BASE_CRITICAL_ROUTES)
      if (!allowSecondaryPrefetch) return

      const scheduleSecondary = () => {
        prefetchBatch(BASE_SECONDARY_ROUTES)
      }

      const scheduleLowPriority = () => {
        if (!allowLowPriorityPrefetch) return
        prefetchBatch(BASE_LOW_PRIORITY_ROUTES)
      }

      if (typeof window !== 'undefined') {
        if ('requestIdleCallback' in window) {
          idleTaskIds.push(window.requestIdleCallback(scheduleSecondary, { timeout: 2000 }))
          idleTaskIds.push(window.requestIdleCallback(scheduleLowPriority, { timeout: 4000 }))
          return
        }
      }

      scheduleSecondary()
      timeoutIds.push(setTimeout(scheduleLowPriority, 2000))
    }

    const prefetchDynamicRoutes = () => {
      if (!currentPath || !allowSecondaryPrefetch) return
      // 如果是目标详情页（如 /target/146），预加载子路由
      const targetIdMatch = currentPath.match(/\/target\/(\d+)$/)
      if (targetIdMatch) {
        const targetId = targetIdMatch[1]
        DETAIL_SUB_ROUTES.forEach((subRoute) => {
          prefetchOnce(`/target/${targetId}/${subRoute}`)
        })
      }
      
      // 如果是扫描历史详情页（如 /scan/history/146），预加载子路由
      const scanIdMatch = currentPath.match(/\/scan\/history\/(\d+)$/)
      if (scanIdMatch) {
        const scanId = scanIdMatch[1]
        DETAIL_SUB_ROUTES.forEach((subRoute) => {
          prefetchOnce(`/scan/history/${scanId}/${subRoute}`)
        })
      }
    }

    const runPrefetch = () => {
      if (!hasPrefetched) {
        prefetchBaseRoutes()
        if (w) {
          w.__lunafoxRoutePrefetchDone = true
          w.dispatchEvent(new Event('lunafox:route-prefetch-done'))
        }
      }
      prefetchDynamicRoutes()
    }

    if (hasPrefetched) {
      runPrefetch()
      return () => {
        idleTaskIds.forEach((id) => {
          if (typeof window !== 'undefined' && 'cancelIdleCallback' in window) {
            window.cancelIdleCallback(id)
          }
        })
        timeoutIds.forEach((id) => clearTimeout(id))
      }
    }

    // 使用 requestIdleCallback 在浏览器空闲时执行，如果不支持则立即执行
    if (typeof window !== 'undefined' && 'requestIdleCallback' in window) {
      const idleId = window.requestIdleCallback(runPrefetch)
      return () => {
        window.cancelIdleCallback(idleId)
        idleTaskIds.forEach((id) => window.cancelIdleCallback(id))
        timeoutIds.forEach((id) => clearTimeout(id))
      }
    }

    runPrefetch()
    return () => {
      timeoutIds.forEach((id) => clearTimeout(id))
    }
  }, [currentPath, prefetchOnce])
}

/**
 * 智能路由预加载 Hook
 * 根据当前路径，预加载用户可能访问的下一个页面
 * @param currentPath 当前页面路径
 */
export function useSmartRoutePrefetch(currentPath: string) {
  const router = useRouter()
  const prefetchedRoutesRef = useRef<Set<string>>(new Set())

  const prefetchOnce = useCallback((path: string) => {
    const normalizedPath = path.startsWith("/") ? path : `/${path}`
    if (prefetchedRoutesRef.current.has(normalizedPath)) return
    prefetchedRoutesRef.current.add(normalizedPath)
    void router.prefetch(normalizedPath)
  }, [router])

  useEffect(() => {
    const timer = setTimeout(() => {
      if (currentPath.includes('/organization')) {
        // 在组织页面，预加载目标页面
        prefetchOnce('/target/')
      } else if (currentPath.includes('/target')) {
        // 在目标页面，预加载扫描和漏洞页面
        prefetchOnce('/scan/history/')
        prefetchOnce('/vulnerabilities/')

        // 如果是目标详情页（如 /target/146），预加载子路由
        const targetIdMatch = currentPath.match(/\/target\/(\d+)$/)
        if (targetIdMatch) {
          const targetId = targetIdMatch[1]
          const subRoutes = ['subdomain', 'endpoints', 'websites', 'vulnerabilities']
          subRoutes.forEach((sub) => {
            prefetchOnce(`/target/${targetId}/${sub}`)
          })
        }
      } else if (currentPath.includes('/scan/history')) {
        // 在扫描历史页面，预加载目标页面
        prefetchOnce('/target/')
        prefetchOnce('/vulnerabilities/')
      } else if (currentPath === '/') {
        // 在首页，预加载主要页面
        prefetchOnce('/dashboard/')
        prefetchOnce('/organization/')
      }
    }, 1500) // 1.5 秒后预加载

    return () => clearTimeout(timer)
  }, [currentPath, prefetchOnce])
}
