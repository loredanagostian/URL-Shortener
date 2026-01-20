import React, { useState, useEffect } from 'react';
import QRCode from 'qrcode';
import { apiService } from '../services/api';
import type { URLHistoryItem } from '../services/api';
import './History.css';

const History: React.FC = () => {
    const [history, setHistory] = useState<URLHistoryItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        loadHistory();
    }, []);

    const loadHistory = async () => {
        try {
            setLoading(true);
            setError(null);
            const historyData = await apiService.getURLHistory();
            setHistory(historyData);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load history');
        } finally {
            setLoading(false);
        }
    };

    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return (
            date.toLocaleDateString() +
            ' ' +
            date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
        );
    };

    const copyToClipboard = async (text: string) => {
        try {
            await navigator.clipboard.writeText(text);
        } catch (err) {
            console.error('Failed to copy:', err);
        }
    };

    const downloadQRCode = async (shortCode: string) => {
        try {
            const shortUrl = apiService.getShortURL(shortCode);
            const dataUrl = await QRCode.toDataURL(shortUrl, {
                width: 200,
                margin: 1,
                color: {
                    dark: '#000000',
                    light: '#FFFFFF',
                },
            });

            const link = document.createElement('a');
            link.download = `qr-code-${shortCode}.png`;
            link.href = dataUrl;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        } catch (error) {
            console.error('Error generating QR code:', error);
        }
    };

    const deleteURL = async (shortCode: string) => {
        if (
            window.confirm(
                'Are you sure you want to delete this URL? This action cannot be undone.',
            )
        ) {
            try {
                await apiService.deleteURL(shortCode);
                // Refresh the history after deletion
                await loadHistory();
            } catch (error) {
                console.error('Error deleting URL:', error);
                setError(error instanceof Error ? error.message : 'Failed to delete URL');
            }
        }
    };

    const getStatusBadge = (status: string) => {
        const className = status === 'active' ? 'status-active' : 'status-expired';
        const icon = status === 'active' ? '✅' : '⏰';
        return (
            <span className={`status-badge ${className}`}>
                {icon} {status.charAt(0).toUpperCase() + status.slice(1)}
            </span>
        );
    };

    if (loading) {
        return (
            <div className="history-container">
                <div className="history-header">
                    <h2>📊 URL History</h2>
                </div>
                <div className="loading">Loading history...</div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="history-container">
                <div className="history-header">
                    <h2>📊 URL History</h2>
                </div>
                <div className="error-message">
                    <span className="error-icon">⚠️</span>
                    {error}
                </div>
            </div>
        );
    }

    return (
        <div className="history-container">
            <div className="history-header">
                <h2>📊 URL History</h2>
                <button onClick={loadHistory} className="refresh-btn">
                    🔄 Refresh
                </button>
            </div>

            {history.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-icon">🔗</div>
                    <h3>No URLs created yet</h3>
                    <p>Start by creating your first short URL!</p>
                </div>
            ) : (
                <div className="history-table-container">
                    <table className="history-table">
                        <thead>
                            <tr>
                                <th>Original URL</th>
                                <th>Short URL</th>
                                <th>Code</th>
                                <th>Created</th>
                                <th>Clicks</th>
                                <th>Status</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {history.map((item) => (
                                <tr key={item.id}>
                                    <td className="original-url-cell">
                                        <div className="url-text" title={item.original_url}>
                                            {item.original_url.length > 40
                                                ? item.original_url.substring(0, 40) + '...'
                                                : item.original_url}
                                        </div>
                                    </td>
                                    <td className="short-url-cell">
                                        <div className="short-url-text">
                                            {apiService.getShortURL(item.short_code)}
                                        </div>
                                    </td>
                                    <td className="code-cell">
                                        <span className="code-badge">{item.short_code}</span>
                                    </td>
                                    <td className="date-cell">{formatDate(item.created_at)}</td>
                                    <td className="clicks-cell">
                                        <span className="clicks-count">{item.click_count}</span>
                                    </td>
                                    <td className="status-cell">{getStatusBadge(item.status)}</td>
                                    <td className="actions-cell">
                                        <button
                                            onClick={() =>
                                                copyToClipboard(
                                                    apiService.getShortURL(item.short_code),
                                                )
                                            }
                                            className="copy-action-btn"
                                            title="Copy short URL"
                                        >
                                            📋
                                        </button>
                                        <button
                                            onClick={() => downloadQRCode(item.short_code)}
                                            className="qr-action-btn"
                                            title="Download QR Code"
                                        >
                                            📱
                                        </button>
                                        <button
                                            onClick={() => deleteURL(item.short_code)}
                                            className="delete-action-btn"
                                            title="Delete URL"
                                        >
                                            🗑️
                                        </button>
                                        {item.status === 'active' && (
                                            <button
                                                onClick={() =>
                                                    window.open(
                                                        apiService.getShortURL(item.short_code),
                                                        '_blank',
                                                    )
                                                }
                                                className="open-action-btn"
                                                title="Open short URL"
                                            >
                                                🔗
                                            </button>
                                        )}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
};

export default History;
