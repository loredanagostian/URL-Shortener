import React, { useState } from 'react';
import { useURLShortener } from '../hooks/useURLShortener';
import type { ShortenRequest, ShortenResponse } from '../services/api';
import './URLForm.css';

interface URLFormProps {
    onSuccess?: (result: ShortenResponse) => void;
}

const URLForm: React.FC<URLFormProps> = ({ onSuccess }) => {
    const [url, setUrl] = useState('');
    const [customCode, setCustomCode] = useState('');
    const [expirationDate, setExpirationDate] = useState('');
    const [showAdvanced, setShowAdvanced] = useState(false);

    const { loading, error, shortenURL } = useURLShortener();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        // Validate URL format
        const urlToShorten = url.trim();
        if (!urlToShorten) {
            return;
        }

        // Ensure URL has protocol
        let validUrl = urlToShorten;
        if (!validUrl.startsWith('http://') && !validUrl.startsWith('https://')) {
            validUrl = 'https://' + validUrl;
        }

        const request: ShortenRequest = {
            url: validUrl,
        };

        if (customCode.trim()) {
            request.custom_code = customCode.trim();
        }

        if (expirationDate) {
            // Convert date to ISO string format
            const expiry = new Date(expirationDate);
            if (expiry > new Date()) {
                request.expires_at = expiry.toISOString();
            }
        }

        console.log('Submitting request:', request);

        try {
            const result = await shortenURL(request);
            console.log('URL shortening successful:', result);

            // Clear form on success
            setUrl('');
            setCustomCode('');
            setExpirationDate('');
            setShowAdvanced(false);

            // Call success callback
            if (onSuccess) {
                onSuccess(result);
            }
        } catch (err) {
            // Error is handled by the hook
            console.error('Failed to shorten URL:', err);
        }
    };

    return (
        <div className="url-form-container">
            <form onSubmit={handleSubmit} className="url-form">
                <div className="form-group">
                    <input
                        type="url"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        placeholder="Enter your URL here (e.g., https://example.com)..."
                        className="url-input"
                        required
                        disabled={loading}
                    />
                </div>

                {showAdvanced && (
                    <div className="advanced-options">
                        <div className="form-group">
                            <input
                                type="text"
                                value={customCode}
                                onChange={(e) => setCustomCode(e.target.value)}
                                placeholder="Custom short code (optional)"
                                className="custom-code-input"
                                pattern="[a-zA-Z0-9-]{3,20}"
                                title="3-20 characters, letters, numbers, and hyphens only"
                                disabled={loading}
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="expiration-date" className="expiration-label">
                                Expiration Date (optional)
                            </label>
                            <input
                                type="datetime-local"
                                id="expiration-date"
                                value={expirationDate}
                                onChange={(e) => setExpirationDate(e.target.value)}
                                min={new Date().toISOString().slice(0, 16)}
                                className="expiration-input"
                                disabled={loading}
                            />
                        </div>
                    </div>
                )}

                <div className="form-actions">
                    <button
                        type="button"
                        onClick={() => setShowAdvanced(!showAdvanced)}
                        className="advanced-toggle"
                        disabled={loading}
                    >
                        {showAdvanced ? 'Hide' : 'Show'} Advanced Options
                    </button>

                    <button type="submit" className="shorten-btn" disabled={loading || !url.trim()}>
                        {loading ? 'Shortening...' : 'Shorten URL'}
                    </button>
                </div>

                {error && (
                    <div className="error-message">
                        <span className="error-icon">⚠️</span>
                        {error}
                    </div>
                )}
            </form>
        </div>
    );
};

export default URLForm;
