import api from './api';
import type { Order, OrderStatus } from '../types';

export interface GetOrdersParams {
  status?: OrderStatus;
  page?: number;
  page_size?: number;
}

export interface GetOrdersResponse {
  orders: Order[];
  total: number;
  page: number;
  pages: number;
}

export interface UpdateStatusRequest {
  status: OrderStatus;
}

export interface AssignStaffRequest {
  staff_id: string;
}

export const orderService = {
  // Get list of orders (users see their own, admins see all)
  async getOrders(params?: GetOrdersParams): Promise<GetOrdersResponse> {
    const response = await api.get<GetOrdersResponse>('/orders', { params });
    return response.data;
  },

  // Get a specific order by ID
  async getOrder(orderId: string): Promise<Order> {
    const response = await api.get<Order>(`/orders/${orderId}`);
    return response.data;
  },

  // Cancel an order
  async cancelOrder(orderId: string): Promise<void> {
    await api.delete(`/orders/${orderId}`);
  },

  // Admin: Update order status
  async updateOrderStatus(orderId: string, status: OrderStatus): Promise<void> {
    await api.put(`/admin/orders/${orderId}/status`, { status });
  },

  // Admin: Assign staff to order
  async assignStaff(orderId: string, staffId: string): Promise<void> {
    await api.put(`/admin/orders/${orderId}/assign`, { staff_id: staffId });
  },

  // Admin: Get all orders with filtering
  async getAdminOrders(params?: GetOrdersParams): Promise<GetOrdersResponse> {
    const response = await api.get<GetOrdersResponse>('/admin/orders', { params });
    return response.data;
  },

  // Admin: Cancel any order
  async adminCancelOrder(orderId: string): Promise<void> {
    await api.delete(`/admin/orders/${orderId}`);
  },
};
