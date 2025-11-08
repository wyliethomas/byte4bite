import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { cartService, type CartResponse } from '../services/cartService';

export const CartPage = () => {
  const navigate = useNavigate();
  const [cartData, setCartData] = useState<CartResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [notes, setNotes] = useState('');
  const [isCheckingOut, setIsCheckingOut] = useState(false);

  useEffect(() => {
    loadCart();
  }, []);

  const loadCart = async () => {
    try {
      setIsLoading(true);
      const data = await cartService.getCurrentCart();
      setCartData(data);
    } catch (err: any) {
      setError('Failed to load cart');
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateQuantity = async (cartItemId: string, newQuantity: number) => {
    if (newQuantity < 0) return;

    try {
      setError('');
      await cartService.updateItemQuantity(cartItemId, newQuantity);
      loadCart();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to update quantity');
    }
  };

  const handleRemoveItem = async (cartItemId: string) => {
    if (!confirm('Remove this item from cart?')) return;

    try {
      setError('');
      await cartService.removeItem(cartItemId);
      loadCart();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to remove item');
    }
  };

  const handleClearCart = async () => {
    if (!confirm('Clear all items from cart?')) return;

    try {
      setError('');
      await cartService.clearCart();
      loadCart();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to clear cart');
    }
  };

  const handleCheckout = async () => {
    try {
      setError('');
      setIsCheckingOut(true);
      const result = await cartService.checkout(notes);
      alert('Order submitted successfully! Order ID: ' + result.order.id);
      navigate('/orders');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to checkout');
    } finally {
      setIsCheckingOut(false);
    }
  };

  if (isLoading) {
    return <div className="p-8">Loading cart...</div>;
  }

  const hasItems = cartData && cartData.cart && cartData.items.length > 0;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-6 flex justify-between items-center">
            <Link to="/items" className="text-blue-600 hover:text-blue-500">
              ← Continue Shopping
            </Link>
            <Link to="/" className="text-blue-600 hover:text-blue-500">
              Home
            </Link>
          </div>

          <h1 className="text-3xl font-bold text-gray-900 mb-6">Your Cart</h1>

          {error && (
            <div className="mb-4 p-4 bg-red-50 text-red-700 rounded-md">{error}</div>
          )}

          {!hasItems ? (
            <div className="bg-white shadow rounded-lg p-12 text-center">
              <p className="text-gray-500 mb-4">Your cart is empty</p>
              <Link
                to="/items"
                className="inline-block px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                Browse Items
              </Link>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Cart Items */}
              <div className="lg:col-span-2">
                <div className="bg-white shadow rounded-lg overflow-hidden">
                  <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
                    <h2 className="text-lg font-medium">Items ({cartData.count})</h2>
                    <button
                      onClick={handleClearCart}
                      className="text-sm text-red-600 hover:text-red-800"
                    >
                      Clear Cart
                    </button>
                  </div>

                  <ul className="divide-y divide-gray-200">
                    {cartData.items.map((cartItem) => (
                      <li key={cartItem.id} className="px-6 py-4">
                        <div className="flex items-center justify-between">
                          <div className="flex-1">
                            <h3 className="text-lg font-medium text-gray-900">
                              {cartItem.item?.name}
                            </h3>
                            {cartItem.item?.description && (
                              <p className="text-sm text-gray-500 mt-1">
                                {cartItem.item.description}
                              </p>
                            )}
                            <p className="text-sm text-gray-500 mt-1">
                              Category: {cartItem.item?.category?.name || 'N/A'}
                            </p>
                          </div>

                          <div className="flex items-center space-x-4">
                            <div className="flex items-center space-x-2">
                              <button
                                onClick={() =>
                                  handleUpdateQuantity(cartItem.id, cartItem.quantity - 1)
                                }
                                className="px-2 py-1 border border-gray-300 rounded hover:bg-gray-100"
                              >
                                -
                              </button>
                              <span className="text-gray-900 font-medium w-12 text-center">
                                {cartItem.quantity}
                              </span>
                              <button
                                onClick={() =>
                                  handleUpdateQuantity(cartItem.id, cartItem.quantity + 1)
                                }
                                className="px-2 py-1 border border-gray-300 rounded hover:bg-gray-100"
                              >
                                +
                              </button>
                            </div>

                            <button
                              onClick={() => handleRemoveItem(cartItem.id)}
                              className="text-red-600 hover:text-red-800"
                            >
                              Remove
                            </button>
                          </div>
                        </div>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>

              {/* Checkout Summary */}
              <div className="lg:col-span-1">
                <div className="bg-white shadow rounded-lg p-6">
                  <h2 className="text-lg font-medium mb-4">Order Summary</h2>

                  <div className="mb-4">
                    <p className="text-sm text-gray-600">Total Items:</p>
                    <p className="text-lg font-medium">{cartData.count}</p>
                  </div>

                  <div className="mb-6">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Notes (optional)
                    </label>
                    <textarea
                      className="w-full px-3 py-2 border border-gray-300 rounded-md"
                      rows={3}
                      placeholder="Any special requests or notes..."
                      value={notes}
                      onChange={(e) => setNotes(e.target.value)}
                    />
                  </div>

                  <button
                    onClick={handleCheckout}
                    disabled={isCheckingOut}
                    className="w-full px-6 py-3 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {isCheckingOut ? 'Processing...' : 'Checkout'}
                  </button>

                  <p className="text-xs text-gray-500 mt-4 text-center">
                    Your order will be prepared for pickup at the pantry
                  </p>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
