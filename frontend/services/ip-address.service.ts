import { api } from "@/lib/api-client"
import type { GetIPAddressesParams, GetIPAddressesResponse } from "@/types/ip-address.types"

// Bulk delete response type
export interface BulkDeleteResponse {
  deletedCount: number
}

export class IPAddressService {
  /**
   * Bulk delete IP addresses
   * POST /api/ip-addresses/bulk-delete/
   * Note: IP addresses are aggregated, so we pass IP strings instead of IDs
   */
  static async bulkDelete(ips: string[]): Promise<BulkDeleteResponse> {
    const response = await api.post<BulkDeleteResponse>(
      `/ip-addresses/bulk-delete/`,
      { ips }
    )
    return response.data
  }

  static async getTargetIPAddresses(
    targetId: number,
    params?: GetIPAddressesParams
  ): Promise<GetIPAddressesResponse> {
    const response = await api.get<GetIPAddressesResponse>(`/targets/${targetId}/ip-addresses/`, {
      params: {
        page: params?.page || 1,
        pageSize: params?.pageSize || 10,
        ...(params?.filter && { filter: params.filter }),
      },
    })
    return response.data
  }

  static async getScanIPAddresses(
    scanId: number,
    params?: GetIPAddressesParams
  ): Promise<GetIPAddressesResponse> {
    const response = await api.get<GetIPAddressesResponse>(`/scans/${scanId}/ip-addresses/`, {
      params: {
        page: params?.page || 1,
        pageSize: params?.pageSize || 10,
        ...(params?.filter && { filter: params.filter }),
      },
    })
    return response.data
  }

  /** Export all IP addresses by target (CSV format) */
  static async exportIPAddressesByTargetId(targetId: number, ips?: string[]): Promise<Blob> {
    const params: Record<string, string> = {}
    if (ips && ips.length > 0) {
      params.ips = ips.join(',')
    }
    const response = await api.get<Blob>(`/targets/${targetId}/ip-addresses/export/`, {
      params,
      responseType: 'blob',
    })
    return response.data
  }

  /** Export all IP addresses by scan task (CSV format) */
  static async exportIPAddressesByScanId(scanId: number): Promise<Blob> {
    const response = await api.get<Blob>(`/scans/${scanId}/ip-addresses/export/`, {
      responseType: 'blob',
    })
    return response.data
  }
}
