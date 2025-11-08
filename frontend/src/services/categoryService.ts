import api from './api';
import type { Category } from '../types';

export interface CreateCategoryRequest {
  name: string;
  description?: string;
  pantry_id: string;
}

export interface UpdateCategoryRequest {
  name?: string;
  description?: string;
}

export interface CategoryListResponse {
  data: Category[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export const categoryService = {
  async list(pantryId?: string, page = 1, pageSize = 20): Promise<CategoryListResponse> {
    const params = new URLSearchParams();
    if (pantryId) params.append('pantry_id', pantryId);
    params.append('page', page.toString());
    params.append('page_size', pageSize.toString());

    const response = await api.get<CategoryListResponse>(
      `/admin/categories?${params.toString()}`
    );
    return response.data;
  },

  async get(id: string): Promise<Category> {
    const response = await api.get<Category>(`/admin/categories/${id}`);
    return response.data;
  },

  async create(data: CreateCategoryRequest): Promise<Category> {
    const response = await api.post<Category>('/admin/categories', data);
    return response.data;
  },

  async update(id: string, data: UpdateCategoryRequest): Promise<Category> {
    const response = await api.put<Category>(`/admin/categories/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/admin/categories/${id}`);
  },
};
