import api from './api';
import type { Cart, CartItem } from '../types';

export interface AddToCartRequest {
  item_id: string;
  quantity: number;
}

export interface UpdateCartItemRequest {
  quantity: number;
}

export interface CheckoutRequest {
  notes?: string;
}

export interface CartResponse {
  cart: Cart | null;
  items: CartItem[];
  count: number;
}

export const cartService = {
  async getCurrentCart(): Promise<CartResponse> {
    const response = await api.get<CartResponse>('/carts/current');
    return response.data;
  },

  async addItem(itemId: string, quantity: number): Promise<Cart> {
    const response = await api.post<Cart>('/carts/items', {
      item_id: itemId,
      quantity,
    });
    return response.data;
  },

  async updateItemQuantity(cartItemId: string, quantity: number): Promise<Cart> {
    const response = await api.put<Cart>(`/carts/items/${cartItemId}`, {
      quantity,
    });
    return response.data;
  },

  async removeItem(cartItemId: string): Promise<Cart> {
    const response = await api.delete<Cart>(`/carts/items/${cartItemId}`);
    return response.data;
  },

  async clearCart(): Promise<void> {
    await api.delete('/carts/current');
  },

  async checkout(notes?: string): Promise<any> {
    const response = await api.post('/carts/checkout', { notes });
    return response.data;
  },
};
