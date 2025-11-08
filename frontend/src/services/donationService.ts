import api from './api';
import type { Donation } from '../types';

export interface CreateDonationRequest {
  pantry_id: string;
  donor_name: string;
  donor_email?: string;
  donor_phone?: string;
  amount?: number;
  description: string;
  donation_date?: string;
}

export interface UpdateDonationRequest {
  donor_name?: string;
  donor_email?: string;
  donor_phone?: string;
  amount?: number;
  description?: string;
  donation_date?: string;
  receipt_sent?: boolean;
}

export interface GetDonationsParams {
  pantry_id?: string;
  receipt_sent?: boolean;
  start_date?: string;
  end_date?: string;
  page?: number;
  page_size?: number;
}

export interface GetDonationsResponse {
  donations: Donation[];
  total: number;
  page: number;
  pages: number;
}

export interface DonationStatsResponse {
  total_donations: number;
  total_amount: number;
  donor_count: number;
  receipts_pending: number;
  monetary_count: number;
  in_kind_count: number;
}

export const donationService = {
  // Public: Create donation
  async createDonation(data: CreateDonationRequest): Promise<Donation> {
    const response = await api.post<Donation>('/donations', data);
    return response.data;
  },

  // Admin: Get list of donations
  async getDonations(params?: GetDonationsParams): Promise<GetDonationsResponse> {
    const response = await api.get<GetDonationsResponse>('/admin/donations', { params });
    return response.data;
  },

  // Admin: Get a specific donation by ID
  async getDonation(donationId: string): Promise<Donation> {
    const response = await api.get<Donation>(`/admin/donations/${donationId}`);
    return response.data;
  },

  // Admin: Update donation
  async updateDonation(donationId: string, data: UpdateDonationRequest): Promise<Donation> {
    const response = await api.put<Donation>(`/admin/donations/${donationId}`, data);
    return response.data;
  },

  // Admin: Delete donation
  async deleteDonation(donationId: string): Promise<void> {
    await api.delete(`/admin/donations/${donationId}`);
  },

  // Admin: Mark receipt as sent
  async markReceiptSent(donationId: string): Promise<Donation> {
    const response = await api.patch<Donation>(`/admin/donations/${donationId}/receipt`);
    return response.data;
  },

  // Admin: Search donations
  async searchDonations(query: string, page?: number, pageSize?: number): Promise<GetDonationsResponse> {
    const response = await api.get<GetDonationsResponse>('/admin/donations/search', {
      params: { q: query, page, page_size: pageSize },
    });
    return response.data;
  },

  // Admin: Get donation statistics
  async getDonationStats(pantryId?: string, startDate?: string, endDate?: string): Promise<DonationStatsResponse> {
    const response = await api.get<DonationStatsResponse>('/admin/donations/stats', {
      params: { pantry_id: pantryId, start_date: startDate, end_date: endDate },
    });
    return response.data;
  },

  // Admin: Get donations by donor email
  async getDonationsByDonor(email: string): Promise<Donation[]> {
    const response = await api.get<Donation[]>('/admin/donations/by-donor', {
      params: { email },
    });
    return response.data;
  },
};
