import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { itemService } from '../../services/itemService';
import type { Item } from '../../types';

export const AdminItems = () => {
  const [items, setItems] = useState<Item[]>([]);
  const [search, setSearch] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadItems();
  }, [search]);

  const loadItems = async () => {
    try {
      setIsLoading(true);
      const response = await itemService.list({ search, page: 1, page_size: 50 });
      setItems(response.data);
    } catch (err: any) {
      setError('Failed to load items');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this item?')) return;

    try {
      await itemService.delete(id);
      loadItems();
    } catch (err: any) {
      setError('Failed to delete item');
    }
  };

  const handleQuantityUpdate = async (id: string, currentQty: number) => {
    const newQty = prompt(`Update quantity (current: ${currentQty}):`, currentQty.toString());
    if (newQty === null) return;

    try {
      await itemService.updateQuantity(id, parseInt(newQty));
      loadItems();
    } catch (err: any) {
      setError('Failed to update quantity');
    }
  };

  if (isLoading) {
    return <div className="p-8">Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-6">
            <Link to="/admin" className="text-blue-600 hover:text-blue-500">
              ← Back to Dashboard
            </Link>
          </div>

          <div className="mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Inventory Items</h1>
          </div>

          {error && (
            <div className="mb-4 p-4 bg-red-50 text-red-700 rounded-md">{error}</div>
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

          <div className="bg-white shadow overflow-hidden sm:rounded-lg">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Category
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Quantity
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Unit
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {items.map((item) => (
                  <tr key={item.id} className={item.quantity <= item.low_stock_threshold ? 'bg-yellow-50' : ''}>
                    <td className="px-6 py-4">
                      <div className="text-sm font-medium text-gray-900">{item.name}</div>
                      {item.description && (
                        <div className="text-sm text-gray-500">{item.description}</div>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {item.category?.name || 'N/A'}
                    </td>
                    <td className="px-6 py-4">
                      <span className={`text-sm font-medium ${
                        item.quantity <= item.low_stock_threshold
                          ? 'text-yellow-600'
                          : 'text-gray-900'
                      }`}>
                        {item.quantity}
                      </span>
                      {item.quantity <= item.low_stock_threshold && (
                        <span className="ml-2 text-xs text-yellow-600">(Low)</span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">{item.unit}</td>
                    <td className="px-6 py-4">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        item.is_available
                          ? 'bg-green-100 text-green-800'
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {item.is_available ? 'Available' : 'Unavailable'}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm space-x-2">
                      <button
                        onClick={() => handleQuantityUpdate(item.id, item.quantity)}
                        className="text-blue-600 hover:text-blue-800"
                      >
                        Update Qty
                      </button>
                      <button
                        onClick={() => handleDelete(item.id)}
                        className="text-red-600 hover:text-red-800"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                ))}
                {items.length === 0 && (
                  <tr>
                    <td colSpan={6} className="px-6 py-8 text-center text-gray-500">
                      No items found. Try adjusting your search.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          <div className="mt-4 text-sm text-gray-500">
            Showing {items.length} items
          </div>
        </div>
      </div>
    </div>
  );
};
