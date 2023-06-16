import { useEffect } from 'react';
import { useResource } from '../../custom/api';
import { useNavigate } from 'react-router-dom';
import { routeFor } from '../../utils';
import Intl from '../../i18n/i18n';

import { MyPlugin } from './MyPlugin';

export interface IMassaPlugin {
  name: string;
  id: string;
  author: string;
  description: string;
  logo: string;
  version: string;
  home: string;
  status: string;
  updatable: boolean;
}

export function MyStation() {
  const navigate = useNavigate();

  const {
    error,
    data: myPlugins,
    isLoading,
    refetch: refetchPlugins,
    isRefetching,
    isSuccess,
  } = useResource<IMassaPlugin[]>('plugin-manager');

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    }
  }, [error, navigate]);

  function refreshList() {
    refetchPlugins();
  }
  console.log(
    'isRefetching',
    isRefetching,
    'isFetched',
    'isSuccess',
    isSuccess,
  );

  return (
    <>
      {isLoading || isRefetching ? (
        <div className="mas-body mb-4 text-neutral">
          {Intl.t('store.mystation.loading')}
        </div>
      ) : (
        <>
          {myPlugins && myPlugins.length > 0 ? (
            <div className="flex gap-4 flex-wrap">
              {isLoading} {isLoading} {isLoading}
              {myPlugins.map((plugin) => (
                <MyPlugin
                  key={plugin.id}
                  plugin={plugin}
                  fetchPlugins={refreshList}
                />
              ))}
            </div>
          ) : (
            <div className="mas-body mb-4 text-neutral">
              {Intl.t('store.mystation.browse')}
            </div>
          )}
        </>
      )}
    </>
  );
}
