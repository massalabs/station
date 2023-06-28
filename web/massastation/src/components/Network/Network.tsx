import { useEffect } from 'react';
import { Outlet, useParams, useNavigate } from 'react-router-dom';
import { useResource, usePost } from '../../custom/api';
import { routeFor } from '../../utils';

import { URL } from '../../const/url/url';
import { NetworkModel } from '../../models';
import { useNetworkStore } from '../../store/store';

interface NetworkRequest {
  network: string;
}

export function Network() {
  const navigate = useNavigate();
  const { network } = useParams();

  const { data: networkData, isLoading } = useResource<NetworkModel>(
    URL.PATH_NETWORKS,
  );

  const { mutate } = usePost<NetworkRequest>(`${URL.PATH_NETWORKS}/${network}`);

  const [setCurrentNetwork, setAvailableNetworks] = useNetworkStore((state) => [
    state.setCurrentNetwork,
    state.setAvailableNetworks,
  ]);

  useEffect(() => {
    if (networkData) {
      const { currentNetwork, availableNetworks } = networkData;

      if (network && !availableNetworks?.includes(network) && currentNetwork)
        navigate(routeFor(currentNetwork));
    }
  }, [networkData, network]);

  useEffect(() => {
    if (!isLoading && network && networkData) {
      setCurrentNetwork(network);
      setAvailableNetworks(networkData.availableNetworks);
      mutate({});
    }
  }, [isLoading, network]);

  return <Outlet />;
}
