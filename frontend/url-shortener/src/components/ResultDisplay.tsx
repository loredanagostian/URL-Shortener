import React, { useState, useEffect } from 'react';
import QRCode from 'qrcode';
import type { ShortenResponse } from '../services/api';
import './ResultDisplay.css';

interface ResultDisplayProps {
    result: ShortenResponse;
    onReset: () => void;
}

const ResultDisplay: React.FC<ResultDisplayProps> = ({ result, onReset }) => {
    const [copied, setCopied] = useState(false);
    const [qrCodeDataUrl, setQrCodeDataUrl] = useState<string>('');

    // Generate QR code when component mounts or result changes
    useEffect(() => {
        const generateQRCode = async () => {
            try {
                const dataUrl = await QRCode.toDataURL(result.short_url, {
                    width: 200,
                    margin: 1,
                    color: {
                        dark: '#000000',
                        light: '#FFFFFF',
                    },
                });
                setQrCodeDataUrl(dataUrl);
            } catch (error) {
                console.error('Error generating QR code:', error);
            }
        };

        generateQRCode();
    }, [result.short_url]);

    const copyToClipboard = async () => {
        try {
            await navigator.clipboard.writeText(result.short_url);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        } catch (err) {
            console.error('Failed to copy to clipboard:', err);
        }
    };

    const openInNewTab = () => {
        window.open(result.short_url, '_blank', 'noopener,noreferrer');
    };

    const downloadQRCode = () => {
        if (!qrCodeDataUrl) return;

        const link = document.createElement('a');
        link.download = `qr-code-${result.code}.png`;
        link.href = qrCodeDataUrl;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    const formatExpirationDate = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleString([], {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    return (
        <div className="result-display">
            <div className="result-header">
                <span className="success-icon">✅</span>
                <h2>URL Shortened Successfully!</h2>
            </div>

            <div className="result-content">
                <div className="url-section">
                    <label>Short URL:</label>
                    <div className="url-display">
                        <input
                            type="text"
                            value={result.short_url}
                            readOnly
                            className="short-url-input"
                        />
                        <button
                            onClick={copyToClipboard}
                            className={`copy-btn ${copied ? 'copied' : ''}`}
                            title="Copy to clipboard"
                        >
                            {copied ? '✓' : '📋'}
                        </button>
                    </div>
                </div>

                <div className="url-section">
                    <label>Original URL:</label>
                    <div className="original-url">{result.original_url}</div>
                </div>

                <div className="url-section">
                    <label>Short Code:</label>
                    <div className="short-code-container">
                        <div className="short-code">{result.code}</div>
                        <div className="qr-code-section">
                            {qrCodeDataUrl && (
                                <div className="qr-code-container">
                                    <img
                                        src={qrCodeDataUrl}
                                        alt="QR Code for short URL"
                                        className="qr-code-image"
                                    />
                                    <button
                                        onClick={downloadQRCode}
                                        className="download-qr-btn"
                                        title="Download QR Code"
                                    >
                                        ⬇️ Download QR Code
                                    </button>
                                </div>
                            )}
                        </div>
                    </div>
                </div>

                <div className="url-section">
                    <label>Expiration Date:</label>
                    <div className="expiration-value">
                        {result.expires_at ? formatExpirationDate(result.expires_at) : 'Unlimited'}
                    </div>
                </div>
            </div>

            <div className="result-actions">
                <button onClick={openInNewTab} className="test-btn">
                    Test Link
                </button>
                <button onClick={onReset} className="new-url-btn">
                    Shorten Another URL
                </button>
            </div>

            {copied && <div className="copy-notification">URL copied to clipboard! 📋</div>}
        </div>
    );
};

export default ResultDisplay;
