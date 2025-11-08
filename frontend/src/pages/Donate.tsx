import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { pantryService } from '../services/pantryService';
import { donationService } from '../services/donationService';
import type { Pantry } from '../types';
import type { CreateDonationRequest } from '../services/donationService';

export const DonatePage = () => {
  const navigate = useNavigate();
  const [pantries, setPantries] = useState<Pantry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [formData, setFormData] = useState<CreateDonationRequest>({
    pantry_id: '',
    donor_name: '',
    donor_email: '',
    donor_phone: '',
    amount: undefined,
    description: '',
  });
  const [donationType, setDonationType] = useState<'monetary' | 'in-kind'>('monetary');

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!formData.pantry_id) {
      setError('Please select a pantry');
      return;
    }

    try {
      setIsSubmitting(true);
      await donationService.createDonation(formData);
      setSuccess(true);
      // Reset form
      setFormData({
        pantry_id: '',
        donor_name: '',
        donor_email: '',
        donor_phone: '',
        amount: undefined,
        description: '',
      });
      setDonationType('monetary');

      // Show success message for 3 seconds then redirect
      setTimeout(() => {
        navigate('/');
      }, 3000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to submit donation');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return <div className="p-8">Loading...</div>;
  }

  if (success) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full text-center">
          <div className="mb-4">
            <svg
              className="mx-auto h-16 w-16 text-green-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Thank You!</h2>
          <p className="text-gray-600 mb-4">
            Your donation has been successfully submitted. We greatly appreciate your
            generosity!
          </p>
          <Link to="/" className="text-blue-600 hover:text-blue-500">
            Return to Home
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-6">
            <Link to="/" className="text-blue-600 hover:text-blue-500">
              ← Back to Home
            </Link>
          </div>

          <div className="bg-white shadow rounded-lg p-6">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">Make a Donation</h1>
            <p className="text-gray-600 mb-6">
              Thank you for your generosity! Your donation helps us serve our community
              better.
            </p>

            {error && (
              <div className="mb-4 p-4 bg-red-50 text-red-700 rounded-md">{error}</div>
            )}

            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Donation Type Selection */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Donation Type *
                </label>
                <div className="grid grid-cols-2 gap-4">
                  <button
                    type="button"
                    onClick={() => {
                      setDonationType('monetary');
                      setFormData({ ...formData, amount: undefined });
                    }}
                    className={`p-4 border-2 rounded-md ${
                      donationType === 'monetary'
                        ? 'border-blue-600 bg-blue-50'
                        : 'border-gray-300'
                    }`}
                  >
                    <div className="text-center">
                      <div className="text-2xl mb-2">💵</div>
                      <div className="font-medium">Monetary Donation</div>
                    </div>
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setDonationType('in-kind');
                      setFormData({ ...formData, amount: undefined });
                    }}
                    className={`p-4 border-2 rounded-md ${
                      donationType === 'in-kind'
                        ? 'border-blue-600 bg-blue-50'
                        : 'border-gray-300'
                    }`}
                  >
                    <div className="text-center">
                      <div className="text-2xl mb-2">📦</div>
                      <div className="font-medium">In-Kind Donation</div>
                    </div>
                  </button>
                </div>
              </div>

              {/* Pantry Selection */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Select Pantry *
                </label>
                <select
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  value={formData.pantry_id}
                  onChange={(e) =>
                    setFormData({ ...formData, pantry_id: e.target.value })
                  }
                >
                  <option value="">Choose a pantry...</option>
                  {pantries.map((pantry) => (
                    <option key={pantry.id} value={pantry.id}>
                      {pantry.name} - {pantry.city}, {pantry.state}
                    </option>
                  ))}
                </select>
              </div>

              {/* Donor Information */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Your Name *
                </label>
                <input
                  type="text"
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  value={formData.donor_name}
                  onChange={(e) =>
                    setFormData({ ...formData, donor_name: e.target.value })
                  }
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Email
                  </label>
                  <input
                    type="email"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                    value={formData.donor_email}
                    onChange={(e) =>
                      setFormData({ ...formData, donor_email: e.target.value })
                    }
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Phone
                  </label>
                  <input
                    type="tel"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                    value={formData.donor_phone}
                    onChange={(e) =>
                      setFormData({ ...formData, donor_phone: e.target.value })
                    }
                  />
                </div>
              </div>

              {/* Amount (for monetary donations) */}
              {donationType === 'monetary' && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Amount ($)
                  </label>
                  <input
                    type="number"
                    min="0"
                    step="0.01"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                    value={formData.amount || ''}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        amount: e.target.value ? parseFloat(e.target.value) : undefined,
                      })
                    }
                  />
                </div>
              )}

              {/* Description */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  {donationType === 'monetary'
                    ? 'Description (optional)'
                    : 'What are you donating? *'}
                </label>
                <textarea
                  required={donationType === 'in-kind'}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  rows={4}
                  placeholder={
                    donationType === 'monetary'
                      ? 'Any additional notes...'
                      : 'Please describe the items you are donating (e.g., 20 cans of vegetables, 10 boxes of pasta...)'
                  }
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                />
              </div>

              <div className="bg-blue-50 border border-blue-200 rounded-md p-4">
                <p className="text-sm text-blue-800">
                  <strong>Note:</strong> After submitting this form, our team will reach
                  out to you to arrange the donation details. Thank you for your support!
                </p>
              </div>

              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? 'Submitting...' : 'Submit Donation'}
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};
