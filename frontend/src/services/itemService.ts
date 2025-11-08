import api from './api';
import type { Item } from '../types';

export interface CreateItemRequest {
  name: string;
  description?: string;
  category_id: string;
  pantry_id: string;
  quantity: number;
  low_stock_threshold: number;
  unit: string;
  image_url?: string;
  is_available: boolean;
}

export interface UpdateItemRequest {
  name?: string;
  description?: string;
  category_id?: string;
  quantity?: number;
  low_stock_threshold?: number;
  unit?: string;
  image_url?: string;
  is_available?: boolean;
}

export interface ItemListResponse {
  data: Item[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ItemListParams {
  pantry_id?: string;
  category_id?: string;
  search?: string;
  available?: boolean;
  low_stock?: boolean;
  page?: number;
  page_size?: number;
}

export const itemService = {
  async list(params: ItemListParams = {}): Promise<ItemListResponse> {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        queryParams.append(key, value.toString());
      }
    });

    const response = await api.get<ItemListResponse>(
      `/admin/items?${queryParams.toString()}`
    );
    return response.data;
  },

  async listPublic(params: ItemListParams = {}): Promise<ItemListResponse> {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        queryParams.append(key, value.toString());
      }
    });

    const response = await api.get<ItemListResponse>(
      `/items?${queryParams.toString()}`
    );
    return response.data;
  },

  async get(id: string): Promise<Item> {
    const response = await api.get<Item>(`/admin/items/${id}`);
    return response.data;
  },

  async create(data: CreateItemRequest): Promise<Item> {
    const response = await api.post<Item>('/admin/items', data);
    return response.data;
  },

  async update(id: string, data: UpdateItemRequest): Promise<Item> {
    const response = await api.put<Item>(`/admin/items/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/admin/items/${id}`);
  },

  async updateQuantity(id: string, quantity: number): Promise<void> {
    await api.patch(`/admin/items/${id}/quantity`, { quantity });
  },

  async getLowStock(pantryId?: string): Promise<{ data: Item[]; count: number }> {
    const params = pantryId ? `?pantry_id=${pantryId}` : '';
    const response = await api.get<{ data: Item[]; count: number }>(
      `/admin/items/low-stock${params}`
    );
    return response.data;
  },
};
