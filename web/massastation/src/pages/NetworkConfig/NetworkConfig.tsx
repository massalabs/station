import { useMemo, useState, useCallback, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Button, Input } from '@massalabs/react-ui-kit';
import {
  FiCheckCircle,
  FiAlertTriangle,
  FiGlobe,
  FiLink,
  FiHash,
  FiPlus,
  FiEdit3,
  FiTrash2,
  FiSettings
} from 'react-icons/fi';
import { URL } from '@/const/url/url';
import { useResource } from '@/custom/api/useResource';
import { usePost } from '@/custom/api/usePost';
import { usePut } from '@/custom/api/usePut';
import { useDelete } from '@/custom/api/useDelete';
import type { NetworkModel } from '@/models/NetworkModel';

type CreateNetworkBody = { name: string; url: string; default?: boolean };
type UpdateNetworkBody = { url?: string; default?: boolean; newName?: string };

export function NetworkConfig() {
  const { data, refetch } = useResource<NetworkModel>(URL.PATH_NETWORKS);
  const [searchParams, setSearchParams] = useSearchParams();

  const [selectedNetwork, setSelectedNetwork] = useState<string>('');
  const [showCreate, setShowCreate] = useState(false);
  const [showEdit, setShowEdit] = useState(false);
  const [showDelete, setShowDelete] = useState(false);

  const [newName, setNewName] = useState('');
  const [newURL, setNewURL] = useState('');
  const [makeDefault, setMakeDefault] = useState(false);
  const [nameError, setNameError] = useState('');

  const createNetwork = usePost<CreateNetworkBody, NetworkModel>(
    URL.PATH_NETWORKS,
  );
  const updateNetwork = usePut<UpdateNetworkBody, NetworkModel>(
    `${URL.PATH_NETWORKS}/${selectedNetwork}`,
  );
  const deleteNetwork = useDelete<NetworkModel>(
    `${URL.PATH_NETWORKS}/${selectedNetwork}`,
  );

  const networks = useMemo(() => data?.availableNetworkInfos || [], [data?.availableNetworkInfos]);
  const currentInfo = useMemo(
    () => networks.find((n) => n.name === data?.currentNetwork),
    [networks, data?.currentNetwork],
  );

  // Validate network name uniqueness
  const validateNetworkName = useCallback((name: string) => {
    if (!name.trim()) {
      return 'Network name is required';
    }
    
    const normalizedName = name.trim().toLowerCase();
    const existingNetwork = networks.find(network => 
      network.name.toLowerCase() === normalizedName
    );
    
    if (existingNetwork) {
      return `Network name "${name}" already exists`;
    }
    
    return '';
  }, [networks]);

  const handleNameChange = useCallback((e: React.ChangeEvent<HTMLInputElement> | string) => {
    const value = typeof e === 'string' ? e : e.target.value;
    setNewName(value);
    
    // Validate name in real-time
    const error = validateNetworkName(value);
    setNameError(error);
  }, [validateNetworkName]);

  const handleURLChange = useCallback((e: React.ChangeEvent<HTMLInputElement> | string) => {
    const value = typeof e === 'string' ? e : e.target.value;
    setNewURL(value);
  }, []);

  const resetForms = useCallback(() => {
    setNewName('');
    setNewURL('');
    setMakeDefault(false);
    setNameError('');
  }, []);

  const onCreate = useCallback(async () => {
    // Validate before creating
    const nameValidationError = validateNetworkName(newName);
    if (nameValidationError) {
      setNameError(nameValidationError);
      return;
    }

    await createNetwork.mutateAsync({ payload: { name: newName.trim(), url: newURL, default: makeDefault } });
    resetForms();
    setShowCreate(false);
    await refetch();
  }, [createNetwork, newName, newURL, makeDefault, resetForms, refetch, validateNetworkName]);

  const onUpdate = useCallback(async () => {
    await updateNetwork.mutateAsync({ url: newURL || undefined, newName: newName || undefined, default: makeDefault });
    resetForms();
    setShowEdit(false);
    await refetch();
  }, [updateNetwork, newURL, newName, makeDefault, resetForms, refetch]);

  const onDelete = useCallback(async () => {
    await deleteNetwork.mutateAsync(undefined as unknown as NetworkModel);
    setShowDelete(false);
    await refetch();
  }, [deleteNetwork, refetch]);

  const handleShowCreate = useCallback(() => {
    resetForms();
    setShowCreate(true);
  }, [resetForms]);

  const handleShowEdit = useCallback(() => {
    resetForms();
    setShowEdit(true);
  }, [resetForms]);

  const handleShowDelete = useCallback(() => {
    setShowDelete(true);
  }, []);

  const handleCloseCreate = useCallback(() => {
    setShowCreate(false);
  }, []);

  const handleCloseEdit = useCallback(() => {
    setShowEdit(false);
  }, []);

  const handleCloseDelete = useCallback(() => {
    setShowDelete(false);
  }, []);

  const handleMakeDefaultChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setMakeDefault(e.target.checked);
  }, []);

  // Handle URL query parameters for prefilling create network form
  useEffect(() => {
    const encodedName = searchParams.get('name');
    const encodedUrl = searchParams.get('url');
    const defaultParam = searchParams.get('default');

    // Decode the name parameter if it exists
    let name = '';
    if (encodedName) {
      try {
        name = decodeURIComponent(encodedName);
      } catch (error) {
        console.warn('Failed to decode name parameter:', encodedName, error);
        name = encodedName; // Fallback to original if decoding fails
      }
    }

    // Decode the URL parameter if it exists
    let url = '';
    if (encodedUrl) {
      try {
        url = decodeURIComponent(encodedUrl);
      } catch (error) {
        console.warn('Failed to decode URL parameter:', encodedUrl, error);
        url = encodedUrl; // Fallback to original if decoding fails
      }
    }

    // If we have query parameters for creating a network
    if (name && url) {
      // Prefill the form fields
      setNewName(name);
      setNewURL(url);
      if (defaultParam === 'true') setMakeDefault(true);
      
      // Open the create modal
      setShowCreate(true);
      
      // Clear the query parameters to prevent re-triggering
      setSearchParams({});
    }
  }, [searchParams, setSearchParams]);

  const SimpleModal = useCallback(({ show, onClose, title, children, footer }:
     { show: boolean; onClose: () => void; title: string; children: React.ReactNode; footer: React.ReactNode }) => {
    if (!show) return null;
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
        <div className="relative bg-secondary border border-tertiary/20 p-8 rounded-2xl shadow-2xl
                     min-w-[480px] max-w-[600px] w-full">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-2 bg-c-primary/10 rounded-lg">
              <FiSettings className="w-6 h-6 text-c-primary" />
            </div>
            <h3 className="text-xl font-bold text-f-primary">{title}</h3>
          </div>
          <div className="mb-8">{children}</div>
          <div className="flex gap-3 justify-end">{footer}</div>
        </div>
      </div>
    );
  }, []);

  return (
    <div className="bg-primary text-f-primary pt-24 px-8 pb-12 min-h-screen">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Header Section */}
        <div className="flex items-center gap-6 mb-8">
          <div className="p-4 bg-secondary rounded-xl">
            <FiSettings className="w-10 h-10 text-c-primary" />
          </div>
          <div className="flex flex-col justify-center">
            <h1 className="text-4xl font-bold text-f-primary">
              Massa Station Network Configuration
            </h1>
            <p className="text-neutral mt-2">Manage your blockchain network connections</p>
          </div>
        </div>

        {/* Active Network Section */}
        <div className="bg-secondary/50 backdrop-blur-sm rounded-2xl border border-tertiary/20 p-8">
          <div className="flex items-center gap-3 mb-6">
            <FiGlobe className="w-6 h-6 text-c-primary" />
            <h2 className="text-xl font-semibold">Active Network</h2>
          </div>

          <div className="bg-secondary rounded-xl p-6 mb-6 border border-tertiary/10">
            <div className="flex items-center gap-4 mb-4">
              <div className="p-2 bg-c-primary/10 rounded-lg">
                <FiGlobe className="w-6 h-6 text-c-primary" />
              </div>
              <div>
                <h3 className="text-2xl font-bold text-f-primary">{data?.currentNetwork || 'No Network'}</h3>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="bg-secondary rounded-xl p-6 border border-tertiary/10
                          hover:border-c-primary/30 transition-colors">
              <div className="flex items-center gap-3 mb-3">
                <div className="p-2 bg-green-500/10 rounded-lg">
                  {currentInfo?.status === 'up' ? (
                    <FiCheckCircle className="w-5 h-5 text-green-500" />
                  ) : (
                    <FiAlertTriangle className="w-5 h-5 text-red-500" />
                  )}
                </div>
                <span className="text-sm font-medium text-neutral">Network Status</span>
              </div>
              <div className="flex items-center gap-2">
                <span className={`text-lg font-semibold ${
                  currentInfo?.status === 'up' ? 'text-green-500' : 'text-red-500'
                }`}>
                  {currentInfo?.status === 'up'
                    ? 'Online'
                    : currentInfo?.status === 'down'
                    ? 'Offline'
                    : 'Unknown'}
                </span>
              </div>
            </div>

            <div className="bg-secondary rounded-xl p-6 border border-tertiary/10
                          hover:border-c-primary/30 transition-colors">
              <div className="flex items-center gap-3 mb-3">
                <div className="p-2 bg-blue-500/10 rounded-lg">
                  <FiLink className="w-5 h-5 text-blue-500" />
                </div>
                <span className="text-sm font-medium text-neutral">RPC Endpoint</span>
              </div>
              <div className="text-sm text-f-primary break-all font-mono bg-tertiary/30 p-2 rounded-lg">
                {currentInfo?.url || 'Not configured'}
              </div>
            </div>

            <div className="bg-secondary rounded-xl p-6 border border-tertiary/10
                          hover:border-c-primary/30 transition-colors">
              <div className="flex items-center gap-3 mb-3">
                <div className="p-2 bg-purple-500/10 rounded-lg">
                  <FiGlobe className="w-5 h-5 text-purple-500" />
                </div>
                <span className="text-sm font-medium text-neutral">Network Name</span>
              </div>
              <div className="text-lg font-semibold text-f-primary">
                {data?.currentNetwork || '-'}
              </div>
            </div>

            <div className="bg-secondary rounded-xl p-6 border border-tertiary/10
                          hover:border-c-primary/30 transition-colors">
              <div className="flex items-center gap-3 mb-3">
                <div className="p-2 bg-orange-500/10 rounded-lg">
                  <FiHash className="w-5 h-5 text-orange-500" />
                </div>
                <span className="text-sm font-medium text-neutral">Chain ID</span>
              </div>
              <div className="text-lg font-semibold text-f-primary font-mono">
                {currentInfo?.chainId || '-'}
              </div>
            </div>
          </div>
        </div>

        {/* Available Networks Section */}
        <div className="bg-secondary/50 backdrop-blur-sm rounded-2xl border border-tertiary/20 p-8">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-3">
              <FiGlobe className="w-6 h-6 text-c-primary" />
              <h2 className="text-xl font-semibold">Available Networks</h2>
            </div>
            <div className="text-sm text-neutral">
              {networks.length} network{networks.length !== 1 ? 's' : ''} configured
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
            {networks.map((net) => {
              const isActive = selectedNetwork === net.name;
              const isCurrentNetwork = data?.currentNetwork === net.name;
              return (
                <div
                  key={net.name}
                  className={`
                    relative bg-secondary rounded-xl p-6 cursor-pointer border-2 transition-all
                    duration-200 hover:scale-[1.02] hover:shadow-lg
                    ${isActive
                    ? 'border-c-primary shadow-lg shadow-c-primary/20'
                    : 'border-tertiary/10 hover:border-c-primary/40'
                }
                    ${isCurrentNetwork ? 'ring-2 ring-c-primary/30' : ''}
                  `}
                  onClick={() => {
                    setSelectedNetwork(net.name);
                    setShowEdit(false);
                    setShowDelete(false);
                  }}
                >

                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-c-primary/10 rounded-lg">
                        <FiGlobe className="w-5 h-5 text-c-primary" />
                      </div>
                      <h3 className="text-lg font-bold text-f-primary">{net.name}</h3>
                    </div>
                    <div className={`
                      flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium
                      ${net.status === 'up'
                      ? 'bg-green-500/10 text-green-500'
                      : 'bg-red-500/10 text-red-500'
                }
                    `}>
                      {net.status === 'up' ? (
                        <FiCheckCircle className="w-3 h-3" />
                      ) : (
                        <FiAlertTriangle className="w-3 h-3" />
                      )}
                      <span>{net.status === 'up' ? 'Online' : 'Offline'}</span>
                    </div>
                  </div>

                  <div className="space-y-3">
                    <div>
                      <div className="text-xs font-medium text-neutral mb-1">RPC Endpoint</div>
                      <div className="text-sm text-f-primary break-all font-mono bg-tertiary/20 p-2 rounded-lg">
                        {net.url}
                      </div>
                    </div>

                    <div className="flex items-center justify-between">
                      <div>
                        <div className="text-xs font-medium text-neutral">Chain ID</div>
                        <div className="text-sm font-bold text-f-primary font-mono">{net.chainId}</div>
                      </div>
                      <div>
                        <div className="text-xs font-medium text-neutral">Version</div>
                        <div className="text-sm font-bold text-f-primary">{net.version}</div>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Action Buttons */}
        <div className="bg-secondary/50 backdrop-blur-sm rounded-2xl border border-tertiary/20 p-8">
          <div className="flex items-center gap-3 mb-6">
            <FiSettings className="w-6 h-6 text-c-primary" />
            <h2 className="text-xl font-semibold">Network Actions</h2>
          </div>

          <div className="flex flex-wrap gap-4">
            <Button
              onClick={handleShowCreate}
              className="flex items-center gap-2 px-6 py-3 bg-c-primary hover:bg-c-primary/80
                         font-medium rounded-xl transition-all"
            >
              <FiPlus className="w-5 h-5" />
              Add New Network
            </Button>

            <Button
              disabled={!selectedNetwork}
              onClick={handleShowEdit}
              variant="secondary"
              className="flex items-center gap-2 px-6 py-3 font-medium rounded-xl transition-all disabled:opacity-50"
            >
              <FiEdit3 className="w-5 h-5" />
              Edit Network
            </Button>

            <Button
              disabled={!selectedNetwork}
              variant="danger"
              onClick={handleShowDelete}
              className="flex items-center gap-2 px-6 py-3 font-medium rounded-xl transition-all disabled:opacity-50"
            >
              <FiTrash2 className="w-5 h-5" />
              Remove Network
            </Button>
          </div>

          {selectedNetwork && (
            <div className="mt-4 p-4 bg-c-primary/10 rounded-xl border border-c-primary/20">
              <p className="text-sm text-neutral">
                <span className="font-medium text-c-primary">{selectedNetwork}</span> is selected.{' '}
                Choose an action above to modify this network configuration.
              </p>
            </div>
          )}
        </div>

        {/* Modals */}
        <SimpleModal
          show={showCreate}
          onClose={handleCloseCreate}
          title="Add New Network"
          children={
            <div className="space-y-6">
              <p className="text-neutral text-sm leading-relaxed">
              Configure a new blockchain network connection. This will allow you to interact with a{' '}
              different network endpoint.
              </p>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-f-primary mb-2">
                  Network Name
                  </label>
                  <Input
                    name="create-network-name"
                    value={newName}
                    onChange={handleNameChange}
                    customClass={`bg-tertiary/20 border rounded-lg p-3 w-full ${
                      nameError ? 'border-red-500' : 'border-tertiary/30'
                    }`}
                    placeholder="e.g., mainnet, testnet, custom-network"
                  />
                  {nameError ? (
                    <p className="text-xs text-red-500 mt-1">{nameError}</p>
                  ) : (
                    <p className="text-xs text-neutral mt-1">Choose a descriptive name for this network</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-f-primary mb-2">
                  RPC Endpoint URL
                  </label>
                  <Input
                    name="create-network-url"
                    value={newURL}
                    onChange={handleURLChange}
                    customClass="bg-tertiary/20 border border-tertiary/30 rounded-lg p-3 w-full font-mono text-sm"
                    placeholder="https://api.example.com/rpc"
                  />
                  <p className="text-xs text-neutral mt-1">The HTTP/HTTPS endpoint for blockchain RPC calls</p>
                </div>

                <div className="p-4 bg-tertiary/10 rounded-xl border border-tertiary/20">
                  <label className="flex items-center gap-3 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={makeDefault}
                      onChange={handleMakeDefaultChange}
                      className="w-4 h-4 text-c-primary bg-tertiary border-tertiary rounded focus:ring-c-primary"
                    />
                    <div>
                      <span className="text-sm font-medium text-f-primary">Set as default network</span>
                      <p className="text-xs text-neutral">Switch to this network immediately after creation</p>
                    </div>
                  </label>
                </div>
              </div>
            </div>
          }
          footer={
            <>
              <Button onClick={handleCloseCreate} variant="secondary" className="px-6 py-2">
              Cancel
              </Button>
              <Button
                onClick={onCreate}
                disabled={!newName || !newURL || !!nameError}
                className="px-6 py-2 bg-c-primary hover:bg-c-primary/80 font-medium disabled:opacity-50"
              >
                <FiPlus className="w-4 h-4 mr-2" />
              Create Network
              </Button>
            </>
          }
        />

        <SimpleModal
          show={showEdit}
          onClose={handleCloseEdit}
          title={`Modify Network Configuration`}
          children={
            <div className="space-y-6">
              <div className="p-4 bg-c-primary/10 rounded-xl border border-c-primary/20">
                <p className="text-sm text-f-primary">
                  <span className="font-bold text-c-primary">{selectedNetwork}</span> network configuration
                </p>
                <p className="text-xs text-neutral mt-1">
                Modify the settings below. Leave fields empty to keep current values.
                </p>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-f-primary mb-2">
                  Network Name
                  </label>
                  <Input
                    name="edit-network-name"
                    value={newName}
                    onChange={handleNameChange}
                    customClass="bg-tertiary/20 border border-tertiary/30 rounded-lg p-3 w-full"
                    placeholder={`Current: ${selectedNetwork}`}
                  />
                  <p className="text-xs text-neutral mt-1">Leave empty to keep current name</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-f-primary mb-2">
                  RPC Endpoint URL
                  </label>
                  <Input
                    name="edit-network-url"
                    value={newURL}
                    onChange={handleURLChange}
                    customClass="bg-tertiary/20 border border-tertiary/30 rounded-lg p-3 w-full font-mono text-sm"
                    placeholder="Enter new RPC URL or leave empty"
                  />
                  <p className="text-xs text-neutral mt-1">Leave empty to keep current endpoint</p>
                </div>

                <div className="p-4 bg-tertiary/10 rounded-xl border border-tertiary/20">
                  <label className="flex items-center gap-3 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={makeDefault}
                      onChange={handleMakeDefaultChange}
                      className="w-4 h-4 text-c-primary bg-tertiary border-tertiary rounded focus:ring-c-primary"
                    />
                    <div>
                      <span className="text-sm font-medium text-f-primary">Set as default network</span>
                      <p className="text-xs text-neutral">Make this the active network after saving</p>
                    </div>
                  </label>
                </div>
              </div>
            </div>
          }
          footer={
            <>
              <Button onClick={handleCloseEdit} variant="secondary" className="px-6 py-2">
              Cancel
              </Button>
              <Button
                onClick={onUpdate}
                disabled={!newName && !newURL && !makeDefault}
                className="px-6 py-2 bg-c-primary hover:bg-c-primary/80 font-medium disabled:opacity-50"
              >
                <FiEdit3 className="w-4 h-4 mr-2" />
              Save Changes
              </Button>
            </>
          }
        />

        <SimpleModal
          show={showDelete}
          onClose={handleCloseDelete}
          title="Remove Network Configuration"
          children={
            <div className="space-y-6">
              <div className="p-6 bg-red-500/10 rounded-xl border border-red-500/20">
                <div className="flex items-center gap-3 mb-3">
                  <div className="p-2 bg-red-500/20 rounded-lg">
                    <FiAlertTriangle className="w-6 h-6 text-red-500" />
                  </div>
                  <h4 className="text-lg font-bold text-red-500">Permanent Deletion Warning</h4>
                </div>
                <p className="text-sm text-f-primary leading-relaxed">
                You are about to permanently remove the network configuration for{' '}
                  <span className="font-bold text-red-500">"{selectedNetwork}"</span>.
                </p>
              </div>

              <div className="space-y-3">
                <p className="text-sm text-neutral leading-relaxed">
                This action will:
                </p>
                <ul className="list-disc list-inside space-y-1 text-sm text-neutral ml-4">
                  <li>Remove the network from your available networks list</li>
                  <li>Delete all stored configuration for this network</li>
                  <li>Cannot be undone without re-adding the network manually</li>
                </ul>
              </div>

              <div className="p-4 bg-tertiary/10 rounded-xl border border-tertiary/20">
                <p className="text-xs text-neutral">
                  <strong>Note:</strong> If this is your currently active network, you'll be switched to{' '}
                  another available network automatically.
                </p>
              </div>
            </div>
          }
          footer={
            <>
              <Button onClick={handleCloseDelete} variant="secondary" className="px-6 py-2">
              Keep Network
              </Button>
              <Button
                variant="danger"
                onClick={onDelete}
                disabled={!selectedNetwork}
                className="px-6 py-2 bg-red-500 hover:bg-red-600 text-white font-medium disabled:opacity-50"
              >
                <FiTrash2 className="w-4 h-4 mr-2" />
              Delete Permanently
              </Button>
            </>
          }
        />
      </div>
    </div>
  );
}


