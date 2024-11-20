import requests
import websockets
import asyncio
import json

class NixDeployClient:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url.rstrip('/')

    def init_configuration(self):
        """Initialize configuration.nix"""
        response = requests.post(f"{self.base_url}/initConfiguration")
        response.raise_for_status()
        return response.ok

    def list_nix_files(self):
        """List all nix files"""
        response = requests.get(f"{self.base_url}/listNixFiles")
        response.raise_for_status()
        return response.json()

    def get_nix_files_contents(self):
        """Get contents of all nix files"""
        response = requests.get(f"{self.base_url}/getNixFilesContents")
        response.raise_for_status()
        return response.text

    def set_nix_files_contents(self, files):
        """Set contents of nix files
        
        Args:
            files (dict): Dictionary mapping filenames to their contents
        """
        response = requests.post(f"{self.base_url}/setNixFilesContents", json=files)
        response.raise_for_status()
        return response.ok

    def run_dry_build(self):
        """Run a dry build"""
        response = requests.post(f"{self.base_url}/runDryBuild")
        response.raise_for_status()
        return response.text

    def run_test(self):
        """Run a test build"""
        response = requests.post(f"{self.base_url}/runTest")
        response.raise_for_status()
        return response.text

    def run_switch(self):
        """Run switch to apply configuration"""
        response = requests.post(f"{self.base_url}/runSwitch")
        response.raise_for_status()
        return response.text

    def liveness_check(self):
        """Check server liveness"""
        response = requests.get(f"{self.base_url}/livenessCheck")
        response.raise_for_status()
        return response.text

    def fetch_local_url(self, url):
        """Fetch a local URL"""
        response = requests.get(f"{self.base_url}/fetchLocalUrl", params={"url": url})
        response.raise_for_status()
        return response.text

    async def connect_ssh(self):
        """Connect to SSH via WebSocket"""
        uri = f"ws://{self.base_url.replace('http://', '')}/ssh"
        async with websockets.connect(uri) as websocket:
            return websocket 