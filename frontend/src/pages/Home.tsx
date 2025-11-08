import { useAuth } from '../context/AuthContext';
import { Link } from 'react-router-dom';

export const Home = () => {
  const { user, logout, isAuthenticated } = useAuth();

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold text-gray-900">Byte4Bite</h1>
            </div>
            <div className="flex items-center space-x-4">
              {isAuthenticated ? (
                <>
                  <span className="text-gray-700">
                    Welcome, {user?.first_name}!
                  </span>
                  {user?.role === 'admin' && (
                    <Link
                      to="/admin"
                      className="text-blue-600 hover:text-blue-500"
                    >
                      Admin Dashboard
                    </Link>
                  )}
                  <Link
                    to="/profile"
                    className="text-blue-600 hover:text-blue-500"
                  >
                    Profile
                  </Link>
                  <button
                    onClick={logout}
                    className="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700"
                  >
                    Logout
                  </button>
                </>
              ) : (
                <>
                  <Link
                    to="/login"
                    className="text-blue-600 hover:text-blue-500"
                  >
                    Sign in
                  </Link>
                  <Link
                    to="/register"
                    className="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
                  >
                    Register
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="text-center">
            <h2 className="text-3xl font-extrabold text-gray-900 sm:text-4xl">
              Welcome to Byte4Bite
            </h2>
            <p className="mt-4 text-lg text-gray-500">
              A free and open platform for community pantries
            </p>

            {isAuthenticated ? (
              <div className="mt-10">
                <div className="rounded-lg bg-white shadow p-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-4">
                    Your Dashboard
                  </h3>
                  <p className="text-gray-600 mb-4">
                    User Role: <span className="font-semibold capitalize">{user?.role}</span>
                  </p>
                  <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
                    <Link
                      to="/items"
                      className="p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 transition-colors"
                    >
                      <h4 className="font-medium text-gray-900">Browse Items</h4>
                      <p className="text-sm text-gray-500 mt-1">
                        View available pantry items
                      </p>
                    </Link>
                    <Link
                      to="/cart"
                      className="p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 transition-colors"
                    >
                      <h4 className="font-medium text-gray-900">Your Cart</h4>
                      <p className="text-sm text-gray-500 mt-1">
                        View and manage your cart
                      </p>
                    </Link>
                    <Link
                      to="/orders"
                      className="p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 transition-colors"
                    >
                      <h4 className="font-medium text-gray-900">Your Orders</h4>
                      <p className="text-sm text-gray-500 mt-1">
                        Track your order history
                      </p>
                    </Link>
                  </div>
                </div>
              </div>
            ) : (
              <div className="mt-10">
                <p className="text-gray-600 mb-6">
                  Please sign in or create an account to get started
                </p>
                <div className="space-x-4">
                  <Link
                    to="/login"
                    className="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
                  >
                    Sign in
                  </Link>
                  <Link
                    to="/register"
                    className="inline-flex items-center px-6 py-3 border border-gray-300 text-base font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                  >
                    Create account
                  </Link>
                </div>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};
