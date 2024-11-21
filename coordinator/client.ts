import { NixServiceClient } from './generated/api';

export class NixDeployClient {
    private client: NixServiceClient;

    constructor(baseURL: string = 'http://localhost:8080') {
        this.client = new NixServiceClient(baseURL);
    }

    async initConfiguration(): Promise<boolean> {
        const response = await this.client.init_config({});
        if (response.error) {
            throw new Error(response.error);
        }
        return response.success;
    }

    async listNixFiles(): Promise<string[]> {
        const response = await this.client.list_files({});
        if (response.error) {
            throw new Error(response.error);
        }
        return response.files;
    }

    async getNixFilesContents(): Promise<string> {
        const response = await this.client.get_contents({});
        if (response.error) {
            throw new Error(response.error);
        }
        return response.contents;
    }

    async setNixFilesContents(files: Record<string, string>): Promise<boolean> {
        const response = await this.client.set_contents({ files });
        if (response.error) {
            throw new Error(response.error);
        }
        return response.success;
    }

    async runDryBuild(): Promise<string> {
        const response = await this.client.dry_build({});
        if (response.error) {
            throw new Error(response.error);
        }
        return response.output;
    }

    async runTest(): Promise<string> {
        const response = await this.client.test({});
        if (response.error) {
            throw new Error(response.error);
        }
        return response.output;
    }

    async runSwitch(): Promise<string> {
        const response = await this.client.switch({});
        if (response.error) {
            throw new Error(response.error);
        }
        return response.output;
    }

    async livenessCheck(): Promise<string> {
        const response = await this.client.liveness({});
        if (response.error) {
            throw new Error(response.error);
        }
        return response.output;
    }
} 