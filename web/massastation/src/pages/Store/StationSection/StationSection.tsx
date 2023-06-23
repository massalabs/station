import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { routeFor } from '../../../utils';
import { sortPlugins } from '../../../utils/sortArray';

import Intl from '../../../i18n/i18n';
import { UseQueryResult } from '@tanstack/react-query';

import StationPlugin from './StationPlugin';
import { IMassaPlugin } from '../../../../../shared/interfaces/IPlugin';

function StationSection({
  getPlugins,
}: {
  getPlugins: UseQueryResult<IMassaPlugin[]>;
}) {
  const navigate = useNavigate();

  const {
    error,
    data: plugins,
    refetch: refetchPlugins,
    isLoading,
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
        <div className="mas-body2 mb-4 text-neutral">
          {Intl.t('store.loading')}
        </div>
      ) : (
        <>
          {plugins && plugins.length ? (
            <div className="flex gap-4 flex-wrap">
              {sortPlugins(plugins)?.map((plugin, index) => (
                <StationPlugin
                  key={index}
                  plugin={plugin}
                  fetchPlugins={() => refetchPlugins()}
                />
              ))}
            </div>
          ) : (
            <div className="mas-body2 mb-4 text-neutral">
              {Intl.t('store.mystation.browse')}
            </div>
          )}
        </>
      )}
    </>
  );
}
export default StationSection;
