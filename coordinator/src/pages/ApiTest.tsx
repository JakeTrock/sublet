import { useState } from 'react';
import * as api from '../connection';

interface PanelProps {
    title: string;
    children: React.ReactNode;
}

function Panel({ title, children }: PanelProps) {
    return (
        <div className="border rounded-lg p-4 mb-4">
            <h2 className="text-xl font-bold mb-3">{title}</h2>
            {children}
        </div>
    );
}

export default function ApiTest() {
    const [results, setResults] = useState<Record<string, any>>({});

    // Helper to update results for a specific endpoint
    const updateResult = (endpoint: string, data: any) => {
        setResults(prev => ({
            ...prev,
            [endpoint]: JSON.stringify(data, null, 2)
        }));
    };

    // Test handlers
    const handleInit = async () => {
        try {
            const response = await api.initConfiguration({ test: true });
            updateResult('init', await response.json());
        } catch (error) {
            updateResult('init', { error: error.message });
        }
    };

    const handleListFiles = async () => {
        try {
            const response = await api.listNixFiles();
            updateResult('files', response);
        } catch (error) {
            updateResult('files', { error: error.message });
        }
    };

    const handleGetContents = async () => {
        try {
            const response = await api.getNixContents('/etc/nixos/configuration.nix');
            updateResult('contents', response);
        } catch (error) {
            updateResult('contents', { error: error.message });
        }
    };

    const handleSetContents = async () => {
        try {
            const response = await api.setNixContents(
                '/etc/nixos/test.nix',
                '{ config, pkgs, ... }: { }'
            );
            updateResult('setContents', await response.json());
        } catch (error) {
            updateResult('setContents', { error: error.message });
        }
    };

    const handleDryBuild = async () => {
        try {
            const response = await api.runDryBuild();
            updateResult('dryBuild', response);
        } catch (error) {
            updateResult('dryBuild', { error: error.message });
        }
    };

    const handleTest = async () => {
        try {
            const response = await api.runTest();
            updateResult('test', response);
        } catch (error) {
            updateResult('test', { error: error.message });
        }
    };

    const handleSwitch = async () => {
        try {
            const response = await api.runSwitch();
            updateResult('switch', response);
        } catch (error) {
            updateResult('switch', { error: error.message });
        }
    };

    const handleHealth = async () => {
        try {
            const response = await api.checkHealth();
            updateResult('health', response);
        } catch (error) {
            updateResult('health', { error: error.message });
        }
    };

    const handleTextLiveness = async () => {
        try {
            const response = await api.textLiveness();
            console.log(response);
            updateResult('textLiveness', response);
        } catch (error) {
            updateResult('textLiveness', { error: error.message });
        }
    };

    const handleFetch = async () => {
        try {
            const response = await api.fetchLocalUrl('http://localhost:3000/health');
            updateResult('fetch', await response.json());
        } catch (error) {
            updateResult('fetch', { error: error.message });
        }
    };

    const handleSSH = () => {
        try {
            const ws = api.connectSSH();
            ws.onopen = () => updateResult('ssh', { status: 'Connected' });
            ws.onclose = () => updateResult('ssh', { status: 'Disconnected' });
            ws.onerror = (error) => updateResult('ssh', { error: 'WebSocket error' });
        } catch (error) {
            updateResult('ssh', { error: error.message });
        }
    };

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-2xl font-bold mb-6">API Test Dashboard</h1>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Panel title="Initialize Configuration">
                    <button onClick={handleInit} className="btn">Test Init</button>
                    <pre className="mt-2">{results.init || 'No result'}</pre>
                </Panel>

                <Panel title="List Nix Files">
                    <button onClick={handleListFiles} className="btn">List Files</button>
                    <pre className="mt-2">{results.files || 'No result'}</pre>
                </Panel>

                <Panel title="Get Nix Contents">
                    <button onClick={handleGetContents} className="btn">Get Contents</button>
                    <pre className="mt-2">{results.contents || 'No result'}</pre>
                </Panel>

                <Panel title="Set Nix Contents">
                    <button onClick={handleSetContents} className="btn">Set Contents</button>
                    <pre className="mt-2">{results.setContents || 'No result'}</pre>
                </Panel>

                <Panel title="Dry Build">
                    <button onClick={handleDryBuild} className="btn">Run Dry Build</button>
                    <pre className="mt-2">{results.dryBuild || 'No result'}</pre>
                </Panel>

                <Panel title="Run Test">
                    <button onClick={handleTest} className="btn">Run Test</button>
                    <pre className="mt-2">{results.test || 'No result'}</pre>
                </Panel>

                <Panel title="System Switch">
                    <button onClick={handleSwitch} className="btn">Run Switch</button>
                    <pre className="mt-2">{results.switch || 'No result'}</pre>
                </Panel>

                <Panel title="Health Check">
                    <button onClick={handleHealth} className="btn">Check Health</button>
                    <pre className="mt-2">{results.health || 'No result'}</pre>
                </Panel>

                <Panel title="Text Liveness">
                    <button onClick={handleTextLiveness} className="btn">Check Text Liveness</button>
                    <pre className="mt-2">{results.textLiveness || 'No result'}</pre>
                </Panel>

                <Panel title="Fetch Local URL">
                    <button onClick={handleFetch} className="btn">Fetch URL</button>
                    <pre className="mt-2">{results.fetch || 'No result'}</pre>
                </Panel>

                <Panel title="SSH Connection">
                    <button onClick={handleSSH} className="btn">Connect SSH</button>
                    <pre className="mt-2">{results.ssh || 'No result'}</pre>
                </Panel>
            </div>
        </div>
    );
} 