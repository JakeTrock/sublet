// Base URL for API requests
const API_BASE_URL = 'http://localhost:3000';

// Type definitions for API responses
interface FileListResponse {
    files: string[];
}

interface FileContentResponse {
    content: string;
}

interface BuildResponse {
    success: boolean;
    output: string;
}

interface HealthCheckResponse {
    status: string;
}

// API functions

/**
 * Initialize configuration
 */
export async function initConfiguration(config: object): Promise<Response> {
    return fetch(`${API_BASE_URL}/init`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(config),
    });
}

/**
 * Get list of Nix files
 */
export async function listNixFiles(): Promise<FileListResponse> {
    const response = await fetch(`${API_BASE_URL}/files`);
    return response.json();
}

/**
 * Get contents of a Nix file
 */
export async function getNixContents(path: string): Promise<FileContentResponse> {
    const response = await fetch(`${API_BASE_URL}/contents?path=${encodeURIComponent(path)}`);
    return response.json();
}

/**
 * Set contents of a Nix file
 */
export async function setNixContents(path: string, content: string): Promise<Response> {
    return fetch(`${API_BASE_URL}/contents`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ path, content }),
    });
}

/**
 * Run a dry build
 */
export async function runDryBuild(): Promise<BuildResponse> {
    const response = await fetch(`${API_BASE_URL}/dry-build`, {
        method: 'POST',
    });
    return response.json();
}

/**
 * Run tests
 */
export async function runTest(): Promise<BuildResponse> {
    const response = await fetch(`${API_BASE_URL}/test`, {
        method: 'POST',
    });
    return response.json();
}

/**
 * Run system switch
 */
export async function runSwitch(): Promise<BuildResponse> {
    const response = await fetch(`${API_BASE_URL}/switch`, {
        method: 'POST',
    });
    return response.json();
}

/**
 * Check system health
 */
export async function checkHealth(): Promise<HealthCheckResponse> {
    const response = await fetch(`${API_BASE_URL}/health`);
    return response.json();
}

/**
 * Check text liveness
 */
export async function textLiveness(): Promise<string> {
    const response = (await fetch(`${API_BASE_URL}/text`)).text();
    return response;
}

/**
 * Fetch local URL
 */
export async function fetchLocalUrl(url: string): Promise<Response> {
    const response = await fetch(`${API_BASE_URL}/fetch?url=${encodeURIComponent(url)}`);
    return response;
}

/**
 * Create WebSocket connection for SSH
 */
export function connectSSH(): WebSocket {
    return new WebSocket(`ws://localhost:3000/ws/ssh`);
}
