import { useEffect } from 'react';
import { useResource } from '../../custom/api';
import { useNavigate } from 'react-router-dom';
import { routeFor } from '../../utils';

import MyPlugin from './MyPlugin';

export interface IMassaPlugin {
  name: string;
  author: string;
  description: string;
  logo: string;
  version: string;
  status: string;
  updatable: boolean;
}

function MyStation() {
  const navigate = useNavigate();
  const {
    error,
    data: plugins,
    status,
  } = useResource<IMassaPlugin[]>('plugin-manager');
  const isLoadingData = status === 'loading';

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    }
  });

  return (
    <>
      {isLoadingData ? (
        <div className="mas-menu-active mb-4 text-neutral">Loading</div>
      ) : plugins && plugins.length > 0 ? (
        <div className="flex gap-2 flex-wrap">
          {plugins.map((plugin) => (
            <MyPlugin key={plugin.name} plugin={plugin} />
          ))}
        </div>
      ) : (
        <div className="mas-body mb-4 text-neutral">
          Browse the store below and manage plugin you've installed in this
          section
        </div>
      )}
    </>
  );
}

export default MyStation;
