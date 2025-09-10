import { Server, Response } from 'miragejs';
import { AppSchema } from '../types';

export function routesForNetwork(server: Server) {
  // GET /network - Fetch network configuration
  server.get('/network', (schema: AppSchema) => {
    const networkData = schema.first('network');
    if (!networkData) {
      return new Response(404, {}, { code: '404', error: 'No network configuration found' });
    }
    return networkData.attrs;
  });

  // POST /network - Create a new network
  server.post('/network', (schema: AppSchema, request) => {
    const requestBody = JSON.parse(request.requestBody);
    const { name, url, default: makeDefault } = requestBody;

    // Get current network data
    const networkData = schema.first('network');
    if (!networkData) {
      return new Response(400, {}, { code: '400', error: 'No network configuration found' });
    }

    const currentNetworks = [...networkData.availableNetworkInfos];
    
    // Check if network name already exists
    const existingNetwork = currentNetworks.find(net => net.name.toLowerCase() === name.toLowerCase());
    if (existingNetwork) {
      return new Response(400, {}, { code: '400', error: 'Network name already exists' });
    }

    // Create new network info
    const newNetwork = {
      name,
      url,
      chainId: Math.floor(Math.random() * 100000000), // Generate random chain ID for mock
      version: '1.0.0',
      status: 'up' as const
    };

    // Add new network to the list
    currentNetworks.push(newNetwork);

    // Update the network data
    const updatedData = {
      availableNetworkInfos: currentNetworks,
      currentNetwork: makeDefault ? name : networkData.currentNetwork
    };

    networkData.update(updatedData);

    return new Response(201, {}, updatedData);
  });

  // PUT /network/:network - Update an existing network
  server.put('/network/:network', (schema: AppSchema, request) => {
    const { network: networkName } = request.params;
    const requestBody = JSON.parse(request.requestBody);
    const { url, newName, default: makeDefault } = requestBody;

    // Get current network data
    const networkData = schema.first('network');
    if (!networkData) {
      return new Response(400, {}, { code: '400', error: 'No network configuration found' });
    }

    const currentNetworks = [...networkData.availableNetworkInfos];
    
    // Find the network to update
    const networkIndex = currentNetworks.findIndex(net => net.name === networkName);
    if (networkIndex === -1) {
      return new Response(404, {}, { code: '404', error: 'Network not found' });
    }

    // Check if new name conflicts with existing networks (excluding current network)
    if (newName && newName !== networkName) {
      const existingNetwork = currentNetworks.find((net, index) => 
        index !== networkIndex && net.name.toLowerCase() === newName.toLowerCase()
      );
      if (existingNetwork) {
        return new Response(400, {}, { code: '400', error: 'Network name already exists' });
      }
    }

    // Update the network
    const updatedNetwork = { ...currentNetworks[networkIndex] };
    if (url) updatedNetwork.url = url;
    if (newName) updatedNetwork.name = newName;
    
    currentNetworks[networkIndex] = updatedNetwork;

    // Update current network if this network was renamed or made default
    let currentNetwork = networkData.currentNetwork;
    if (makeDefault) {
      currentNetwork = newName || networkName;
    } else if (newName && networkData.currentNetwork === networkName) {
      currentNetwork = newName;
    }

    const updatedData = {
      availableNetworkInfos: currentNetworks,
      currentNetwork
    };

    networkData.update(updatedData);

    return new Response(200, {}, updatedData);
  });

  // DELETE /network/:network - Delete a network
  server.delete('/network/:network', (schema: AppSchema, request) => {
    const { network: networkName } = request.params;

    // Get current network data
    const networkData = schema.first('network');
    if (!networkData) {
      return new Response(400, {}, { code: '400', error: 'No network configuration found' });
    }

    const currentNetworks = [...networkData.availableNetworkInfos];
    
    // Find the network to delete
    const networkIndex = currentNetworks.findIndex(net => net.name === networkName);
    if (networkIndex === -1) {
      return new Response(404, {}, { code: '404', error: 'Network not found' });
    }

    // Prevent deletion if it's the only network
    if (currentNetworks.length === 1) {
      return new Response(400, {}, { code: '400', error: 'Cannot delete the last remaining network' });
    }

    // Remove the network
    currentNetworks.splice(networkIndex, 1);

    // If the deleted network was the current network, switch to the first available
    let currentNetwork = networkData.currentNetwork;
    if (currentNetwork === networkName) {
      currentNetwork = currentNetworks[0].name;
    }

    const updatedData = {
      availableNetworkInfos: currentNetworks,
      currentNetwork
    };

    networkData.update(updatedData);

    return new Response(200, {}, updatedData);
  });

  // POST /network/:network - Switch to a network (existing functionality)
  server.post('/network/:network', (schema: AppSchema, request) => {
    const { network } = request.params;
    const storedNetwork = schema.first('network');

    if (!storedNetwork)
      return new Response(404, {}, { code: '404', error: 'Not Found' });

    // Check if the network exists in available networks
    const networkExists = storedNetwork.availableNetworkInfos.some(net => net.name === network);
    if (!networkExists) {
      return new Response(404, {}, { code: '404', error: 'Network not found in available networks' });
    }

    storedNetwork.update({ currentNetwork: network });

    return new Response(200, {}, storedNetwork.attrs);
  });
}
