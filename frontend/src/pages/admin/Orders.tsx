import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { orderService } from '../../services/orderService';
import type { Order, OrderStatus } from '../../types';

export const AdminOrders = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [statusFilter, setStatusFilter] = useState<OrderStatus | ''>('');

  useEffect(() => {
    loadOrders();
  }, [statusFilter]);

  const loadOrders = async () => {
    try {
      setIsLoading(true);
      const params = statusFilter ? { status: statusFilter } : {};
      const response = await orderService.getAdminOrders(params);
      setOrders(response.orders);
    } catch (err: any) {
      setError('Failed to load orders');
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateStatus = async (orderId: string, newStatus: OrderStatus) => {
    try {
      setError('');
      setSuccess('');
      await orderService.updateOrderStatus(orderId, newStatus);
      setSuccess('Order status updated successfully');
      loadOrders();
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to update order status');
      setTimeout(() => setError(''), 3000);
    }
  };

  const handleCancelOrder = async (orderId: string) => {
    if (!confirm('Are you sure you want to cancel this order?')) return;

    try {
      setError('');
      setSuccess('');
      await orderService.adminCancelOrder(orderId);
      setSuccess('Order cancelled successfully');
      loadOrders();
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to cancel order');
      setTimeout(() => setError(''), 3000);
    }
  };

  const getStatusBadge = (status: OrderStatus) => {
    const statusColors = {
      pending: 'bg-yellow-100 text-yellow-800',
      preparing: 'bg-blue-100 text-blue-800',
      ready: 'bg-green-100 text-green-800',
      picked_up: 'bg-gray-100 text-gray-800',
      cancelled: 'bg-red-100 text-red-800',
    };

    return (
      <span
        className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${statusColors[status]}`}
      >
        {status.replace('_', ' ').toUpperCase()}
      </span>
    );
  };

  const getNextStatus = (currentStatus: OrderStatus): OrderStatus | null => {
    const statusFlow: Record<OrderStatus, OrderStatus | null> = {
      pending: 'preparing',
      preparing: 'ready',
      ready: 'picked_up',
      picked_up: null,
      cancelled: null,
    };
    return statusFlow[currentStatus];
  };

  const getStatusButtonText = (status: OrderStatus): string => {
    const buttonText: Record<OrderStatus, string> = {
      pending: 'Start Preparing',
      preparing: 'Mark as Ready',
      ready: 'Mark as Picked Up',
      picked_up: '',
      cancelled: '',
    };
    return buttonText[status];
  };

  if (isLoading) {
    return <div className="p-8">Loading orders...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-6">
            <Link to="/admin" className="text-blue-600 hover:text-blue-500">
              ← Back to Admin Dashboard
            </Link>
          </div>

          <h1 className="text-3xl font-bold text-gray-900 mb-6">Order Management</h1>

          {error && (
            <div className="mb-4 p-4 bg-red-50 text-red-700 rounded-md">{error}</div>
          )}

          {success && (
            <div className="mb-4 p-4 bg-green-50 text-green-700 rounded-md">{success}</div>
          )}

          {/* Summary Stats */}
          <div className="grid grid-cols-1 gap-5 sm:grid-cols-5 mb-6">
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="px-4 py-5 sm:p-6">
                <dt className="text-sm font-medium text-gray-500 truncate">Pending</dt>
                <dd className="mt-1 text-3xl font-semibold text-yellow-600">
                  {orders.filter((o) => o.status === 'pending').length}
                </dd>
              </div>
            </div>
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="px-4 py-5 sm:p-6">
                <dt className="text-sm font-medium text-gray-500 truncate">Preparing</dt>
                <dd className="mt-1 text-3xl font-semibold text-blue-600">
                  {orders.filter((o) => o.status === 'preparing').length}
                </dd>
              </div>
            </div>
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="px-4 py-5 sm:p-6">
                <dt className="text-sm font-medium text-gray-500 truncate">Ready</dt>
                <dd className="mt-1 text-3xl font-semibold text-green-600">
                  {orders.filter((o) => o.status === 'ready').length}
                </dd>
              </div>
            </div>
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="px-4 py-5 sm:p-6">
                <dt className="text-sm font-medium text-gray-500 truncate">Picked Up</dt>
                <dd className="mt-1 text-3xl font-semibold text-gray-600">
                  {orders.filter((o) => o.status === 'picked_up').length}
                </dd>
              </div>
            </div>
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="px-4 py-5 sm:p-6">
                <dt className="text-sm font-medium text-gray-500 truncate">Cancelled</dt>
                <dd className="mt-1 text-3xl font-semibold text-red-600">
                  {orders.filter((o) => o.status === 'cancelled').length}
                </dd>
              </div>
            </div>
          </div>

          {/* Filter */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Filter by Status
            </label>
            <select
              className="w-full md:w-64 px-3 py-2 border border-gray-300 rounded-md"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value as OrderStatus | '')}
            >
              <option value="">All Orders</option>
              <option value="pending">Pending</option>
              <option value="preparing">Preparing</option>
              <option value="ready">Ready for Pickup</option>
              <option value="picked_up">Picked Up</option>
              <option value="cancelled">Cancelled</option>
            </select>
          </div>

          {/* Orders List */}
          {orders.length === 0 ? (
            <div className="bg-white shadow rounded-lg p-12 text-center">
              <p className="text-gray-500">No orders found</p>
            </div>
          ) : (
            <div className="bg-white shadow overflow-hidden sm:rounded-md">
              <ul className="divide-y divide-gray-200">
                {orders.map((order) => {
                  const nextStatus = getNextStatus(order.status);
                  const buttonText = getStatusButtonText(order.status);

                  return (
                    <li key={order.id} className="px-6 py-4 hover:bg-gray-50">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center justify-between mb-2">
                            <div>
                              <p className="text-sm font-medium text-blue-600 truncate">
                                Order #{order.id.slice(0, 8)}
                              </p>
                              <p className="text-sm text-gray-500">
                                Customer: {order.user?.first_name} {order.user?.last_name}
                              </p>
                              <p className="text-sm text-gray-500">
                                Email: {order.user?.email}
                              </p>
                            </div>
                            {getStatusBadge(order.status)}
                          </div>

                          <div className="mt-2 text-sm text-gray-500">
                            <p>
                              {order.cart?.items?.length || 0} item(s) | Submitted:{' '}
                              {new Date(order.submitted_at).toLocaleString()}
                            </p>
                          </div>

                          {order.notes && (
                            <p className="mt-2 text-sm text-gray-600">
                              <span className="font-medium">Note:</span> {order.notes}
                            </p>
                          )}

                          {order.ready_at && (
                            <p className="mt-2 text-sm text-green-600">
                              Ready at: {new Date(order.ready_at).toLocaleString()}
                            </p>
                          )}

                          {order.picked_up_at && (
                            <p className="mt-2 text-sm text-gray-600">
                              Picked up at: {new Date(order.picked_up_at).toLocaleString()}
                            </p>
                          )}

                          {/* Order Items */}
                          {order.cart?.items && order.cart.items.length > 0 && (
                            <div className="mt-3 pt-3 border-t border-gray-200">
                              <p className="text-sm font-medium text-gray-700 mb-2">
                                Items:
                              </p>
                              <ul className="space-y-1">
                                {order.cart.items.map((cartItem) => (
                                  <li key={cartItem.id} className="text-sm text-gray-600">
                                    {cartItem.item?.name} - Quantity: {cartItem.quantity}
                                  </li>
                                ))}
                              </ul>
                            </div>
                          )}
                        </div>

                        {/* Action buttons */}
                        <div className="ml-4 flex flex-col space-y-2">
                          {nextStatus && buttonText && (
                            <button
                              onClick={() => handleUpdateStatus(order.id, nextStatus)}
                              className="px-4 py-2 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700"
                            >
                              {buttonText}
                            </button>
                          )}
                          {(order.status === 'pending' ||
                            order.status === 'preparing' ||
                            order.status === 'ready') && (
                            <button
                              onClick={() => handleCancelOrder(order.id)}
                              className="px-4 py-2 text-sm text-red-600 border border-red-300 rounded-md hover:bg-red-50"
                            >
                              Cancel Order
                            </button>
                          )}
                        </div>
                      </div>
                    </li>
                  );
                })}
              </ul>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
