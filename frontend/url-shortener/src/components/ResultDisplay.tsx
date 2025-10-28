import React, { useState } from "react";
import type { ShortenResponse } from "../services/api";
import "./ResultDisplay.css";

interface ResultDisplayProps {
  result: ShortenResponse;
  onReset: () => void;
}

const ResultDisplay: React.FC<ResultDisplayProps> = ({ result, onReset }) => {
  const [copied, setCopied] = useState(false);

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(result.short_url);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy to clipboard:", err);
    }
  };

  const openInNewTab = () => {
    window.open(result.short_url, "_blank", "noopener,noreferrer");
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
              className={`copy-btn ${copied ? "copied" : ""}`}
              title="Copy to clipboard"
            >
              {copied ? "✓" : "📋"}
            </button>
          </div>
        </div>

        <div className="url-section">
          <label>Original URL:</label>
          <div className="original-url">{result.original_url}</div>
        </div>

        <div className="url-section">
          <label>Short Code:</label>
          <div className="short-code">{result.code}</div>
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

      {copied && (
        <div className="copy-notification">URL copied to clipboard! 📋</div>
      )}
    </div>
  );
};

export default ResultDisplay;
