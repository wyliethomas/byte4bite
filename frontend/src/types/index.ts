// User types
export type UserRole = 'admin' | 'user';

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  role: UserRole;
  pantry_id?: string;
  created_at: string;
  updated_at: string;
}

// Pantry types
export interface Pantry {
  id: string;
  name: string;
  address: string;
  city: string;
  state: string;
  zip_code: string;
  contact_email: string;
  contact_phone?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// Category types
export interface Category {
  id: string;
  name: string;
  description?: string;
  pantry_id: string;
  created_at: string;
  updated_at: string;
}

// Item types
export interface Item {
  id: string;
  name: string;
  description?: string;
  category_id: string;
  category?: Category;
  pantry_id: string;
  quantity: number;
  low_stock_threshold: number;
  unit: string;
  image_url?: string;
  is_available: boolean;
  created_at: string;
  updated_at: string;
}

// Cart types
export type CartStatus = 'active' | 'submitted' | 'cancelled';

export interface CartItem {
  id: string;
  cart_id: string;
  item_id: string;
  item?: Item;
  quantity: number;
  created_at: string;
  updated_at: string;
}

export interface Cart {
  id: string;
  user_id: string;
  pantry_id: string;
  status: CartStatus;
  items?: CartItem[];
  created_at: string;
  updated_at: string;
}

// Order types
export type OrderStatus = 'pending' | 'preparing' | 'ready' | 'picked_up' | 'cancelled';

export interface Order {
  id: string;
  cart_id: string;
  cart?: Cart;
  user_id: string;
  user?: User;
  pantry_id: string;
  status: OrderStatus;
  notes?: string;
  assigned_to_id?: string;
  assigned_to?: User;
  submitted_at: string;
  ready_at?: string;
  picked_up_at?: string;
  created_at: string;
  updated_at: string;
}

// Donation types
export interface Donation {
  id: string;
  pantry_id: string;
  pantry?: Pantry;
  donor_name: string;
  donor_email?: string;
  donor_phone?: string;
  amount?: number;
  description: string;
  donation_date: string;
  receipt_sent: boolean;
  created_at: string;
  updated_at: string;
}

// Auth types
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
  pantry_id?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}
