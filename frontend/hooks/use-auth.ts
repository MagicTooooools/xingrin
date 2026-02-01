/**
 * Authentication-related hooks
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useLocale } from 'next-intl'
import { useRouter } from 'next/navigation'
import { login, logout, getMe, changePassword } from '@/services/auth.service'
import { useToastMessages } from '@/lib/toast-helpers'
import { getErrorCode, getErrorResponseData } from '@/lib/response-parser'
import type { LoginRequest, ChangePasswordRequest } from '@/types/auth.types'

/**
 * Get current user information
 */
export function useAuth() {
  const skipAuth = process.env.NEXT_PUBLIC_SKIP_AUTH === 'true'
  
  return useQuery({
    queryKey: ['auth', 'me'],
    queryFn: skipAuth 
      ? () => Promise.resolve({ authenticated: true } as Awaited<ReturnType<typeof getMe>>)
      : getMe,
    staleTime: 1000 * 60 * 5, // Don't re-request within 5 minutes
    retry: false,
  })
}

/**
 * User login
 */
export function useLogin() {
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: (data: LoginRequest) => login(data),
    onSuccess: () => {
      // Navigation and data prefetch are handled by the login page.
      toastMessages.success('toast.auth.login.success')
    },
    onError: (error: unknown) => {
      const errorCode = getErrorCode(getErrorResponseData(error))
      if (errorCode) {
        toastMessages.errorFromCode(errorCode)
      } else {
        toastMessages.error('auth.loginFailed')
      }
    },
  })
}

/**
 * User logout
 */
export function useLogout() {
  const queryClient = useQueryClient()
  const router = useRouter()
  const toastMessages = useToastMessages()
  const locale = useLocale()

  return useMutation({
    mutationFn: logout,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth', 'me'] })
      toastMessages.success('toast.auth.logout.success')
      router.push(`/${locale}/login/`)
    },
    onError: (error: unknown) => {
      const errorCode = getErrorCode(getErrorResponseData(error))
      if (errorCode) {
        toastMessages.errorFromCode(errorCode)
      } else {
        toastMessages.error('errors.unknown')
      }
    },
  })
}

/**
 * Change password
 */
export function useChangePassword() {
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: (data: ChangePasswordRequest) => changePassword(data),
    onSuccess: () => {
      toastMessages.success('toast.auth.changePassword.success')
    },
    onError: (error: unknown) => {
      const errorCode = getErrorCode(getErrorResponseData(error))
      if (errorCode) {
        toastMessages.errorFromCode(errorCode)
      } else {
        toastMessages.error('toast.auth.changePassword.error')
      }
    },
  })
}
