// API service for communicating with the backend
//
// In dev we rely on Vite's proxy for API calls, so the API base is '' (same-origin).
// Redirect links (/{shortCode}) should still point to the backend server.
const API_BASE_URL = import.meta.env.DEV ? '' : 'http://localhost:8080';
const SHORTLINK_BASE_URL =
    import.meta.env.VITE_SHORTLINK_BASE_URL ??
    (import.meta.env.DEV ? 'http://localhost:8080' : API_BASE_URL);

export interface ShortenRequest {
    url: string;
    custom_code?: string;
    expires_at?: string;
}

export interface ShortenResponse {
    short_url: string;
    original_url: string;
    code: string;
    expires_at?: string;
    qr_code?: string;
}

export interface URLDetails {
    id: number;
    short_code: string;
    original_url: string;
    created_at: string;
    expires_at?: string;
    click_count: number;
    last_clicked?: string;
}

export interface URLHistoryItem extends URLDetails {
    status: string;
}

class ApiService {
    private baseUrl: string;
    private shortlinkBaseUrl: string;

    constructor() {
        this.baseUrl = API_BASE_URL;
        this.shortlinkBaseUrl = SHORTLINK_BASE_URL;
    }

    // Create a short URL
    async createShortURL(request: ShortenRequest): Promise<ShortenResponse> {
        try {
            console.log('Making request to:', `${this.baseUrl}/api/shorten`);
            console.log('Request payload:', request);

            const response = await fetch(`${this.baseUrl}/api/shorten`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(request),
            });

            console.log('Response status:', response.status);
            console.log('Response headers:', response.headers);

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Error response:', errorText);
                throw new Error(`Failed to create short URL: ${errorText}`);
            }

            const result = await response.json();
            console.log('Success response:', result);
            return result;
        } catch (error) {
            console.error('Fetch error:', error);
            throw error;
        }
    }

    // Get URL details
    async getURLDetails(shortCode: string): Promise<URLDetails> {
        const response = await fetch(`${this.baseUrl}/api/urls/${shortCode}`);

        if (!response.ok) {
            throw new Error(`Failed to get URL details: ${response.statusText}`);
        }

        return response.json();
    }

    // Get all URLs
    async getAllURLs(): Promise<URLDetails[]> {
        const response = await fetch(`${this.baseUrl}/api/urls`);

        if (!response.ok) {
            throw new Error(`Failed to get URLs: ${response.statusText}`);
        }

        return response.json();
    }

    // Delete a URL
    async deleteURL(shortCode: string): Promise<void> {
        const response = await fetch(`${this.baseUrl}/api/urls/${shortCode}`, {
            method: 'DELETE',
        });

        if (!response.ok) {
            throw new Error(`Failed to delete URL: ${response.statusText}`);
        }
    }

    // Get URL history
    async getURLHistory(): Promise<URLHistoryItem[]> {
        const response = await fetch(`${this.baseUrl}/api/history`);

        if (!response.ok) {
            throw new Error(`Failed to get URL history: ${response.statusText}`);
        }

        return response.json();
    }

    // Build the full short URL for display
    getShortURL(shortCode: string): string {
        return `${this.shortlinkBaseUrl}/${shortCode}`;
    }
}

export const apiService = new ApiService();
