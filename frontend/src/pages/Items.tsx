import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { itemService } from '../services/itemService';
import { cartService } from '../services/cartService';
import { pantryService } from '../services/pantryService';
import type { Item, Pantry } from '../types';

export const ItemsPage = () => {
  const [items, setItems] = useState<Item[]>([]);
  const [selectedPantry, setSelectedPantry] = useState<Pantry | null>(null);
  const [search, setSearch] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    loadPantryAndItems();
  }, [search]);

  const loadPantryAndItems = async () => {
    try {
      setIsLoading(true);

      // Load selected pantry
      const pantryId = localStorage.getItem('selectedPantryId');
      if (pantryId) {
        const pantry = await pantryService.getPantry(pantryId);
        setSelectedPantry(pantry);
      }

      const response = await itemService.listPublic({ search, page: 1, page_size: 50 });
      setItems(response.data);
    } catch (err: any) {
      setError('Failed to load items');
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddToCart = async (itemId: string, itemName: string) => {
    try {
      setError('');
      setSuccess('');
      await cartService.addItem(itemId, 1);
      setSuccess(`Added "${itemName}" to cart!`);
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to add item to cart');
      setTimeout(() => setError(''), 3000);
    }
  };

  if (isLoading) {
    return <div className="p-8">Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-6 flex justify-between items-center">
            <Link to="/" className="text-blue-600 hover:text-blue-500">
              ← Back to Home
            </Link>
            <Link to="/pantries" className="text-blue-600 hover:text-blue-500">
              Change Pantry
            </Link>
          </div>

          {selectedPantry && (
            <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-md">
              <h2 className="text-sm font-medium text-blue-900 mb-1">
                Shopping at: {selectedPantry.name}
              </h2>
              <p className="text-sm text-blue-700">
                {selectedPantry.city}, {selectedPantry.state}
              </p>
            </div>
          )}

          <div className="flex justify-between items-center mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Browse Items</h1>
            <Link
              to="/cart"
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              View Cart
            </Link>
          </div>

          {error && (
            <div className="mb-4 p-4 bg-red-50 text-red-700 rounded-md">{error}</div>
          )}

          {success && (
            <div className="mb-4 p-4 bg-green-50 text-green-700 rounded-md">{success}</div>
          )}

          <div className="mb-6">
            <input
              type="text"
              placeholder="Search items..."
              className="w-full px-4 py-2 border border-gray-300 rounded-md"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>

          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {items.map((item) => (
              <div
                key={item.id}
                className="bg-white overflow-hidden shadow rounded-lg hover:shadow-md transition-shadow"
              >
                <div className="p-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-2">{item.name}</h3>
                  {item.description && (
                    <p className="text-sm text-gray-500 mb-4">{item.description}</p>
                  )}

                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <span className="text-sm text-gray-500">Category: </span>
                      <span className="text-sm font-medium text-gray-900">
                        {item.category?.name || 'N/A'}
                      </span>
                    </div>
                  </div>

                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <span className="text-sm text-gray-500">Available: </span>
                      <span className="text-sm font-medium text-gray-900">
                        {item.quantity} {item.unit}
                      </span>
                    </div>
                    <span
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        item.is_available
                          ? 'bg-green-100 text-green-800'
                          : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {item.is_available ? 'Available' : 'Unavailable'}
                    </span>
                  </div>

                  <button
                    onClick={() => handleAddToCart(item.id, item.name)}
                    disabled={!item.is_available || item.quantity === 0}
                    className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {item.is_available && item.quantity > 0
                      ? 'Add to Cart'
                      : 'Out of Stock'}
                  </button>
                </div>
              </div>
            ))}

            {items.length === 0 && (
              <div className="col-span-3 text-center py-12 text-gray-500">
                No items found. Try adjusting your search.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
