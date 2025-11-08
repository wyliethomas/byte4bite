import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { donationService } from '../../services/donationService';
import type { Donation } from '../../types';
import type { DonationStatsResponse } from '../../services/donationService';

export const AdminDonations = () => {
  const [donations, setDonations] = useState<Donation[]>([]);
  const [stats, setStats] = useState<DonationStatsResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [receiptFilter, setReceiptFilter] = useState<'all' | 'sent' | 'pending'>('all');

  useEffect(() => {
    loadData();
  }, [receiptFilter]);

  const loadData = async () => {
    try {
      setIsLoading(true);

      // Load donations
      const receiptSentParam =
        receiptFilter === 'sent' ? true : receiptFilter === 'pending' ? false : undefined;
      const donationsResponse = await donationService.getDonations({
        receipt_sent: receiptSentParam,
      });
      setDonations(donationsResponse.donations);

      // Load statistics
      const statsData = await donationService.getDonationStats();
      setStats(statsData);
    } catch (err: any) {
      setError('Failed to load donations');
    } finally {
      setIsLoading(false);
    }
  };

  const handleMarkReceiptSent = async (donationId: string) => {
    try {
      setError('');
      setSuccess('');
      await donationService.markReceiptSent(donationId);
      setSuccess('Receipt marked as sent');
      loadData();
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError('Failed to mark receipt as sent');
      setTimeout(() => setError(''), 3000);
    }
  };

  const handleDelete = async (donationId: string, donorName: string) => {
    if (!confirm(`Are you sure you want to delete the donation from "${donorName}"?`))
      return;

    try {
      setError('');
      setSuccess('');
      await donationService.deleteDonation(donationId);
      setSuccess('Donation deleted successfully');
      loadData();
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete donation');
      setTimeout(() => setError(''), 3000);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };

  if (isLoading) {
    return <div className="p-8">Loading donations...</div>;
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

          <h1 className="text-3xl font-bold text-gray-900 mb-6">
            Donation Management
          </h1>

          {error && (
            <div className="mb-4 p-4 bg-red-50 text-red-700 rounded-md">{error}</div>
          )}

          {success && (
            <div className="mb-4 p-4 bg-green-50 text-green-700 rounded-md">{success}</div>
          )}

          {/* Statistics */}
          {stats && (
            <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3 mb-6">
              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="px-4 py-5 sm:p-6">
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Total Donations
                  </dt>
                  <dd className="mt-1 text-3xl font-semibold text-gray-900">
                    {stats.total_donations}
                  </dd>
                  <div className="mt-2 text-sm text-gray-600">
                    {stats.monetary_count} monetary, {stats.in_kind_count} in-kind
                  </div>
                </div>
              </div>

              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="px-4 py-5 sm:p-6">
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Total Amount
                  </dt>
                  <dd className="mt-1 text-3xl font-semibold text-green-600">
                    {formatCurrency(stats.total_amount)}
                  </dd>
                  <div className="mt-2 text-sm text-gray-600">
                    Monetary donations only
                  </div>
                </div>
              </div>

              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="px-4 py-5 sm:p-6">
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Unique Donors
                  </dt>
                  <dd className="mt-1 text-3xl font-semibold text-blue-600">
                    {stats.donor_count}
                  </dd>
                </div>
              </div>

              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="px-4 py-5 sm:p-6">
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Receipts Pending
                  </dt>
                  <dd className="mt-1 text-3xl font-semibold text-yellow-600">
                    {stats.receipts_pending}
                  </dd>
                </div>
              </div>
            </div>
          )}

          {/* Filters */}
          <div className="mb-6 bg-white shadow rounded-lg p-4">
            <div className="flex items-center space-x-4">
              <label className="text-sm font-medium text-gray-700">Filter:</label>
              <select
                className="px-3 py-2 border border-gray-300 rounded-md"
                value={receiptFilter}
                onChange={(e) =>
                  setReceiptFilter(e.target.value as 'all' | 'sent' | 'pending')
                }
              >
                <option value="all">All Donations</option>
                <option value="pending">Receipt Pending</option>
                <option value="sent">Receipt Sent</option>
              </select>
            </div>
          </div>

          {/* Donations List */}
          {donations.length === 0 ? (
            <div className="bg-white shadow rounded-lg p-12 text-center">
              <p className="text-gray-500">No donations found</p>
            </div>
          ) : (
            <div className="bg-white shadow overflow-hidden sm:rounded-lg">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Date
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Donor
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Type
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Amount/Description
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Pantry
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Receipt
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {donations.map((donation) => (
                    <tr key={donation.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {formatDate(donation.donation_date)}
                      </td>
                      <td className="px-6 py-4">
                        <div className="text-sm font-medium text-gray-900">
                          {donation.donor_name}
                        </div>
                        {donation.donor_email && (
                          <div className="text-sm text-gray-500">
                            {donation.donor_email}
                          </div>
                        )}
                        {donation.donor_phone && (
                          <div className="text-sm text-gray-500">
                            {donation.donor_phone}
                          </div>
                        )}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                            donation.amount
                              ? 'bg-green-100 text-green-800'
                              : 'bg-blue-100 text-blue-800'
                          }`}
                        >
                          {donation.amount ? 'Monetary' : 'In-Kind'}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        {donation.amount ? (
                          <div className="text-sm font-medium text-green-600">
                            {formatCurrency(donation.amount)}
                          </div>
                        ) : null}
                        <div className="text-sm text-gray-600">
                          {donation.description}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {donation.pantry?.name || 'N/A'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {donation.receipt_sent ? (
                          <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                            Sent
                          </span>
                        ) : (
                          <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-yellow-100 text-yellow-800">
                            Pending
                          </span>
                        )}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        {!donation.receipt_sent && (
                          <button
                            onClick={() => handleMarkReceiptSent(donation.id)}
                            className="text-green-600 hover:text-green-900 mr-4"
                          >
                            Mark Sent
                          </button>
                        )}
                        <button
                          onClick={() => handleDelete(donation.id, donation.donor_name)}
                          className="text-red-600 hover:text-red-900"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
