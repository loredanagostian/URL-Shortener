import { useState } from "react";
import { apiService } from "../services/api";
import type { ShortenRequest, ShortenResponse } from "../services/api";

export const useURLShortener = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<ShortenResponse | null>(null);

  const shortenURL = async (request: ShortenRequest) => {
    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await apiService.createShortURL(request);
      setResult(response);
      return response;
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "An unknown error occurred";
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const reset = () => {
    setError(null);
    setResult(null);
  };

  return {
    loading,
    error,
    result,
    shortenURL,
    reset,
  };
};
