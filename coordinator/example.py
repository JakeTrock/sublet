from client import NixDeployClient
import asyncio

async def main():
    # Create client instance
    client = NixDeployClient()

    # Initialize configuration
    print("Initializing configuration...")
    client.init_configuration()

    # List nix files
    print("\nListing Nix files:")
    files = client.list_nix_files()
    print(files)

    # Get file contents
    print("\nGetting Nix file contents:")
    contents = client.get_nix_files_contents()
    print(contents)

    # Run a dry build
    print("\nRunning dry build:")
    build_output = client.run_dry_build()
    print(build_output)

    # Check liveness
    print("\nChecking liveness:")
    liveness = client.liveness_check()
    print(liveness)

    # Example SSH connection (just demonstration)
    print("\nConnecting to SSH...")
    ws = await client.connect_ssh()
    # Here you would handle the websocket connection for SSH
    # This is just a basic example

if __name__ == "__main__":
    asyncio.run(main()) 