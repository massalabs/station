import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { routeFor } from '../../utils';
import Intl from '../../i18n/i18n';

import { MyPlugin } from './MyPlugin';
import { UseQueryResult } from '@tanstack/react-query';

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

export function MyStation({
  getPlugins,
}: {
  getPlugins: UseQueryResult<IMassaPlugin[]>;
}) {
  const navigate = useNavigate();

  const {
    error,
    data: myPlugins,
    isLoading,
    refetch: refetchPlugins,
    isRefetching,
  } = getPlugins;

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    }
  }, [error, navigate]);

  return (
    <>
      {isLoading || isRefetching ? (
        <div className="mas-body mb-4 text-neutral">
          {Intl.t('store.mystation.loading')}
        </div>
      ) : (
        <>
          {myPlugins && myPlugins.length ? (
            <div className="flex gap-4 flex-wrap">
              {myPlugins.map((plugin) => (
                <MyPlugin
                  key={plugin.id}
                  plugin={plugin}
                  fetchPlugins={() => refetchPlugins()}
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
