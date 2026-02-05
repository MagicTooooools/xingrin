import { useEffect, useCallback } from 'react'
import { useLocale } from 'next-intl'
import { useRouter } from '@/i18n/navigation'

/**
 * 路由预加载 Hook
 * 在页面加载完成后，后台预加载其他页面的 JS/CSS 资源
 * 不会发送 API 请求，只加载页面组件
 * @param currentPath 当前页面路径（可选），如果提供则会智能预加载相关动态路由
 */
export function useRoutePrefetch(currentPath?: string) {
  const router = useRouter()
  const locale = useLocale()
  const withLocale = useCallback((path: string) => {
    if (path.startsWith(`/${locale}/`)) return path
    const normalized = path.startsWith("/") ? path : `/${path}`
    return `/${locale}${normalized}`
  }, [locale])

  useEffect(() => {
    const w = typeof window !== 'undefined'
      ? (window as Window & { __lunafoxRoutePrefetchDone?: boolean })
      : null
    const hasPrefetched = !!w?.__lunafoxRoutePrefetchDone

    // 使用 requestIdleCallback 在浏览器空闲时预加载，不影响当前页面渲染
    const prefetchBaseRoutes = () => {
      const criticalRoutes = ['/dashboard/', '/organization/', '/target/']
      const secondaryRoutes = ['/scan/history/', '/vulnerabilities/']
      const lowPriorityRoutes = [
        '/scan/scheduled/',
        '/scan/engine/',
        '/tools/',
        '/settings/workers/',
      ]

      criticalRoutes.forEach((route) => {
        router.prefetch(withLocale(route))
      })

      const scheduleSecondary = () => {
        secondaryRoutes.forEach((route) => {
          router.prefetch(withLocale(route))
        })
      }

      const scheduleLowPriority = () => {
        lowPriorityRoutes.forEach((route) => {
          router.prefetch(withLocale(route))
        })
      }

      if (typeof window !== 'undefined') {
        if ('requestIdleCallback' in window) {
          window.requestIdleCallback(scheduleSecondary, { timeout: 2000 })
          window.requestIdleCallback(scheduleLowPriority, { timeout: 4000 })
          return
        }
      }

      scheduleSecondary()
      setTimeout(scheduleLowPriority, 1200)
    }

    const prefetchDynamicRoutes = () => {
      if (!currentPath) return
      // 如果是目标详情页（如 /target/146），预加载子路由
      const targetIdMatch = currentPath.match(/\/target\/(\d+)$/)
      if (targetIdMatch) {
        const targetId = targetIdMatch[1]
        const subRoutes = ['subdomain', 'endpoints', 'websites', 'vulnerabilities', 'directories', 'ip-addresses']
        subRoutes.forEach(sub => {
          router.prefetch(withLocale(`/target/${targetId}/${sub}`))
        })
      }
      
      // 如果是扫描历史详情页（如 /scan/history/146），预加载子路由
      const scanIdMatch = currentPath.match(/\/scan\/history\/(\d+)$/)
      if (scanIdMatch) {
        const scanId = scanIdMatch[1]
        const subRoutes = ['subdomain', 'endpoints', 'websites', 'vulnerabilities', 'directories', 'ip-addresses']
        subRoutes.forEach(sub => {
          router.prefetch(withLocale(`/scan/history/${scanId}/${sub}`))
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
      return
    }

    // 使用 requestIdleCallback 在浏览器空闲时执行，如果不支持则立即执行
    if (typeof window !== 'undefined' && 'requestIdleCallback' in window) {
      const idleId = window.requestIdleCallback(runPrefetch)
      return () => window.cancelIdleCallback(idleId)
    }

    runPrefetch()
    return
  }, [router, currentPath, withLocale])
}

/**
 * 智能路由预加载 Hook
 * 根据当前路径，预加载用户可能访问的下一个页面
 * @param currentPath 当前页面路径
 */
export function useSmartRoutePrefetch(currentPath: string) {
  const router = useRouter()
  const locale = useLocale()
  const withLocale = useCallback((path: string) => {
    if (path.startsWith(`/${locale}/`)) return path
    const normalized = path.startsWith("/") ? path : `/${path}`
    return `/${locale}${normalized}`
  }, [locale])

  useEffect(() => {
    const timer = setTimeout(() => {
      if (currentPath.includes('/organization')) {
        // 在组织页面，预加载目标页面
        router.prefetch(withLocale('/target/'))
      } else if (currentPath.includes('/target')) {
        // 在目标页面，预加载扫描和漏洞页面
        router.prefetch(withLocale('/scan/history/'))
        router.prefetch(withLocale('/vulnerabilities/'))

        // 如果是目标详情页（如 /target/146），预加载子路由
        const targetIdMatch = currentPath.match(/\/target\/(\d+)$/)
        if (targetIdMatch) {
          const targetId = targetIdMatch[1]
          const subRoutes = ['subdomain', 'endpoints', 'websites', 'vulnerabilities']
          subRoutes.forEach(sub => {
            router.prefetch(withLocale(`/target/${targetId}/${sub}`))
          })
        }
      } else if (currentPath.includes('/scan/history')) {
        // 在扫描历史页面，预加载目标页面
        router.prefetch(withLocale('/target/'))
        router.prefetch(withLocale('/vulnerabilities/'))
      } else if (currentPath === '/') {
        // 在首页，预加载主要页面
        router.prefetch(withLocale('/dashboard/'))
        router.prefetch(withLocale('/organization/'))
      }
    }, 1500) // 1.5 秒后预加载

    return () => clearTimeout(timer)
  }, [currentPath, router, withLocale])
}
