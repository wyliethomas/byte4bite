import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { pantryService } from '../services/pantryService';
import type { Pantry } from '../types';

export const PantriesPage = () => {
  const navigate = useNavigate();
  const [pantries, setPantries] = useState<Pantry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    loadPantries();
  }, []);

  const loadPantries = async () => {
    try {
      setIsLoading(true);
      const response = await pantryService.getPantries({ is_active: true });
      setPantries(response.pantries);
    } catch (err: any) {
      setError('Failed to load pantries');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      loadPantries();
      return;
    }

    try {
      setIsLoading(true);
      setError('');
      const results = await pantryService.searchPantries(searchQuery);
      setPantries(results);
    } catch (err: any) {
      setError('Failed to search pantries');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSelectPantry = (pantryId: string) => {
    // Store selected pantry in localStorage
    localStorage.setItem('selectedPantryId', pantryId);
    // Navigate to items page
    navigate('/items');
  };

  if (isLoading) {
    return <div className="p-8">Loading pantries...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-6">
            <Link to="/" className="text-blue-600 hover:text-blue-500">
              ← Back to Home
            </Link>
          </div>

          <h1 className="text-3xl font-bold text-gray-900 mb-2">Select a Pantry</h1>
          <p className="text-gray-600 mb-6">
            Choose a community pantry to browse available items
          </p>

          {error && (
            <div className="mb-4 p-4 bg-red-50 text-red-700 rounded-md">{error}</div>
          )}

          {/* Search */}
          <div className="mb-6">
            <div className="flex gap-2">
              <input
                type="text"
                placeholder="Search by name or city..."
                className="flex-1 px-4 py-2 border border-gray-300 rounded-md"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              />
              <button
                onClick={handleSearch}
                className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                Search
              </button>
              {searchQuery && (
                <button
                  onClick={() => {
                    setSearchQuery('');
                    loadPantries();
                  }}
                  className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  Clear
                </button>
              )}
            </div>
          </div>

          {pantries.length === 0 ? (
            <div className="bg-white shadow rounded-lg p-12 text-center">
              <p className="text-gray-500">
                {searchQuery
                  ? 'No pantries found matching your search'
                  : 'No active pantries available'}
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {pantries.map((pantry) => (
                <div
                  key={pantry.id}
                  className="bg-white overflow-hidden shadow rounded-lg hover:shadow-lg transition-shadow cursor-pointer"
                  onClick={() => handleSelectPantry(pantry.id)}
                >
                  <div className="p-6">
                    <h3 className="text-xl font-semibold text-gray-900 mb-2">
                      {pantry.name}
                    </h3>

                    <div className="space-y-2 text-sm text-gray-600 mb-4">
                      <p className="flex items-start">
                        <svg
                          className="h-5 w-5 mr-2 text-gray-400"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                          />
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                          />
                        </svg>
                        <span>
                          {pantry.address}
                          <br />
                          {pantry.city}, {pantry.state} {pantry.zip_code}
                        </span>
                      </p>

                      <p className="flex items-center">
                        <svg
                          className="h-5 w-5 mr-2 text-gray-400"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                          />
                        </svg>
                        {pantry.contact_email}
                      </p>

                      {pantry.contact_phone && (
                        <p className="flex items-center">
                          <svg
                            className="h-5 w-5 mr-2 text-gray-400"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"
                            />
                          </svg>
                          {pantry.contact_phone}
                        </p>
                      )}
                    </div>

                    <button className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors">
                      Browse Items
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
