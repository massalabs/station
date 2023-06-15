import { useEffect } from 'react';
import { useResource } from '../../../custom/api';
import { useNavigate } from 'react-router-dom';
import { routeFor } from '../../../utils';
import Intl from '../../../i18n/i18n';

import StationPlugin from './StationPlugin';

export interface IMassaPlugin {
  name: string;
  author: string;
  description: string;
  logo: string;
  version: string;
  status: string;
  updatable: boolean;
}

function StationSection() {
  const navigate = useNavigate();
  const {
    error,
    data: plugins,
    isLoading,
  } = useResource<IMassaPlugin[]>('plugin-manager');

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    }
  });

  return (
    <>
      {isLoading ? (
        <div className="mas-body mb-4 text-neutral">
          {Intl.t('store.mystation.loading')}
        </div>
      ) : plugins && plugins.length > 0 ? (
        <div className="flex gap-4 flex-wrap">
          {plugins.map((plugin) => (
            <StationPlugin key={plugin.name} plugin={plugin} />
          ))}
        </div>
      ) : (
        <div className="mas-body mb-4 text-neutral">
          {Intl.t('store.mystation.browse')}
        </div>
      )}
    </>
  );
}

export default StationSection;
