import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { routeFor } from '@/utils/utils';
import { sortPlugins } from '@/utils/sortArray';

import Intl from '@/i18n/i18n';
import { UseQueryResult } from '@tanstack/react-query';

import StationPlugin from './StationPlugin';
import { MassaPluginModel } from '@/models';

function StationSection({
  getPlugins,
}: {
  getPlugins: UseQueryResult<MassaPluginModel[]>;
}) {
  const navigate = useNavigate();

  const { error, data: plugins, refetch, isLoading, isRefetching } = getPlugins;

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    }
  }, [error, navigate]);

  return (
    <>
      {isLoading || isRefetching ? null : (
        <>
          {plugins && plugins.length ? (
            <div className="flex gap-4 flex-wrap">
              {sortPlugins(plugins)?.map((plugin, index) => (
                <StationPlugin key={index} plugin={plugin} refetch={refetch} />
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
