import api from './api';
import type { Pantry } from '../types';

export interface GetPantriesParams {
  is_active?: boolean;
  page?: number;
  page_size?: number;
}

export interface GetPantriesResponse {
  pantries: Pantry[];
  total: number;
  page: number;
  pages: number;
}

export interface CreatePantryRequest {
  name: string;
  address: string;
  city: string;
  state: string;
  zip_code: string;
  contact_email: string;
  contact_phone?: string;
  is_active: boolean;
}

export interface UpdatePantryRequest {
  name?: string;
  address?: string;
  city?: string;
  state?: string;
  zip_code?: string;
  contact_email?: string;
  contact_phone?: string;
  is_active?: boolean;
}

export const pantryService = {
  // Get list of pantries
  async getPantries(params?: GetPantriesParams): Promise<GetPantriesResponse> {
    const response = await api.get<GetPantriesResponse>('/pantries', { params });
    return response.data;
  },

  // Get a specific pantry by ID
  async getPantry(pantryId: string): Promise<Pantry> {
    const response = await api.get<Pantry>(`/pantries/${pantryId}`);
    return response.data;
  },

  // Search pantries by name or city
  async searchPantries(query: string): Promise<Pantry[]> {
    const response = await api.get<Pantry[]>('/pantries/search', {
      params: { q: query },
    });
    return response.data;
  },

  // Get pantries by city
  async getPantriesByCity(city: string): Promise<Pantry[]> {
    const response = await api.get<Pantry[]>('/pantries/by-city', {
      params: { city },
    });
    return response.data;
  },

  // Get pantries by zip code
  async getPantriesByZipCode(zipCode: string): Promise<Pantry[]> {
    const response = await api.get<Pantry[]>('/pantries/by-zip', {
      params: { zip: zipCode },
    });
    return response.data;
  },

  // Admin: Create pantry
  async createPantry(data: CreatePantryRequest): Promise<Pantry> {
    const response = await api.post<Pantry>('/admin/pantries', data);
    return response.data;
  },

  // Admin: Update pantry
  async updatePantry(pantryId: string, data: UpdatePantryRequest): Promise<Pantry> {
    const response = await api.put<Pantry>(`/admin/pantries/${pantryId}`, data);
    return response.data;
  },

  // Admin: Delete pantry
  async deletePantry(pantryId: string): Promise<void> {
    await api.delete(`/admin/pantries/${pantryId}`);
  },

  // Admin: Toggle pantry active status
  async togglePantryStatus(pantryId: string): Promise<Pantry> {
    const response = await api.patch<Pantry>(`/admin/pantries/${pantryId}/toggle`);
    return response.data;
  },
};
